package lidlib

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

type Bild struct {
	Src      string `json:"src"`
	Filename string `json:"filename"`
}

// Frontmatter stores encodeable data
type Frontmatter map[string]interface{}

type Attribute struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is-dir"`
	Size     int64       `json:"size"`
	ModTime  string      `json:"mod-time"`
	Metadata Frontmatter `json:"metadata,omitempty"`
	Content  string      `json:"content,omitempty"`
	Bilder   []*Bild     `json:"images,omitempty"`
}

type Relationship struct {
	Region    `json:"region"`
	Thema     `json:"themen"`
	Exkursion `json:"exkursionen"`
}

// File represents a file within directory context
type File struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Link         `json:"links"`
	Attribute    `json:"attributes"`
	Relationship `json:"relationships"`
}

// Page represents a markdown file
type PageFile struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Link         `json:"links"`
	Attribute    `json:"attributes"`
	Relationship `json:"relationships"`
}

type PageFileJSON struct {
	PageFile `json:"page"`
}

type Page struct{}

type Files []*File

const TOML = '+'
const YAML = '-'
