package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET request received")
}

func postRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST request received")
}

func deleteRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DELETE request received")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getRequest).Methods("GET")
	r.HandleFunc("/", postRequest).Methods("POST")
	r.HandleFunc("/", deleteRequest).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
