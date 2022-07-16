package function

import (
	"net/http"

	postalcodesearch "github.com/aratama/go-postalcode-search/postalcodesearch"
)

func PostalCodeSearch(w http.ResponseWriter, r *http.Request) {
	postalcodesearch.PostalCodeSearch(w, r)
}
