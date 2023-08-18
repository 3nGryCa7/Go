package main

import (
	"example.com/backend"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	a := backend.App{}
	a.Port = ":8080"
	a.Initialize()
	a.Run()
}
