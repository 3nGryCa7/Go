package main

import (
	"example.com/rest_api_server/backend"
)

func main() {
	a := backend.App{}
	a.Port = ":8080"
	a.Initailize()
	a.Run()
}
