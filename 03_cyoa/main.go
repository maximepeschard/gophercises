package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/maximepeschard/gophercises/03_cyoa/story"
)

var usage = `Usage: cyoa [options...] <story file>

Options:
  -templates  The directory containing the templates (default: tmpl).
`

type storyHandler struct {
	stry *story.Story
	tmpl *template.Template
}

func (sh storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arcName := strings.TrimPrefix(r.URL.Path, "/")
	if arcName == "" {
		arcName = "intro"
	}

	arc, prs := sh.stry.Arcs[arcName]
	if !prs {
		sh.tmpl.ExecuteTemplate(w, "404", "intro")
	} else {
		sh.tmpl.ExecuteTemplate(w, "arc", arc)
	}

}

func main() {
	templatesDir := flag.String("templates", "tmpl", "the directory containing the templates")
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()
	if flag.NArg() != 1 {
		printUsageAndExit()
	}

	storyFilename := flag.Arg(0)
	storyData, err := ioutil.ReadFile(storyFilename)
	if err != nil {
		printErrorAndExit(err)
	}

	story, err := story.NewStory(storyFilename).ParseJSON(storyData)
	if err != nil {
		printErrorAndExit(err)
	}

	templatesPattern := filepath.Join(*templatesDir, "*.tmpl")
	tmpl := template.Must(template.ParseGlob(templatesPattern))

	sh := storyHandler{story, tmpl}
	http.ListenAndServe(":8080", sh)

}

func printUsageAndExit() {
	flag.Usage()
	os.Exit(1)
}

func printErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
