package lidlib

import (
	"os"
	"regexp"
	"strings"
)

type File struct {
	Id       string `json:"id"`
	FileName string `json:"fileName"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	Size     int64  `json:"size"`
	ModTime  string `json:"modTime"`
	Edition  string `json:"edition"`
	Link     string `json:"link"`
}

type Files []*File

// NewFile constructs a new File based on a path and file info
func NewFile(path string, info os.FileInfo) *File {
	file := &File{Path: path}
	file.Load(info)
	return file
}
func (f *File) Load(info os.FileInfo) {

	searchTerm := `(([\s\S]+?)[-]{1})`
	re := regexp.MustCompile(searchTerm)
	prefix := re.FindStringSubmatch(string(info.Name()))[1]

	f.Id = info.Name()
	f.FileName = info.Name()
	f.IsDir = info.IsDir()
	f.Size = info.Size()
	f.ModTime = info.ModTime().Format("02/01/2006")
	f.Edition = strings.Split(strings.Split(info.Name(), "-")[0], "_")[0]
	_nameTrimSuffix := strings.TrimSuffix(info.Name(), ".md")
	_nameTrimPrefix := strings.TrimPrefix(_nameTrimSuffix, prefix)
	_nameTrimDash := strings.Replace(_nameTrimPrefix, "-", " ", -1)
	_nameToTitle := strings.ToTitle(_nameTrimDash)
	f.Name = _nameToTitle

}
