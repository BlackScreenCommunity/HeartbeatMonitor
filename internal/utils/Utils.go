package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
)

// Checks if a string exists in an array of strings.
func IsArrayContainString(slice []string, str string) bool {
	return slices.Contains(slice, str)
}

// Convert oject to strings array
func ConvertInterfaceArrayToStringArray(interfaceArray []interface{}) []string {
	var stringSlice []string
	for _, v := range interfaceArray {
		if str, ok := v.(string); ok {
			stringSlice = append(stringSlice, str)
		}
	}

	return stringSlice
}

// Merges all CSS files in a directory
func MergeFilesWithExtension(dir string, fileExtension string) ([]byte, error) {
	var buffer bytes.Buffer

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == fileExtension {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer func() {
				if err := file.Close(); err != nil {
					log.Printf("Error while closing file: %v", err)
				}
			}()

			_, err = io.Copy(&buffer, file)
			if err != nil {
				log.Printf("Error while copying content: %v", err)
			}
			buffer.WriteString("\n")
		}
		return nil
	})

	return buffer.Bytes(), err
}
