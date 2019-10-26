package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/maximepeschard/gophercises/02_urlshort/urlshort"
)

func readFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	return content, err
}

func main() {
	yamlFilename := flag.String("yaml", "", "provide URL mappings using a YAML file")
	flag.Parse()

	// YAML configuration switch
	var yaml []byte
	if *yamlFilename != "" {
		content, err := ioutil.ReadFile(*yamlFilename)
		if err != nil {
			log.Fatal(err)
		}
		yaml = content
	} else {
		yamlString := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		yaml = []byte(yamlString)
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Build the JSONHandler using the yamlHandler as the
	// fallback
	jsn := `	
[
	{"path": "/hn", "url": "https://news.ycombinator.com"},
	{"path": "/theverge", "url": "https://theverge.com"}
]
`
	jsonHandler, err := urlshort.JSONHandler([]byte(jsn), yamlHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Build the RedisHandler using the jsonHandler as the
	// fallback
	redisHandler := urlshort.RedisHandler("localhost", 6379, 0, "urls", jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", redisHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
