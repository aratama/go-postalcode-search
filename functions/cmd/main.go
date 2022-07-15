package main

import (
	"log"
	"net/http"

	postalcodeSearch "example.com/hello"
)

func main() {

	// fmt.Printf("%+v\n", res)

	handler := func(w http.ResponseWriter, req *http.Request) {
		postalcodeSearch.HelloWorld(w, req)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
