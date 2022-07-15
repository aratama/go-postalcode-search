package postalcodeSearch

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type result struct {
	Page  int        `json:"page"`
	Q     string     `json:"q"`
	Limit int        `json:"limit"`
	Hits  [][]string `json:"hits"`
	Time  int        `json:"time"`
}

var kenall [][]string

func init() {

	dir := "./"
	functionTarget := os.Getenv("FUNCTION_TARGET")
	// fmt.Printf("FUNCTION_TARGET=%s\n", functionTarget)
	if functionTarget != "" {
		dir = "/workspace/serverless_function_source_code/"
	}

	csvLoadStart := time.Now()

	shiftJisCSVFileHandle, err := os.Open(dir + "x-ken-all.csv")
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

	functions.HTTP("PostalCodeSearch", PostalCodeSearch)
}

func intParam(r *http.Request, name string, defaultValue int, maxValue int) int {
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

func postalCodeSearchByQuery(w http.ResponseWriter, r *http.Request) result {

	query := KatakanaToHiragana(HankakuKatakanaToKatakana(r.URL.Query().Get("q")))

	limit := intParam(r, "lmit", 10, 20)

	page := intParam(r, "page", 0, 100)

	skip := limit * page

	var res result = result{
		Q:     query,
		Limit: limit,
		Page:  page,
	}

	for _, row := range kenall {
		found := false

		for k := 3; k <= 8; k++ {
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

	return res

}

func postalCodeSearchByPostalCode(w http.ResponseWriter, r *http.Request, postalCode string) result {

	limit := intParam(r, "lmit", 10, 20)

	var res result = result{
		Q:     postalCode,
		Limit: limit,
		Page:  0,
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

	var res result

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
