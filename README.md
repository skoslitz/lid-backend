# lid-backend

## content api for lid online editor

### Windows Build
- Rename lid-backend._syso to *.syso
- Compile with `env GOOS=windows GOARCH=amd64 go build -o lid-backend.exe`
- NOTE: If you move lid-backend exec, always take config.toml along

### Endpoints:

#### Basepath: /api

#### Directories (Items in your repo `/content` folder)
- Read Dir: GET /dir/{dirname}/
- Read Region related content Filter: GET /regionen/{filename.md}/[themen exkursionen]
- Create Dir: POST /dir/{dirname}
- Update Dir: OPTIONS /dir/{dirname}
- Delete Dir: DELETE /dir/{dirname}

Example*: `OPTIONS dir/themen --form dir[name]="themen-neu"`

*: under usage of [http-prompt](https://github.com/eliangcs/http-prompt)

#### Pages
- Read Page: GET /page/{path}/
- Create Page: POST /page/{path}
- Update Page: OPTIONS /page/{path}
- Delete Page: DELETE /page/{path}

Example*: `OPTIONS page/themen/77_B_000.md --form page[meta]={\"title\": \"Beitragstitel\"} page[content]="text string"`

Example (REST client usage):
`POST page/themen/79_B_100-titel.md`
Raw payload:
`{ "page": { "id": "77_B_100-titel.md", "type": "themen", "links": {}, "attributes": {"metadata": {"title": "79_B_100-titel"}} }}`

### Development

## Installation:

- First, install [golang](https://golang.org/doc/install#install)
- Then: `go get github.com/skoslitz/lid-backend`
- Edit edit `config.toml` with required path infos
- Compile with `go build`
- Run with `./lid-backend`
