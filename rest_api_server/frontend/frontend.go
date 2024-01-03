package frontend

import (
	"backend"

	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	a := backend.App{}
	a.Port = ":8080"
	a.Initialize()
	a.Run()
}
