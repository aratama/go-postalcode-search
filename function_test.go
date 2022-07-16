package postalcodeSearch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostalCodeSearch(t *testing.T) {

	t.Run("PostalCodeSearch", func(t *testing.T) {

		reqBody := bytes.NewBufferString("request body")
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/?q=長野", reqBody)
		recorder := httptest.NewRecorder()

		PostalCodeSearch(recorder, req)

		result := recorder.Result()
		defer result.Body.Close()
		if result.Status != "200 OK" {
			t.Errorf("expected = %v, want %v", "200 OK", recorder.Result().Status)
		}

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Errorf("error: %v", err)
		}

		var resJson PostalCodeSearchResult
		err = json.Unmarshal(body, &resJson)
		if err != nil {
			t.Errorf("json Unmarshal error: %v", err)
		}

		if len(resJson.Hits) != 10 {
			t.Errorf("not fouund")
		}

		if resJson.Q != "長野" {
			t.Errorf("q not match, expected 長野, found %v", resJson.Q)
		}

	})

}
