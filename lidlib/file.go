package lidlib

import (
	"os"
	"strings"
)

type File struct {
	Id    string `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

type Files []*File

// NewFile constructs a new File based on a path and file info
func NewFile(path string, info os.FileInfo) *File {
	file := &File{Path: path}
	file.Load(info)
	return file
}
func (f *File) Load(info os.FileInfo) {
	id := strings.Split(info.Name(), "-")
	f.Id = 	id[0]
	f.Name = info.Name()
	f.IsDir = info.IsDir()
	f.Size = info.Size()
	f.ModTime = info.ModTime().Format("02/01/2006")
}
