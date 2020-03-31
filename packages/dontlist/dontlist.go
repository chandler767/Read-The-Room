// Package dontlist disables directory listing with http.FileServer using custom filesystem. Returns error if not found or directory.
package dontlist

import (
	"net/http"
	"os"
)

type DontListFiles struct {
	Fs http.FileSystem
}

func (fs DontListFiles) Open(name string) (http.File, error) {
	f, err := fs.Fs.Open(name)
	if err != nil {
		return nil, os.ErrNotExist
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, os.ErrNotExist
}
