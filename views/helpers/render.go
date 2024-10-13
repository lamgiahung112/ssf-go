package views

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	expressgo "github.com/lamgiahung112/express-go"
)

var tmplt *template.Template

func PrepareTemplates() {
	tmplt = template.New("")
	err := filepath.Walk("./views", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = tmplt.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}
		return err
	})

	if err != nil {
		panic(err)
	}
}

func Render(path string, data *expressgo.RequestContext, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmplt.ExecuteTemplate(w, path, data.GetRaw())
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Sorry, the application encountered an unexpected error", http.StatusInternalServerError)
		return
	}
}
