package goefs

import (
	"embed"
	"io/fs"
	"os"
	"path"
)

func get_dir(f *embed.FS, folder string) []fs.DirEntry {
	if list, err := f.ReadDir(folder); err != nil {
		panic(err)
	} else {
		return list
	}
}

func is_exist(fpath string) bool {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}
	return true
}

func createDirFromEmbedFS(f *embed.FS, parent string) error {
	files := get_dir(f, parent)

	for _, info := range files {
		fpath := path.Join(parent, info.Name())
		if !is_exist(fpath) {
			if info.IsDir() {
				if err := os.Mkdir(fpath, fs.ModePerm); err != nil {
					return err
				}
			} else if data, err := f.ReadFile(fpath); err != nil {
				return err
			} else if err := os.WriteFile(fpath, data, fs.ModePerm); err != nil {
				return err
			}
		}

		if info.IsDir() {
			if err := createDirFromEmbedFS(f, fpath); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create directory from embed FS
func CreateDirFromEmbedFS(f *embed.FS) error {
	return createDirFromEmbedFS(f, ".")
}
