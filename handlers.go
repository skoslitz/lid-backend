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
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kennygrant/sanitize"
	"github.com/skoslitz/lid-backend/lidlib"
	"github.com/spf13/viper"
)

type Handlers struct {
	Config lidlib.ConfigManager
	Dir    lidlib.DirManager
	Page   lidlib.PageManager

	ContentDir string
	AssetsDir  string
	PreviewDir string
}

/*
  Directories
*/

type readDirResponse struct {
	Dir lidlib.Files `json:"data"`
}

type createDirResponse struct {
	Dir *lidlib.File `json:"data"`
}

type updateDirResponse struct {
	Dir *lidlib.File `json:"data"`
}

// readDir reads contents of a directory
func (h Handlers) ReadDir(w http.ResponseWriter, r *http.Request) {

	var ApiPageUrl = strings.Join([]string{"http://", r.Host, "/api/page/"}, "")
	var ApiUrl = strings.Join([]string{"http://", r.Host, "/api/"}, "")

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

		// TODO
		// Convert item.Path to raw string, then strings.Replace(item.Path, "\\", "/", -1)
		// or regexp.Replace \ with /
		// see also func ReadRegionRelationships

		item.Path = strings.TrimPrefix(item.Path, h.ContentDir)
		item.Self = strings.Join([]string{ApiPageUrl, item.Path}, "")
		// set type from path (content folder) to ember model names
		searchTerm := `(([\s\S]+?)[/]{1}([\s\S]+?)[.md])`
		re := regexp.MustCompile(searchTerm)
		reSlice := re.FindStringSubmatch(string(item.Path))

		if len(reSlice) > 0 {
			it := re.FindStringSubmatch(string(item.Path))[2]
			switch it {
			case "regionen":
				item.Type = "region-list"
			case "themen":
				item.Type = "topic-list"
			case "exkursionen":
				item.Type = "excursion-list"
			case "reihe":
				item.Type = "series"
			case "meta":
				item.Type = "meta"
			}

		}

		item.SetRelationship(ApiUrl)

	}

	printJson(w, &readDirResponse{Dir: contents})
}

// reads content of a directory and filters by region
func (h Handlers) ReadRegionRelationships(w http.ResponseWriter, r *http.Request) {

	ctype := mux.Vars(r)["type"]
	if containsContentType(ctype) {
		cid := strings.Split(mux.Vars(r)["id"], "-")[0]
		var cTypePath = strings.Join([]string{h.ContentDir, ctype}, "")
		var ApiPageUrl = strings.Join([]string{"http://", r.Host, "/api/page/"}, "")
		var ApiUrl = strings.Join([]string{"http://", r.Host, "/api/"}, "")

		// try and read contents of dir
		var contents lidlib.Files
		contents, err := h.Dir.Read(cTypePath)
		if err != nil {
			errDirNotFound.Write(w)
			return
		}

		var filteredContent lidlib.Files
		// trim content prefix
		for _, item := range contents {
			var cRegionId = strings.Split(item.Id, "_")[0]
			if cRegionId == cid {
				item.Path = strings.TrimPrefix(item.Path, h.ContentDir)
				item.Self = strings.Join([]string{ApiPageUrl, item.Path}, "")
				searchTerm := `(([\s\S]+?)[/]{1}([\s\S]+?)[.md])`
				re := regexp.MustCompile(searchTerm)
				it := re.FindStringSubmatch(string(item.Path))[2]
				switch it {
				case "regionen":
					item.Type = "region-list"
				case "themen":
					item.Type = "topic-list"
				case "exkursionen":
					item.Type = "excursion-list"
				case "reihe":
					item.Type = "series"
				case "meta":
					item.Type = "meta"
				}
				filteredContent = append(filteredContent, item)
				item.SetRelationship(ApiUrl)
			}

		}
		printJson(w, &readDirResponse{Dir: filteredContent})
	} else {
		errInvalidDir.Write(w)
		return
	}

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

// TODO: make this more safe with captcha or similar
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
	Page *lidlib.PageFile `json:"data"`
}

type createPageResponse struct {
	Page *lidlib.PageFile `json:"data"`
}

type updatePageResponse struct {
	Page *lidlib.PageFile `json:"data"`
}

// readPage reads page data
func (h Handlers) ReadPage(w http.ResponseWriter, r *http.Request) {
	var ApiUrl = strings.Join([]string{"http://", r.Host, "/api/"}, "")
	var ApiPageUrl = strings.Join([]string{"http://", r.Host, "/api/page/"}, "")

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
	page.Self = strings.Join([]string{ApiPageUrl, page.Path}, "")

	searchTerm := `(([\s\S]+?)[/]{1}([\s\S]+?)[.md])`
	re := regexp.MustCompile(searchTerm)
	reSlice := re.FindStringSubmatch(string(page.Path))

	if len(reSlice) > 0 {
		pt := re.FindStringSubmatch(string(page.Path))[2]
		switch pt {
		case "regionen":
			page.Type = "region"
		case "themen":
			page.Type = "topic"
		case "exkursionen":
			page.Type = "excursion"
		case "reihe":
			page.Type = "series"
		case "meta":
			page.Type = "meta"
		}

	}

	page.SetRelationship(ApiUrl)

	// Append content page related images info
	// -------------------------------------------

	var assets lidlib.Files
	var assetsFilePath string
	var assetsPath string
	assetsDir := viper.GetString("assetspath")

	// compose assetsFilePath from request path
	apiAssetsPath := strings.Join([]string{"http://", r.Host, "/assets"}, "")
	requestPath := mux.Vars(r)["path"]
	contentType := strings.Split(requestPath, "/")[0]
	fileName := strings.Split(requestPath, "/")[1]
	regionNr := strings.Split(fileName, "_")[0]
	hugoID := strings.Split(fileName, "-")[0]

	switch contentType {
	case "themen":
		assetsFilePath = strings.Join([]string{assetsDir, "img", regionNr, contentType, hugoID}, "/")
		assetsPath = strings.Join([]string{apiAssetsPath, "img", regionNr, contentType, hugoID}, "/")
	case "exkursionen":
		assetsFilePath = strings.Join([]string{assetsDir, "img", regionNr, contentType, hugoID}, "/")
		assetsPath = strings.Join([]string{apiAssetsPath, "img", regionNr, contentType, hugoID}, "/")
	case "regionen":
		regionNr := strings.Split(fileName, "-")[0]
		assetsFilePath = strings.Join([]string{assetsDir, "img", regionNr}, "/")
		assetsPath = strings.Join([]string{apiAssetsPath, "img", regionNr}, "/")
	default:
		assetsFilePath = ""
		assetsPath = ""
	}

	if assetsFilePath != "" {
		assets, err = h.Dir.Read(assetsFilePath)
		if err != nil {
			errInvalidAssetDir.Write(w)
			return
		}

		// write assetsDir content to page.Bilder slice
		for _, item := range assets {

			if !item.IsDir {
				bild := &lidlib.Bild{
					Src:      strings.Join([]string{assetsPath, item.Name}, "/"),
					Filename: item.Name,
				}

				page.Bilder = append(page.Bilder, bild)
			}
		}
	}
	// -------------------------------------------

	// print json pageresponse
	printJson(w, &readPageResponse{Page: page})
}

// createPage creates a new page
func (h Handlers) CreatePage(w http.ResponseWriter, r *http.Request) {
	// parse the incoming pageFile
	var pageFileJSON lidlib.PageFileJSON
	err := json.NewDecoder(r.Body).Decode(&pageFileJSON)

	fp, err := h.fixPathWithDir(mux.Vars(r)["path"], h.ContentDir)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	metadata := pageFileJSON.PageFile.Metadata
	content := []byte(pageFileJSON.PageFile.Content)

	page, err := h.Page.Create(fp, metadata, content)
	if err != nil {
		wrapError(err).Write(w)
		return
	}

	// trim content prefix from path
	page.Path = strings.TrimPrefix(page.Path, h.ContentDir)

	// prepare json resource type
	searchTerm := `(([\s\S]+?)[/]{1}([\s\S]+?)[.md])`
	re := regexp.MustCompile(searchTerm)
	reSlice := re.FindStringSubmatch(string(page.Path))

	// prepare page static img folder check
	hugoId := fmt.Sprint(page.Metadata["id"])
	regionNumber := fmt.Sprint(page.Metadata["bandnummer"])
	bandnummer := hugoId[:2]
	regionAssetfolder := h.AssetsDir + "img/" + regionNumber
	topicAssetfolder := h.AssetsDir + "img/" + bandnummer + "/themen/" + hugoId
	excursionAssetfolder := h.AssetsDir + "img/" + bandnummer + "/exkursionen/" + hugoId

	// set json resource type
	// check img static folder
	if len(reSlice) > 0 {
		pt := re.FindStringSubmatch(string(page.Path))[2]
		switch pt {
		case "regionen":
			page.Type = "regions"

			if dirExists(regionAssetfolder) || fileExists(regionAssetfolder) == false {
				os.MkdirAll(regionAssetfolder, 0755)
				fmt.Println("Anhangsverzeichnis wird erstellt ", regionAssetfolder)
			}
		case "themen":
			page.Type = "topics"

			if dirExists(topicAssetfolder) || fileExists(topicAssetfolder) == false {
				os.MkdirAll(topicAssetfolder, 0755)
				fmt.Println("Anhangsverzeichnis wird erstellt ", topicAssetfolder)
			}
		case "exkursionen":
			page.Type = "excursions"

			if dirExists(excursionAssetfolder) || fileExists(excursionAssetfolder) == false {
				os.MkdirAll(excursionAssetfolder, 0755)
				fmt.Println("Anhangsverzeichnis wird erstellt ", excursionAssetfolder)
			}
		case "reihe":
			page.Type = "series"
		case "meta":
			page.Type = "meta"
		}

	}

	printJson(w, &createPageResponse{Page: page})
}

// updatePage writes page data to a file
func (h Handlers) UpdatePage(w http.ResponseWriter, r *http.Request) {
	// parse the incoming pageFile
	var pageFileJSON lidlib.PageFileJSON
	err := json.NewDecoder(r.Body).Decode(&pageFileJSON)

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

	metadata := pageFileJSON.PageFile.Metadata

	content := []byte(pageFileJSON.PageFile.Content)

	page, err := h.Page.Update(fp, metadata, content)
	if err != nil {
		wrapError(err).Write(w)
		return
	}

	var ApiUrl = strings.Join([]string{"http://", r.Host, "/api/"}, "")
	var ApiPageUrl = strings.Join([]string{"http://", r.Host, "/api/page/"}, "")
	// trim content prefix from path
	page.Path = strings.TrimPrefix(page.Path, h.ContentDir)
	page.Self = strings.Join([]string{ApiPageUrl, page.Path}, "")

	searchTerm := `(([\s\S]+?)[/]{1}([\s\S]+?)[.md])`
	re := regexp.MustCompile(searchTerm)
	reSlice := re.FindStringSubmatch(string(page.Path))

	if len(reSlice) > 0 {
		pt := re.FindStringSubmatch(string(page.Path))[2]
		switch pt {
		case "regionen":
			page.Type = "region"
		case "themen":
			page.Type = "topic"
		case "exkursionen":
			page.Type = "excursion"
		case "reihe":
			page.Type = "series"
		case "meta":
			page.Type = "meta"
		}

	}

	page.SetRelationship(ApiUrl)

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

	// Get file form request
	file, fileheader, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	// Sanitize file name
	filename := sanitize.Path(fileheader.Filename)
	fp := path.Join(dir, filename)

	// Check file name doesn't already exist

	// TODO: save to path based on page type and id
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
	repoPath := viper.GetString("repopath")
	webFolder := viper.GetString("webfolderpath")

	output, err := lidlib.RunHugo(repoPath, webFolder)
	if err != nil {
		wrapError(err).Write(w)
	}

	printJson(w, struct {
		Output string `json:"output"`
	}{
		Output: string(output),
	})
}

func (h Handlers) PreviewSite(w http.ResponseWriter, r *http.Request) {
	baseUrlPrefix := strings.Join([]string{"http://", r.Host}, "")
	repoPath := viper.GetString("repopath")

	output, err := lidlib.RunHugoPreview(baseUrlPrefix, repoPath)
	if err != nil {
		fmt.Println("hugo Stdout should appear as message")
		wrapError(err).Write(w)
		fmt.Println(output)
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
