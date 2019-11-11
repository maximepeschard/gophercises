package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maximepeschard/gophercises/08_phone/gorm"
	"github.com/maximepeschard/gophercises/08_phone/sql"
	"github.com/maximepeschard/gophercises/08_phone/sqlx"
)

var usage = `Usage: phone [options...] {sql, sqlx, gorm}`

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "sql":
		sql.Main()
	case "sqlx":
		sqlx.Main()
	case "gorm":
		gorm.Main()
	default:
		flag.Usage()
		os.Exit(1)
	}
}
