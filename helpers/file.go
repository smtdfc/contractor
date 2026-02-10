package helpers

import (
	"os"
	"path/filepath"
)

func ReadTextFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func WriteTextFile(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}
