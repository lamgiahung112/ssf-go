package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ssf/models"
	"ssf/storage"
	views "ssf/views/helpers"

	expressgo "github.com/lamgiahung112/express-go"
)

type FilesController struct {
	storedFileRepo *models.StoredFileRepository
	fileStorage    *storage.Storage
}

func NewFilesController() *FilesController {
	// Get the current working directory (project root)
	rootPath, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	opts := storage.StorageOptions{
		RootPath:          filepath.Join(rootPath, "files"),
		FilenameConverter: storage.EncodeFilenameConverterFunc,
	}
	storage := storage.Storage{
		StorageOptions: opts,
	}
	return &FilesController{
		storedFileRepo: models.NewStoredFileRepository(),
		fileStorage:    &storage,
	}
}

func (fc *FilesController) HandleSubmitShareFile(w http.ResponseWriter, r *http.Request, context *expressgo.RequestContext) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(50 << 20) // 50 MB max memory
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	defer file.Close()

	// Get the key from form data
	password := r.FormValue("password")
	if password == "" {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// Create a new StoredFile
	storedFile := fc.storedFileRepo.NewStoredFile(header.Filename)

	// Save the file metadata
	err = fc.storedFileRepo.Save(storedFile)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	err = fc.fileStorage.Write(storedFile.Slug, password, file)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	context.Set("file_link", "https://lghung.io.vn/files/"+storedFile.Slug)
	context.Set("file_password", password)
	context.Set("title_for_layout", "Your file is secured!")
	views.Render("file_store_success.html", context, w, r)
}

func (fc *FilesController) HandleShowSharedFile(w http.ResponseWriter, r *http.Request, context *expressgo.RequestContext) {
	slug := r.PathValue("slug")
	if len(slug) == 0 {
		context.Set("error", "File is not found")
		context.Set("title_for_layout", "Your file is broken or not found!")
		views.Render("files_by_slug.html", context, w, r)
		return
	}

	fileEntity, err := fc.storedFileRepo.FindBySlug(slug)
	if err != nil {
		context.Set("error", "File is not found")
		context.Set("title_for_layout", "Your file is broken or not found!")
		views.Render("files_by_slug.html", context, w, r)
		return
	}

	context.Set("file", fileEntity)
	context.Set("title_for_layout", "Secure file sharing")
	views.Render("files_by_slug.html", context, w, r)
}

func (fc *FilesController) CheckFilePassword(w http.ResponseWriter, r *http.Request, context *expressgo.RequestContext) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Error parsing data: " + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slug := r.FormValue("slug")
	password := r.FormValue("password")

	if len(slug) == 0 || len(password) == 0 {
		w.Write([]byte("File or password missing"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = fc.storedFileRepo.FindBySlug(slug)
	if err != nil {
		w.Write([]byte("Error retrieving file: " + err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedFilename := fc.fileStorage.FilenameConverter(slug, password)

	// Check if the file exists
	filePath := filepath.Join(fc.fileStorage.StorageOptions.RootPath, hashedFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found or incorrect password", http.StatusNotFound)
		return
	}

	// If we've reached this point, the file exists and the password is correct
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File exists and password is correct"))
}

func (fc *FilesController) Download(w http.ResponseWriter, r *http.Request, context *expressgo.RequestContext) {
	slug := r.PathValue("slug")
	password := r.URL.Query().Get("password")
	if len(slug) == 0 || len(password) == 0 {
		w.Write([]byte("Unexpected error"))
		return
	}

	storedFile, err := fc.storedFileRepo.FindBySlug(slug)
	if err != nil {
		w.Write([]byte("Error retrieving file: " + err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedFilename := fc.fileStorage.FilenameConverter(slug, password)

	// Check if the file exists
	filePath := filepath.Join(fc.fileStorage.StorageOptions.RootPath, hashedFilename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found or incorrect password", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+storedFile.OriginalFilename)
	w.Header().Set("Content-Type", "application/octet-stream")
	err = fc.fileStorage.Read(slug, password, w)
	if err != nil {
		w.Write([]byte("Error retrieving file"))
	}
}
