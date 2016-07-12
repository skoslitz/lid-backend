package lidlib

import (
	"os"
)

type Region struct {
	Self string `json:"self"`
}

type Thema struct {
	Self string `json:"self"`
}

type Exkursion struct {
	Self string `json:"self"`
}
type Link struct {
	Self string `json:"self"`
}

type Attribute struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

type Relationship struct {
	Regions     Region    `json:"region"`
	Themen      Thema     `json:"themen"`
	Exkursionen Exkursion `json:"exkursionen"`
}

type File struct {
	Id            string       `json:"id"`
	Type          string       `json:"type"`
	Links         Link         `json:"links"`
	Attributes    Attribute    `json:"attributes"`
	Relationships Relationship `json:"relationships"`
}

type Files []*File

// NewFile constructs a new File based on a path and file info
func NewFile(path string, info os.FileInfo) *File {
	file := new(File)
	file.Attributes.Path = path
	file.Load(info)
	return file
}
func (f *File) Load(info os.FileInfo) {
	f.Id = info.Name()
	f.Attributes.Name = info.Name()
	f.Attributes.IsDir = info.IsDir()
	f.Attributes.Size = info.Size()
	f.Attributes.ModTime = info.ModTime().Format("02/01/2006")
}
