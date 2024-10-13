package models

import (
	"ssf/db"
	"testing"

	"github.com/joho/godotenv"
)

func TestStoredFiles(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf(err.Error())
	}
	db.ConnectDatabase()
	storedFileRepository := NewStoredFileRepository()
	storedFile := storedFileRepository.NewStoredFile("lamgiahung.html")

	err = storedFileRepository.Save(storedFile)

	if err != nil {
		t.Errorf(err.Error())
	}

	savedStoreFile, err := storedFileRepository.FindBySlug(storedFile.Slug)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = storedFileRepository.DeleteById(savedStoreFile.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
}
