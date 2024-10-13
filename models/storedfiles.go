package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"mime"
	"path/filepath"
	"ssf/db"
	"strings"
	"time"
	"unicode"
)

type StoredFile struct {
	ID               int
	Slug             string
	OriginalFilename string
	MimeType         string
	UploadedAt       time.Time
}

type StoredFileRepository struct{}

func NewStoredFileRepository() *StoredFileRepository {
	return &StoredFileRepository{}
}

func (repo *StoredFileRepository) NewStoredFile(
	originalFilename string,
) *StoredFile {
	return &StoredFile{
		Slug:             createSlugWithRandomString(originalFilename),
		OriginalFilename: originalFilename,
		MimeType:         getMimeType(originalFilename),
		UploadedAt:       time.Now(),
	}
}

func (repo *StoredFileRepository) Save(f *StoredFile) error {
	query := `
		INSERT INTO stored_files (slug, original_filename, mime_type, uploaded_at)
		VALUES (?, ?, ?, ?)
	`

	// Execute the SQL statement
	result, err := db.DB.Exec(query, f.Slug, f.OriginalFilename, f.MimeType, f.UploadedAt)
	if err != nil {
		return fmt.Errorf("error saving stored file: %v", err)
	}

	// Get the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Set the ID of the StoredFile
	f.ID = int(lastInsertID)

	return nil
}

func (repo *StoredFileRepository) FindBySlug(slug string) (*StoredFile, error) {
	query := `
		SELECT id, slug, original_filename, mime_type, uploaded_at 
		FROM stored_files 
		WHERE slug = ? 
		LIMIT 1
	`

	var file StoredFile
	err := db.DB.QueryRow(query, slug).Scan(
		&file.ID,
		&file.Slug,
		&file.OriginalFilename,
		&file.MimeType,
		&file.UploadedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("error retrieving file: %v", err)
	}

	return &file, nil
}

func (repo *StoredFileRepository) DeleteById(id int) error {
	query := `
		DELETE FROM stored_files 
		WHERE id = ?
	`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting StoredFiles (id: %d): %s", id, err)
	}
	return nil
}

func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// If MIME type is not found, default to "application/octet-stream"
		mimeType = "application/octet-stream"
	}
	return mimeType
}

func createSlugWithRandomString(s string) string {
	slug := createSlug(s)
	randomString := generateRandomString(6) // Generate a 6-character random string
	return fmt.Sprintf("%s-%s", slug, randomString)
}

func createSlug(s string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(s)
	slug = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		return -1
	}, slug)

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
