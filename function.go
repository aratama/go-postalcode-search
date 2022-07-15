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

	// for i := 0; i < 100; i++ {
	// 	fmt.Printf("%v\n", kenallRecords[i])
	// }

	fmt.Printf("csv load time: %d msecs\n", time.Since(csvLoadStart).Milliseconds())

	functions.HTTP("PostalCodeSearch", PostalCodeSearch)
}

func PostalCodeSearch(w http.ResponseWriter, r *http.Request) {

	query := KatakanaToHiragana(HankakuKatakanaToKatakana(r.URL.Query().Get("q")))

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
	res.Q = query
	res.Limit = limit
	res.Page = page

	for i := 1; i < len(kenall); i++ {
		row := kenall[i]
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
