package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var (
	user      = flag.String("u", "", "")
	password  = flag.String("p", "", "")
	query     = flag.String("q", "", "")
	tableName = flag.String("t", "", "")
)

var usage = `Usage: myquerydump [options...] <database>
Options:
  -u  User for login.
  -p  Password to use when connecting to server.
  -q  SQL Query to execute for dumping.
  -t  Table name using when importing.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()
	if flag.NArg() < 1 {
		usageAndExit("")
	}

	database := flag.Args()[0]
	db, err := sql.Open("mysql", *user+":"+*password+"@/"+database)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query(*query)
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	sql := "INSERT INTO `" + *tableName + "` VALUES "
	for rows.Next() {
		sqlRecord := "("
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for _, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			sqlRecord = sqlRecord + "'" + value + "',"
		}
		sql = sql + sqlRecord[0:len(sqlRecord)-1] + "),"
	}
	fmt.Println(sql[0:len(sql)-1] + ";")
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
