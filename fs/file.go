package fs

import (
	"encoding/base64"
	"os"
)

// Decodes the base64 data and writes it to a file
func SaveBase64File(filename string, data string) error {
	// Decodes the file
	dec, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	// Creates the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writes the file
	if _, err := file.Write(dec); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

// Removes a file
func RemoveFile(filename string) error {
	return os.Remove(filename)
}
