package files

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func IsEmpty(dir string) (bool, error) {
	paths, err := GetPaths(dir)
	if err != nil {
		return false, err
	}
	return len(paths) < 1, nil
}

func Move(src, dst string) error {
	return os.Rename(src, dst)
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func MarshalJSON[T any](dataStruct T) ([]byte, error) {
	return json.MarshalIndent(dataStruct, "", "  ")
}

func GetPaths(folder string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not walk over %s: %v", folder, err)
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return err
	})
	return files, err
}
