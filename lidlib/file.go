package lidlib

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Link struct {
	Self    string `json:"self,omitempty"`
	Related string `json:"related,omitempty"`
}

type Region struct {
	Link `json:"links"`
}

type Thema struct {
	Link `json:"links"`
}

type Exkursion struct {
	Link `json:"links"`
}

type Attribute struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Size     int64       `json:"size"`
	ModTime  string      `json:"modTime"`
	Metadata Frontmatter `json:"metadata,omitempty"`
	Content  string      `json:"content,omitempty"`
}

type Relationship struct {
	Region    `json:"region"`
	Thema     `json:"themen"`
	Exkursion `json:"exkursionen"`
}

type File struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Link         `json:"links"`
	Attribute    `json:"attributes"`
	Relationship `json:"relationships"`
}

type Files []*File

// NewFile constructs a new File based on a path and file info
func NewFile(path string, info os.FileInfo) *File {
	file := new(File)
	file.Path = path
	file.Load(info)
	return file
}

func (f *File) Load(info os.FileInfo) {
	f.Id = info.Name()
	f.Name = info.Name()
	f.IsDir = info.IsDir()
	f.Size = info.Size()
	f.ModTime = info.ModTime().Format("02/01/2006")
}

func (f *File) SetRelationship(ApiUrl string) {
	switch f.Type {
	case "regionen":
		f.Thema.Related = strings.Join([]string{ApiUrl, f.Path, "/themen"}, "")
		f.Exkursion.Related = strings.Join([]string{ApiUrl, f.Path, "/exkursionen"}, "")
	case "themen":
		var cRegionId = strings.Split(f.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("ContentDir"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				f.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("ContentDir"))}, "")
			}
		}

	case "exkursionen":
		var cRegionId = strings.Split(f.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("ContentDir"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				f.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("ContentDir"))}, "")
			}
		}
	}

}
