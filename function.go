package postalcodeSearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type kenall = [][]string

var kennallArray kenall

type result struct {
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Hits  [][]string `json:"hits"`
	Time  int        `json:"time"`
}

func init() {
	loadStart := time.Now()

	dir := "./"
	functionTarget := os.Getenv("FUNCTION_TARGET")
	// fmt.Printf("FUNCTION_TARGET=%s\n", functionTarget)
	if functionTarget != "" {
		dir = "/workspace/serverless_function_source_code/"
	}

	file, err := ioutil.ReadFile(dir + "x-ken-all-hiragana.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &kennallArray)
	if err != nil {
		panic(err)
	}

	fmt.Printf("load time: %d msecs\n", time.Since(loadStart).Milliseconds())

	functions.HTTP("PostalCodeSearch", PostalCodeSearch)
}

func PostalCodeSearch(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = parsed
		}
	}
	limit = int(math.Min(float64(limit), 20))

	page := 0
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		parsed, err := strconv.Atoi(pageStr)
		if err == nil {
			page = parsed
		}
	}

	searchStart := time.Now()
	skip := limit * page

	var res result
	res.Limit = limit
	res.Page = page

	for i := 1; i < len(kennallArray); i++ {
		row := kennallArray[i]
		found := false
		for k := 0; k < len(row); k++ {
			value := row[k]
			if 0 <= strings.Index(value, query) {
				found = true
				break
			}
		}
		if found {
			if skip <= 0 {
				res.Hits = append(res.Hits, row)
				limit -= 1
				if limit == 0 {
					break
				}
			} else {
				skip -= 1
			}
		}
	}

	res.Time = int(time.Since(searchStart).Milliseconds())

	resultString, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, "error\n")
	} else {
		fmt.Fprintf(w, "%s\n", resultString)
	}

}
