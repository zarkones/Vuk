package helpers

import (
	"os"
	"path/filepath"
)

func WalkDir(rootDirectoryPath string) (files []string, err error) {
	records, err := os.ReadDir(rootDirectoryPath)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		filePath := filepath.Join(rootDirectoryPath, record.Name())

		if record.IsDir() {
			subfiles, err := WalkDir(filePath)
			if err != nil {
				return files, err
			}

			files = append(files, subfiles...)
			continue
		}

		files = append(files, filePath)
	}

	return files, nil
}
