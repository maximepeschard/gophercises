package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/maximepeschard/gophercises/05_sitemap/sitemap"
)

var urlWithScheme = regexp.MustCompile(`^https?://.*$`)

var usage = `Usage: sitemap [options...] <url>

Options:
  -depth   Maximum number of links to follow (default: -1 for no maximum).
  -output  Name of a file to write the sitemap to (default: print to stdout).
`

func main() {
	depth := flag.Int("depth", -1, "maximum number of links to follow (default: -1 for no maximum)")
	output := flag.String("output", "", "name of a file to write the sitemap to (default: print to stdout)")
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()
	if flag.NArg() != 1 {
		printUsageAndExit()
	}

	siteURL, err := normalizeURL(flag.Arg(0))
	if err != nil {
		printErrorAndExit(err)
	}

	outputWriter := os.Stdout
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			printErrorAndExit(err)
		}
		defer closeFile(file)
		outputWriter = file
	}

	siteMap := sitemap.New(siteURL).Build(*depth)
	err = siteMap.WriteXML(outputWriter)
	if err != nil {
		printErrorAndExit(err)
	}
}

func normalizeURL(rawURL string) (*url.URL, error) {
	// if rawURL does not start with "http[s]://", we add it
	if !urlWithScheme.MatchString(rawURL) {
		rawURL = "https://" + rawURL
	}

	return url.Parse(rawURL)
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		printErrorAndExit(err)
	}
}

func printUsageAndExit() {
	flag.Usage()
	os.Exit(1)
}

func printErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
