package cmd

import (
	"bytes"
	"os"
	"path/filepath"
)

func replace(path, pattern, old, new string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return replaceWalkFn(path, info, pattern, []byte(old), []byte(new))
	})
}

func replaceWalkFn(path string, info os.FileInfo, pattern string, old, new []byte) (err error) {
	var matched bool
	if matched, err = filepath.Match(pattern, info.Name()); err != nil {
		return
	}

	if matched {
		cleanedPath := filepath.Clean(path)

		var oldContent []byte
		if oldContent, err = os.ReadFile(cleanedPath); err != nil {
			return
		}

		if err = os.WriteFile(cleanedPath, bytes.ReplaceAll(oldContent, old, new), 0); err != nil {
			return
		}
	}

	return
}
