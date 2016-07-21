# lid-backend - content api for lid

## Installation:

- First, install [golang](https://golang.org/doc/install#install)
- Then: `go get github.com/skoslitz/lid-backend`
- Edit `main.go` on line 20 to specify content repo path
- Compile with `go build`
- Run with `./lid-backend`

### Endpoints:

#### Basepath: /api

#### Directories (Items in your repo `/content` folder)
- Read Dir: GET /dir/{dirname}/
- Read Region related content Filter: GET /regionen/{filename.md}/[themen exkursionen]
- Create Dir: POST /dir/{dirname}
- Update Dir: PUT /dir/{dirname}
- Delete Dir: DELETE /dir/{dirname}

Example*: `PUT dir/themen --form dir[name]="themen-neu"`

#### Pages
- Read Page: GET /page/{path}/
- Create Page: POST /page/{path}
- Update Page: PUT /page/{path}
- Delete Page: DELETE /page/{path}

Example*: `PUT page/themen/77_B_000.md --form page[meta]={\"title\": \"Beitragstitel\"} page[content]="text string"`

Example (REST client usage):
POST `page/themen/79_B_100-titel.md`
Raw payload: 
`{ "page": { "id": "77_B_100-titel.md", "type": "themen", "links": {}, "attributes": {"metadata": {"title": "79_B_100-titel"}} }}`

#### Config
- Read Config: GET /config
- Update Config: PUT /config


*: under usage of [http-prompt](https://github.com/eliangcs/http-prompt)
