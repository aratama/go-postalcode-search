package main

import (
	"log"
	"net/http"

	postalcodeSearch "example.com/postalcode-search"
)

func main() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		postalcodeSearch.PostalCodeSearch(w, req)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
