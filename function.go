package function

import (
	"net/http"

	postalcodeSearch "example.com/postalcode-search/postalcodesearch"
)

func PostalCodeSearch(w http.ResponseWriter, r *http.Request) {
	postalcodeSearch.PostalCodeSearch(w, r)
}
