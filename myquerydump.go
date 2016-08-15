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
	host           = flag.String("h", "localhost", "")
	passFlg        = flag.Bool("p", false, "")
	port           = flag.String("P", "3306", "")
	tableName      = flag.String("t", "", "")
	user           = flag.String("u", os.Getenv("USER"), "")
	deleteTableFlg = flag.Bool("add-delete-table", false, "")
)

var usage = `Usage: myquerydump [options...] <database> <query>
Options:
  -h		Connect to host.
  -p		Password to use when connecting to server. It's solicited on the tty.
  -P		Port number to use for connection.
  -t		Table name using when importing.
  -u		User for login if not current user.
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
		password = inputPassword()
	}

	database := flag.Args()[0]
	query := flag.Args()[1]
	db, err := sql.Open("mysql", *user+":"+password+"@tcp("+*host+":"+*port+")/"+database)
	if err != nil {
		errorAndExit(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		errorAndExit(err)
	}

	columns, err := rows.Columns()
	if err != nil {
		errorAndExit(err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if *tableName == "" {
		*tableName = parseTableName(query)
	}
	var sql string
	if *deleteTableFlg {
		sql = "DELETE FROM `" + *tableName + "`;\n"
	}
	sql = sql + "INSERT INTO `" + *tableName + "` VALUES "
	for rows.Next() {
		sqlTmp := "("
		err = rows.Scan(scanArgs...)
		if err != nil {
			errorAndExit(err)
		}

		var value string
		for _, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			sqlTmp = sqlTmp + "'" + value + "',"
		}
		sql = sql + sqlTmp[0:len(sqlTmp)-1] + "),"
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

func errorAndExit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func parseTableName(query string) string {
	r := regexp.MustCompile(`.+(FROM|from)\s(\w+)($|\s.*)`)
	return r.ReplaceAllString(query, "$2")
}

func inputPassword() string {
	fmt.Print("Enter password: ")
	str, _ := gopass.GetPasswd()
	return string(str)
}
