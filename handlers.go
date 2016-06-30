package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kennygrant/sanitize"
	"github.com/skoslitz/lid-backend/lidlib"
)

type Handlers struct {
	Config lidlib.ConfigManager
	Dir    lidlib.DirManager
	Page   lidlib.PageManager

	ContentDir string
	AssetsDir  string
}

/*
  Directories
*/

type readDirResponse struct {
	Data lidlib.Files `json:"data"`
}

type createDirResponse struct {
	Dir *lidlib.File `json:"dir"`
}

type updateDirResponse struct {
	Dir *lidlib.File `json:"dir"`
}

// readDir reads contents of a directory
func (h Handlers) ReadDir(w http.ResponseWriter, r *http.Request) {

	var ApiPageUrl = strings.Join([]string{r.Host, "/api/page/"}, "")

	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// try and read contents of dir
	var contents lidlib.Files
	contents, err = h.Dir.Read(fp)
	if err != nil {
		errDirNotFound.Write(w)
		return
	}

	// trim content prefix
	for _, item := range contents {
		item.Path = strings.TrimPrefix(item.Path, h.ContentDir)
		item.Link = strings.Join([]string{ApiPageUrl, item.Path}, "")
	}

	printJson(w, &readDirResponse{Data: contents})
}

// readDir reads contents of a directory
func (h Handlers) ReadDirEdition(w http.ResponseWriter, r *http.Request) {

	var ApiPageUrl = strings.Join([]string{r.Host, "/api/page/"}, "")

	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// read edition to filter
	var editionNumber string
	editionNumber = mux.Vars(r)["edition"]
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// try and read contents of dir
	var contents lidlib.Files
	contents, err = h.Dir.Read(fp)
	if err != nil {
		errDirNotFound.Write(w)
		return
	}

	// trim content prefix
	for i, item := range contents {
		if editionNumber == item.Edition {
			item.Path = strings.TrimPrefix(item.Path, h.ContentDir)
			item.Link = strings.Join([]string{ApiPageUrl, item.Path}, "")
		} else {
			contents[i] = nil
		}
	}

	printJson(w, &readDirResponse{Data: contents})

}

// createDir creates a directory
func (h Handlers) CreateDir(w http.ResponseWriter, r *http.Request) {

	// combine parent and dirname
	parent := mux.Vars(r)["path"]
	dirname := sanitize.Path(r.FormValue("dir[name]"))
	fp := filepath.Join(parent, dirname)

	// check that it is a valid path
	fp, err := h.fixPathWithDir(fp, h.ContentDir)
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// check if dir already exists
	if fileExists(fp) || dirExists(fp) {
		errDirConflict.Write(w)
		return
	}

	// make directory
	dir, err := h.Dir.Create(fp)
	if err != nil {
		wrapError(err).Write(w)
		return
	}

	// trim content prefix
	dir.Path = strings.TrimPrefix(dir.Path, h.ContentDir)

	// print info
	printJson(w, &createDirResponse{Dir: dir})
}

// updateDir renames a directory
func (h Handlers) UpdateDir(w http.ResponseWriter, r *http.Request) {
	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// check that the specified directory is not the root content folder
	if fp == h.ContentDir {
		errInvalidDir.Write(w)
		return
	}

	// check that directory exists
	if dirExists(fp) == false {
		errDirNotFound.Write(w)
		return
	}

	// combine parent dir with dir name
	parent := filepath.Dir(fp)
	dirname := sanitize.Path(r.FormValue("dir[name]"))
	dest := filepath.Join(parent, dirname)

	// rename directory
	dir, err := h.Dir.Update(fp, dest)
	if err != nil {
		wrapError(err).Write(w)
		return
	}

	// print info
	printJson(w, &updateDirResponse{Dir: dir})
}

// TODO: make this more safe with captcha
// DeleteDir deletes a directory
func (h Handlers) DeleteDir(w http.ResponseWriter, r *http.Request) {
	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// check that the specified directory is not the root content folder
	if fp == h.ContentDir {
		errInvalidDir.Write(w)
		return
	}

	// remove directory
	if err = h.Dir.Delete(fp); err != nil {
		errDirNotFound.Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
  Pages
*/

type readPageResponse struct {
	Page *lidlib.PageFile `json:"page"`
}

type createPageResponse struct {
	Page *lidlib.PageFile `json:"page"`
}

type updatePageResponse struct {
	Page *lidlib.PageFile `json:"page"`
}

// readPage reads page data
func (h Handlers) ReadPage(w http.ResponseWriter, r *http.Request) {
	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		errInvalidDir.Write(w)
		return
	}

	// read page from disk
	page, err := h.Page.Read(fp)
	if err != nil {
		errPageNotFound.Write(w)
		return
	}

	// trim content prefix from path
	page.Path = strings.TrimPrefix(page.Path, h.ContentDir)

	// add asset path to relationships
	bandNummer := fmt.Sprintf("%d", page.Metadata["bandnummer"])
	var ApiAssetUrl = strings.Join([]string{r.Host, "/assets/img/", bandNummer}, "")
	page.Relationships.Assets = ApiAssetUrl

	// Relationships ----------------------------------------------------------------------

	// search for region related topics/excursions
	// TODO: make this a function
	if strings.Contains(page.Path, "regionen") {

		bandNummer := fmt.Sprintf("%d", page.Metadata["bandnummer"])
		var ApiPageUrl = strings.Join([]string{r.Host, "/api/page/"}, "")
		var themenPath = strings.Join([]string{h.ContentDir, "themen"}, "")
		var exkursionenPath = strings.Join([]string{h.ContentDir, "exkursionen"}, "")

		// try and read contents of themen dir
		var themenContents lidlib.Files
		themenContents, err = h.Dir.Read(themenPath)
		if err != nil {
			errDirNotFound.Write(w)
			return
		}

		// append related contents to ressource
		for _, item := range themenContents {

			if item.Edition == bandNummer {
				page.Relationships.Themen = append(page.Relationships.Themen, strings.Join([]string{ApiPageUrl, strings.TrimPrefix(item.Path, h.ContentDir)}, ""))
			}
		}

		// try and read contents of exkursionen dir
		var exkursionenContents lidlib.Files
		exkursionenContents, err = h.Dir.Read(exkursionenPath)
		if err != nil {
			errDirNotFound.Write(w)
			return
		}

		// append related contents to ressource
		for _, item := range exkursionenContents {
			if item.Edition == bandNummer {
				page.Relationships.Exkursionen = append(page.Relationships.Exkursionen, strings.Join([]string{ApiPageUrl, strings.TrimPrefix(item.Path, h.ContentDir)}, ""))
			}
		}

	}

	// search for region related excursions
	// TODO: make this a function
	if strings.Contains(page.Path, "exkursionen") {

		_bandNummer := fmt.Sprint(page.Metadata["id"])
		bandNummer := strings.Split(_bandNummer, "_")[0]
		var ApiPageUrl = strings.Join([]string{r.Host, "/api/page/"}, "")
		var regionenPath = strings.Join([]string{h.ContentDir, "regionen"}, "")

		// try and read contents of exkursionen dir
		var regionenContents lidlib.Files
		regionenContents, err = h.Dir.Read(regionenPath)
		if err != nil {
			errDirNotFound.Write(w)
			return
		}

		// append related contents to ressource
		for _, item := range regionenContents {

			if item.Edition == bandNummer {
				page.Relationships.Region = append(page.Relationships.Region, strings.Join([]string{ApiPageUrl, strings.TrimPrefix(item.Path, h.ContentDir)}, ""))
			}
		}

	}

	// search for region related excursions
	// TODO: make this a function
	if strings.Contains(page.Path, "themen") {

		_bandNummer := fmt.Sprint(page.Metadata["id"])
		bandNummer := strings.Split(_bandNummer, "_")[0]
		var ApiPageUrl = strings.Join([]string{r.Host, "/api/page/"}, "")
		var regionenPath = strings.Join([]string{h.ContentDir, "regionen"}, "")

		// try and read contents of themen dir
		var regionenContents lidlib.Files
		regionenContents, err = h.Dir.Read(regionenPath)
		if err != nil {
			errDirNotFound.Write(w)
			return
		}

		// append related contents to ressource
		for _, item := range regionenContents {

			if item.Edition == bandNummer {
				page.Relationships.Region = append(page.Relationships.Region, strings.Join([]string{ApiPageUrl, strings.TrimPrefix(item.Path, h.ContentDir)}, ""))
			}
		}

	}

	// print json
	printJson(w, &readPageResponse{Page: page})
}

// createPage creates a new page
func (h Handlers) CreatePage(w http.ResponseWriter, r *http.Request) {
	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// check that parent dir exists
	if fileExists(fp) || dirExists(fp) == false {
		errDirNotFound.Write(w)
		return
	}

	metastring := r.FormValue("page[meta]")
	if len(metastring) == 0 {
		errNoMeta.Write(w)
	}

	metadata := lidlib.Frontmatter{}
	err = json.Unmarshal([]byte(metastring), &metadata)
	if err != nil {
		errInvalidJson.Write(w)
		return
	}

	content := []byte(r.FormValue("page[content]"))

	page, err := h.Page.Create(fp, metadata, content)
	if err != nil {
		wrapError(err).Write(w)
		return
	}

	// trim content prefix from path
	page.Path = strings.TrimPrefix(page.Path, h.ContentDir)

	printJson(w, &createPageResponse{Page: page})
}

// updatePage writes page data to a file
func (h Handlers) UpdatePage(w http.ResponseWriter, r *http.Request) {
	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// check that existing page exists
	if dirExists(fp) || fileExists(fp) == false {
		errPageNotFound.Write(w)
		return
	}

	metastring := r.FormValue("page[meta]")
	if len(metastring) == 0 {
		errNoMeta.Write(w)
	}

	metadata := lidlib.Frontmatter{}
	err = json.Unmarshal([]byte(metastring), &metadata)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	content := []byte(r.FormValue("page[content]"))

	page, err := h.Page.Update(fp, metadata, content)
	if err != nil {
		wrapError(err).Write(w)
		return
	}

	// trim content prefix from path
	page.Path = strings.TrimPrefix(page.Path, h.ContentDir)

	printJson(w, &updatePageResponse{Page: page})
}

// TODO: make this more safe with captcha
// DeletePage deletes a page
func (h Handlers) DeletePage(w http.ResponseWriter, r *http.Request) {
	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// delete page
	if err = h.Page.Delete(fp); err != nil {
		errPageNotFound.Write(w)
		return
	}

	// don't need to send anything back
	w.WriteHeader(http.StatusNoContent)
}

/*
  Config
*/

// readConfig reads data from a config
func (h Handlers) ReadConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.Config.Parse()
	if err != nil {
		errNoConfig.Write(w)
		return
	}

	printJson(w, config)
}

// updateConfig writes json data to a config file
func (h Handlers) UpdateConfig(w http.ResponseWriter, r *http.Request) {

	// parse the config
	config := &lidlib.ConfigMap{}
	err := json.Unmarshal([]byte(r.FormValue("config")), config)
	if err != nil {
		errInvalidJson.Write(w)
		return
	}

	// save config
	if err := h.Config.Save(config); err != nil {
		wrapError(err).Write(w)
		return
	}

	// don't need to send anything back
	w.WriteHeader(http.StatusNoContent)
}

/*
  Assets
*/

type createAssetResponse struct {
	Asset *lidlib.Asset `json:"asset"`
}

// CreateAsset uploads a file into the assets directory
func (h Handlers) CreateAsset(w http.ResponseWriter, r *http.Request) {

	// get path to store file in
	dir, err := h.fixPathWithDir(mux.Vars(r)["path"], h.AssetsDir)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	// Check page exists [optional]

	// Remove .md extension
	ext := path.Ext(dir)
	dir = dir[0 : len(dir)-len(ext)]

	// Create folder structure in assets folder
	os.MkdirAll(dir, 0755)

	// Get file form request
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	// Sanitize file name
	filename := sanitize.Path(header.Filename)
	fp := path.Join(dir, filename)

	// Check file name doesn't already exist

	// TODO: save to path based on page name and sanitized file name
	out, err := os.Create(fp)
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing.")
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	asset := &lidlib.Asset{
		Name: filename,
		Path: dir,
	}

	asset.Resample()

	asset.Path = strings.TrimPrefix(asset.Path, h.AssetsDir)

	// TODO: print out proper status message
	printJson(w, &createAssetResponse{Asset: asset})

	// Write filename into page [optional]
}

/*
  Hugo
*/

func (h Handlers) PublishSite(w http.ResponseWriter, r *http.Request) {
	output, err := lidlib.RunHugo()
	if err != nil {
		wrapError(err).Write(w)
	}

	printJson(w, struct {
		Output string `json:"output"`
	}{
		Output: string(output),
	})
}

func (h Handlers) fixPathWithDir(p string, dir string) (string, error) {
	err := errors.New("invalid path")

	// join path with content folder
	fp := path.Join(dir, p)

	fmt.Println(fp)

	// check that path still starts with content dir
	if !strings.HasPrefix(fp, dir) {
		return fp, err
	}

	// check that path doesn't contain any ..
	if strings.Contains(fp, "..") {
		return fp, err
	}

	return fp, nil
}
