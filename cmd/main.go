package main

import (
	"log"
	"net/http"

	postalcodeSearch "example.com/postalcode-search/postalcodesearch"
)

func main() {
	http.HandleFunc("/", postalcodeSearch.PostalCodeSearch)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
