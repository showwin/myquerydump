package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"

	"github.com/howeyc/gopass"
)

var db *sql.DB

var (
	user           = flag.String("u", os.Getenv("USER"), "")
	passFlg        = flag.Bool("p", false, "")
	tableName      = flag.String("t", "", "")
	deleteTableFlg = flag.Bool("add-delete-table", false, "")
)

var usage = `Usage: myquerydump [options...] <database> <query>
Options:
  -u		User for login if not current user.
  -p		Password to use when connecting to server. It's solicited on the tty.
  -t		Table name using when importing.
  -add-delete-table
		Add DELETE FROM before INSERT.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()
	if flag.NArg() < 2 {
		usageAndExit("")
	}

	var password string
	if *passFlg {
		fmt.Print("Enter password: ")
		input, _ := gopass.GetPasswd()
		password = string(input)
	}

	database := flag.Args()[0]
	query := flag.Args()[1]
	db, err := sql.Open("mysql", *user+":"+password+"@/"+database)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if *tableName == "" {
		*tableName = parseTableName(query)
	}

	rows, err := db.Query(query)
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

	var sql string
	if *deleteTableFlg {
		sql = "DELETE FROM `" + *tableName + "`;\n"
	}
	sql = sql + "INSERT INTO `" + *tableName + "` VALUES "
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

func parseTableName(query string) string {
	r := regexp.MustCompile(`.+(FROM|from)\s(\w+)($|\s.*)`)
	return r.ReplaceAllString(query, "$2")
}
