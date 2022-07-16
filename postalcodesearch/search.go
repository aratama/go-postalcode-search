package postalcodeSearch

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type PostalCodeSearchResult struct {
	Page       int        `json:"page,omitempty"`
	Q          string     `json:"q,omitempty"`
	PostalCode string     `json:"postalcode,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Hits       [][]string `json:"hits"`
	Time       int        `json:"time"`
}

var kenall [][]string

func getIntParam(r *http.Request, name string, defaultValue int, maxValue int) int {
	param := defaultValue
	paramStr := r.URL.Query().Get(name)
	if paramStr != "" {
		parsed, err := strconv.Atoi(paramStr)
		if err == nil {
			param = parsed
		}
	}
	return int(math.Min(float64(param), float64(maxValue)))
}

func init() {
	// Get CSV file path
	// Note the static files are located at "/workspace/serverless_function_source_code/" in Cloud Functions
	// and `go test` also run tests on their own directory.
	// SO relative path from current working directory will not work
	_, b, _, _ := runtime.Caller(0)
	csvFilePath := filepath.Join(filepath.Dir(b), "../x-ken-all.csv")

	csvLoadStart := time.Now()

	fmt.Printf("Loading CSV from %s...\n", csvFilePath)

	shiftJisCSVFileHandle, err := os.Open(csvFilePath)
	if err != nil {
		panic(err)
	}

	// encode shift-jis to UTF8
	utf8CSVFileReader := transform.NewReader(shiftJisCSVFileHandle, japanese.ShiftJIS.NewDecoder())

	csvReader := csv.NewReader(utf8CSVFileReader)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	kenall = records

	for _, row := range kenall {
		row[3] = KatakanaToHiragana(HankakuKatakanaToKatakana(row[3]))
		row[4] = KatakanaToHiragana(HankakuKatakanaToKatakana(row[4]))
		row[5] = KatakanaToHiragana(HankakuKatakanaToKatakana(row[5]))
	}

	fmt.Printf("csv load time: %d msecs\n", time.Since(csvLoadStart).Milliseconds())
}

func searchInRow(row []string, query string) bool {
	for k := 3; k <= 8; k++ {
		value := row[k]
		if 0 <= strings.Index(value, query) {
			return true

		}
	}
	return false
}

func postalCodeSearchByQuery(w http.ResponseWriter, r *http.Request) PostalCodeSearchResult {

	query := KatakanaToHiragana(HankakuKatakanaToKatakana(r.URL.Query().Get("q")))

	limit := getIntParam(r, "lmit", 10, 20)

	page := getIntParam(r, "page", 0, 100)

	skip := limit * page

	res := PostalCodeSearchResult{
		Q:     query,
		Limit: limit,
		Page:  page,
	}

	for _, row := range kenall {
		if searchInRow(row, query) {
			if skip <= 0 {
				res.Hits = append(res.Hits, row)
				if limit <= len(res.Hits) {
					return res
				}
			} else {
				skip -= 1
			}
		}
	}

	return res

}

func postalCodeSearchByPostalCode(w http.ResponseWriter, r *http.Request, postalCodeWithHyphen string) PostalCodeSearchResult {

	postalCode := strings.Replace(postalCodeWithHyphen, "-", "", -1)

	limit := getIntParam(r, "lmit", 10, 20)

	var res PostalCodeSearchResult = PostalCodeSearchResult{
		PostalCode: postalCode,
		Limit:      limit,
	}

	for _, row := range kenall {
		if 0 <= strings.Index(row[2], postalCode) {
			res.Hits = append(res.Hits, row)
			if limit <= len(res.Hits) {
				break
			}
		}
	}

	return res

}

func PostalCodeSearch(w http.ResponseWriter, r *http.Request) {

	postalcode := r.URL.Query().Get("postalcode")

	searchStart := time.Now()

	var res PostalCodeSearchResult

	if postalcode == "" {
		res = postalCodeSearchByQuery(w, r)
	} else {
		res = postalCodeSearchByPostalCode(w, r, postalcode)
	}

	res.Time = int(time.Since(searchStart).Milliseconds())

	resultString, err := json.Marshal(res)

	if err != nil {
		fmt.Fprintf(w, "error\n")
	} else {
		fmt.Fprintf(w, "%s\n", resultString)
	}
}
