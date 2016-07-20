package lidlib

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"github.com/spf13/viper"
)

const TOML = '+'
const YAML = '-'

// Frontmatter stores encodeable data
type Frontmatter map[string]interface{}

// Page represents a markdown file
type PageFile struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Link         `json:"links"`
	Attribute    `json:"attributes"`
	Relationship `json:"relationships"`
}

type PageManager interface {
	Read(fp string) (*PageFile, error)
	Create(fp string, fm Frontmatter, content []byte) (*PageFile, error)
	Update(fp string, fm Frontmatter, content []byte) (*PageFile, error)
	Delete(fp string) error
}

type Page struct{}

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

	// create and assemble a new Page instance
	pagefile := new(PageFile)
	pagefile.assemblePageInstance(fp, metadata, parser.Content())

	return pagefile, nil

}

// CreatePage creates a new file and saves page content to it
func (p Page) Create(dirname string, fm Frontmatter, content []byte) (*PageFile, error) {

	// get title from metadata
	title, err := getTitle(fm)
	if err != nil {
		return nil, err
	}

	// the filepath for the page
	fp := generateFilePath(dirname, title)

	// create and assemble a new Page instance
	pagefile := new(PageFile)
	pagefile.assemblePageInstance(fp, fm, content)

	// save page to disk
	err = pagefile.Save()
	if err != nil {
		return nil, err
	}

	return pagefile, nil
}

// UpdatePage changes the content of an existing page
func (p Page) Update(fp string, fm Frontmatter, content []byte) (*PageFile, error) {

	// get title from metadata
	title, err := getTitle(fm)
	if err != nil {
		return nil, err
	}

	// delete existing page
	err = p.Delete(fp)
	if err != nil {
		return nil, err
	}

	// the filepath for the page
	dirname := filepath.Dir(fp)
	fp = generateFilePath(dirname, title)

	// create and assemble a new Page instance
	pagefile := new(PageFile)
	pagefile.assemblePageInstance(fp, fm, content)

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

func (p *PageFile) SetRelationship(ApiUrl string) {
	switch p.Type {
	case "regionen":
		p.Thema.Related = strings.Join([]string{ApiUrl, p.Path, "/themen"}, "")
		p.Exkursion.Related = strings.Join([]string{ApiUrl, p.Path, "/exkursionen"}, "")
		p.Region.Related = strings.Join([]string{ApiUrl, "dir", "/regionen"}, "")
	case "themen":
		var cRegionId = strings.Split(p.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("ContentDir"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				p.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("ContentDir"))}, "")
			}
		}
	case "exkursionen":
		var cRegionId = strings.Split(p.Id, "_")[0]
		regionContentDir := strings.Join([]string{viper.GetString("ContentDir"), "regionen"}, "")

		// read contents of regionContentDir
		var contents Files
		var rc = new(Dir)
		contents, _ = rc.Read(regionContentDir)

		for _, item := range contents {
			var cid = strings.Split(item.Id, "-")[0]
			if cRegionId == cid {
				p.Region.Related = strings.Join([]string{ApiUrl, "page/", strings.TrimPrefix(item.Path, viper.GetString("ContentDir"))}, "")
			}
		}
	}
}

func (p *PageFile) assemblePageInstance(fp string, fm Frontmatter, content []byte) {
	p.Id = filepath.Base(fp)
	p.Path = fp
	p.Metadata = fm
	p.Content = string(content)

	fileStat, _ := os.Stat(fp)
	p.ModTime = fileStat.ModTime().Format("02/01/2006")
	p.Size = fileStat.Size()
}

// generateFilePath generates a filepath based on a page title
// if the filename already exists, add a number on the end
// if that exists, increment the number by one until we find a filename
// that doesn't exist
func generateFilePath(dirname, title string) (fp string) {
	count := 0

	for {

		// combine title with count
		name := title
		if count != 0 {
			name += " " + strconv.Itoa(count)
		}

		// join filename with dirname
		filename := sanitize.Path(name + ".md")
		fp = filepath.Join(dirname, filename)

		// only stop looping when file doesn't already exist
		if _, err := os.Stat(fp); err != nil {
			break
		}

		// try again with a different number
		count += 1
	}

	return fp
}

func getTitle(fm Frontmatter) (string, error) {

	// check that title has been specified
	t, ok := fm["title"]
	if ok == false {
		return "", errors.New("page[meta].title must be specified")
	}

	// check that title is a string
	title, ok := t.(string)
	if ok == false {
		return "", errors.New("page[meta].title must be a string")
	}

	return title, nil
}
