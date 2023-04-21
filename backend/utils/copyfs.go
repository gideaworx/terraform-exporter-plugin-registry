package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ExportFS(src fs.FS, targetDir string) error {
	return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, e error) error {
		itemPath := filepath.Join(targetDir, path)
		if d.IsDir() {
			return os.MkdirAll(itemPath, 0777)
		}

		inputBytes, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}

		return os.WriteFile(itemPath, inputBytes, 0666)
	})
}
