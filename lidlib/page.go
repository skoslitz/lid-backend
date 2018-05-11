package lidlib

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/parser"
	"github.com/spf13/viper"
)

type PageManager interface {
	Read(fp string) (*PageFile, error)
	Create(fp string, fm Frontmatter, content []byte) (*PageFile, error)
	Update(fp string, fm Frontmatter, content []byte) (*PageFile, error)
	Delete(fp string) error
}

func NewPage() *Page {
	return &Page{}
}

// ReadPage reads a page from disk
func (p Page) Read(fp string) (*PageFile, error) {
	// open the file for reading
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// use the Hugo parser lib to read the contents
	parser, err := parser.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	// get the metadata
	rawdata, err := parser.Metadata()
	if err != nil {
		return nil, err
	}

	// convert the interface{} into map[string]interface{}
	metadata, err := cast.ToStringMapE(rawdata)
	if err != nil {
		return nil, err
	}

	// retrieve file infos
	fileStat, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}

	// create and assemble a new Page instance
	pagefile := &PageFile{
		Id: filepath.Base(fp),
		Attribute: Attribute{
			Name:     filepath.Base(fp),
			Path:     filepath.ToSlash(fp),
			Size:     fileStat.Size(),
			ModTime:  fileStat.ModTime().Format("02/01/2006"),
			Metadata: metadata,
			Content:  string(parser.Content()),
		},
	}

	return pagefile, nil

}

// CreatePage creates a new file and saves page content to it
func (p Page) Create(fp string, fm Frontmatter, content []byte) (*PageFile, error) {

	// create and assemble a new Page instance
	pagefile := &PageFile{
		Id: filepath.Base(fp),
		Attribute: Attribute{
			Path:     fp,
			Metadata: fm,
			Content:  string(content),
		},
	}

	// save page to disk
	err := pagefile.Save()
	if err != nil {
		return nil, err
	}

	// NOTE
	// pagefile os.Stat() is not performed
	// -> size and modTime are not set

	return pagefile, nil
}

// UpdatePage changes the content of an existing page
func (p Page) Update(fp string, fm Frontmatter, content []byte) (*PageFile, error) {

	// delete existing page
	err := p.Delete(fp)
	if err != nil {
		return nil, err
	}

	// create and assemble a new Page instance
	pagefile := &PageFile{
		Id: filepath.Base(fp),
		Attribute: Attribute{
			Path:     fp,
			Metadata: fm,
			Content:  string(content),
		},
	}

	// save page to disk
	err = pagefile.Save()
	if err != nil {
		return nil, err
	}

	return pagefile, nil
}

// Delete deletes a page
func (p Page) Delete(fp string) error {

	// check that file exists
	info, err := os.Stat(fp)
	if err != nil {
		return err
	}

	// that file is a directory
	if info.IsDir() {
		return errors.New("DeletePage cannot delete directories")
	}

	// remove the directory
	return os.Remove(fp)
}

// Saves a page
func (p *PageFile) Save() error {
	// create new hugo page
	page, err := hugolib.NewPage(p.Path)
	if err != nil {
		return err
	}

	// set attributes
	page.SetSourceMetaData(p.Metadata, TOML)
	page.SetSourceContent([]byte(p.Content))

	// save page
	return page.SafeSaveSourceAs(p.Path)
}

// Sets relationship links per type and id
func (p *PageFile) SetRelationship(ApiUrl string) {
	switch p.Type {
	case "region":
		p.Thema.Related = strings.Join([]string{ApiUrl, p.Path, "/themen"}, "")
		p.Exkursion.Related = strings.Join([]string{ApiUrl, p.Path, "/exkursionen"}, "")
		p.Region.Related = strings.Join([]string{ApiUrl, "dir", "/regionen"}, "")
	case "topic":
		var cRegionId = strings.Split(p.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("contentpath"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				p.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("contentpath"))}, "")
			}
		}
	case "excursion":
		var cRegionId = strings.Split(p.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("contentpath"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				p.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("contentpath"))}, "")
			}
		}
	}
}
