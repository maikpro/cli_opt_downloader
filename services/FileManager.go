package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func SaveFile(path string, filename string, file []byte) (string, error) {
	fullpath := fmt.Sprintf("%s/%s", path, filename)

	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(fullpath), 0o755)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Write the data to a file
	err = os.WriteFile(fullpath, file, 0o644)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return filepath.Abs(fullpath)
}
