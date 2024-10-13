package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type StorageOptions struct {
	RootPath          string
	FilenameConverter FilenameConverterFunc
}

type Storage struct {
	StorageOptions
}

func prepareAES(key string) cipher.AEAD {
	hashedKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashedKey[:])
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	return gcm
}

func (s *Storage) Read(slug string, key string, w io.Writer) error {
	filename := s.FilenameConverter(slug, key)

	f, err := os.Open(s.RootPath + "/" + filename)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer f.Close()

	gcm := prepareAES(key)
	nonceSize := gcm.NonceSize()

	for {
		// Read nonce
		nonce := make([]byte, nonceSize)
		_, err := io.ReadFull(f, nonce)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading nonce: %s", err)
		}

		// Read encrypted chunk
		ciphertext := make([]byte, 512+gcm.Overhead())
		n, err := f.Read(ciphertext)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading ciphertext: %s", err)
		}

		// Decrypt chunk
		plaintext, err := gcm.Open(nil, nonce, ciphertext[:n], nil)
		if err != nil {
			return fmt.Errorf("error decrypting: %s", err)
		}

		w.Write(plaintext)

		if err == io.EOF {
			break
		}
	}

	return nil
}

func (s *Storage) Write(slug string, key string, r io.Reader) error {
	filename := s.FilenameConverter(slug, key)

	f, err := os.Create(s.RootPath + "/" + filename)
	if err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}
	defer f.Close()

	buf := make([]byte, 512)
	gcm := prepareAES(key)
	nonce := make([]byte, gcm.NonceSize())

	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error writing %s: %s", slug, err)
		}

		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return fmt.Errorf("error writing %s: %s", slug, err)
		}
		f.Write(gcm.Seal(nonce, nonce, buf[:n], nil))
	}
	return nil
}
