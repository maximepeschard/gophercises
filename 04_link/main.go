package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/maximepeschard/gophercises/04_link/link"
)

var commands = [...]string{"html", "url"}

var usage = `Usage:
  links html <HTML file>
  links url <url>
`

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()
	if flag.NArg() != 2 || !validCommand(flag.Arg(0)) {
		printUsageAndExit()
	}
	command := flag.Arg(0)

	var reader io.Reader
	var err error
	if command == "html" {
		reader, err = os.Open(flag.Arg(1))
		if err != nil {
			printErrorAndExit(err)
		}
	} else {
		response, err := http.Get(flag.Arg(1))
		if err != nil {
			printErrorAndExit(err)
		}
		reader = response.Body
	}

	links, err := link.Parse(reader)
	if err != nil {
		printErrorAndExit(err)
	}

	for _, link := range links {
		fmt.Printf("\"%s\" -> %s\n", link.Text, link.Href)
	}
}

func validCommand(cmd string) bool {
	for _, c := range commands {
		if cmd == c {
			return true
		}
	}
	return false
}

func printUsageAndExit() {
	flag.Usage()
	os.Exit(1)
}

func printErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
