package main

import (
	backend "backend"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := backend.App{}
	app.Port = ":8080"
	app.Initialize()
	app.Run()
}
