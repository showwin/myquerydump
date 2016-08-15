# myquerydump
`myquerydump` dumps records from MySQL with any complicated `SELECT` query.  
`myquerydump` is NOT all-powerful. Recommend to use this or `mysqldump` depending on a purpose.

## USAGE
```
$ myquerydump --help
Usage: myquerydump [options...] <database> <query>
Options:
  -h		Connect to host.
  -p		Password to use when connecting to server. It's solicited on the tty.
  -P		Port number to use for connection.
  -t		Table name using when importing. Default value: after first FROM.
  -u		User for login if not current user.
  -add-delete-table
		Add DELETE FROM before INSERT.

$ myquerydump -u myuser -p -t users_with_history mydatabase "SELECT * FROM users INNER JOIN histories ON users.id = histories.user_id WHERE users.id > 10 ORDER BY users.updated_at" > myquery.dump

$ cat myquery.dump
INSERT INTO `users_with_history` VALUES (foo),(bar)…
```

## DETAIL
With `-add-delete-table` OPTION, empty the table before INSERT.
```
$ myquerydump -add-delete-table mydatabase "SELECT * FROM users ORDER BY users.updated_at"
DELETE FROM `users`;
INSERT INTO `users` VALUES (…)
```
## NOT SUPPORT
* `CREATE TABLE`: Cannot specify the table schema for records selected by any query.
* database dump: Using `mysqldump` is better. `myquerydump` only support dumping with single SQL query.

## ToDo
* `--extended-insert` OPTION (by default)
* `--skip-extended-insert` OPTION
* `--fields-terminated-by=name` OPTION
* `--fields-enclosed-by=name` OPTION
* `--insert-ignore` OPTION
* `--lines-terminated-by=name` OPTION
* `--single-transaction` OPTION
* `-S, --socket=name` OPTION
* `--ssl-mode=name` OPTION
* `-V, --version` OPTION
* Install from Homebrew
* Performance measurement (vs mysqldump)
