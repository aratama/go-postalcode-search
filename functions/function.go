package postalcodeSearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type kenall = [][]string

var kennallArray kenall

func init() {
	loadStart := time.Now()

	dir := "./"
	functionTarget := os.Getenv("FUNCTION_TARGET")
	fmt.Printf("FUNCTION_TARGET=%s\n", functionTarget)
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

	functions.HTTP("HelloWorld", HelloWorld)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	searchStart := time.Now()

	limit := 20

	for i := 0; i < len(kennallArray); i++ {
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
			fmt.Fprintf(w, "%+v\n", row)
			limit -= 1
			if limit == 0 {
				break
			}
		}
	}

	fmt.Printf("search time: %d msecs\n", time.Since(searchStart).Milliseconds())
}
