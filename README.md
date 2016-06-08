# lid-backend - content api for lid

## Installation:

- First, install [golang](https://golang.org/doc/install#install)
- Then: `go get github.com/skoslitz/lid-backend`
- Edit `main.go` on line 20 to specify content repo path
- Compile with `go build`
- Run with `./lid-backend`

### Endpoints:

#### Basepath: /api

#### Directories
- Read Dir: GET /dir/{dirname}/
- Read Dir | Edition Filter: GET /dir/{dirname}/{editionNumber}
- Create Dir: POST /dir/{dirname}
- Update Dir: PUT /dir/{dirname}
- Delete Dir: DELETE /dir/{dirname}

#### Pages
- Read Page: GET /page/{path}/
- Create Page: POST /page/{path}
- Update Page: PUT /page/{path}
- Delete Page: DELETE /page/{path}

#### Config
- Read Config: GET /config
- Update Config: PUT /config
