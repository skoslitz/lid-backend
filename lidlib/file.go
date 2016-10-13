package lidlib

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// NewFile constructs a new File based on a path and file info
func NewFile(path string, info os.FileInfo) *File {
	file := new(File)
	file.Path = filepath.ToSlash(path)
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
	case "region-list":
		f.Thema.Related = strings.Join([]string{ApiUrl, f.Path, "/themen"}, "")
		f.Exkursion.Related = strings.Join([]string{ApiUrl, f.Path, "/exkursionen"}, "")
	case "topic-list":
		var cRegionId = strings.Split(f.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("contentpath"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				f.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("contentpath"))}, "")
			}
		}

	case "excursion-list":
		var cRegionId = strings.Split(f.Id, "-")[0]
		regionContentDir := strings.Join([]string{viper.GetString("contentpath"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				f.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("contentpath"))}, "")
			}
		}
	}

}
