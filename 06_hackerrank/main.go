package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var usage = `Usage:
  hackerrank camelcase <string>
  hackerrank cesarcipher <shift> <string>
`

var upperCaseRegexp = regexp.MustCompile(`[A-Z]`)

func camelCase(s string) int {
	return len(upperCaseRegexp.FindAllStringIndex(s, -1)) + 1
}

func cesarCipher(shift int32, s string) string {
	rot := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+shift)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+shift)%26
		default:
			return r
		}
	}

	return strings.Map(rot, s)
}

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	command := flag.Arg(0)
	switch command {
	case "camelcase":
		if flag.NArg() != 2 {
			printUsageAndExit()
		}
		fmt.Println(camelCase(flag.Arg(1)))
	case "cesarcipher":
		if flag.NArg() != 3 {
			printUsageAndExit()
		}
		shift, err := strconv.ParseInt(flag.Arg(1), 10, 32)
		if err != nil {
			printErrorAndExit(err)
		}
		fmt.Println(cesarCipher(int32(shift), flag.Arg(2)))
	default:
		printUsageAndExit()
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
