package storage

import (
	"bytes"
	"testing"
)

func TestStorage(t *testing.T) {
	opts := &StorageOptions{
		RootPath:          "/home/jmt/projects/ssF/files",
		FilenameConverter: EncodeFilenameConverterFunc,
	}
	storage := &Storage{
		StorageOptions: *opts,
	}

	data := []byte(`
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		Hello hello This is a text larger than 512 bytes to test if it actually works...
		`)
	err := storage.Write("lamgiahung.html", "1233", bytes.NewReader(data))

	if err != nil {
		t.Errorf(err.Error())
	}

	// Test Read
	var buf bytes.Buffer
	err = storage.Read("lamgiahung.html", "1233", &buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Compare read data with original data
	if buf.String() != string(data) {
		t.Errorf("Read data does not match original data")
	}
}
