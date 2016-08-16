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
  -t		Table name using when importing.
  -u		User for login if not current user.
  -add-delete-table
		Add DELETE before INSERT.
  -extended-insert
		Use multiple-row INSERT syntax that include 10000 VALUES lists.
		(Defaults to on; use -skip-extended-insert to disable.)
  -insert-ignore
		Insert rows with INSERT IGNORE.

$ myquerydump -u myuser -p -t users_with_history mydatabase "SELECT * FROM users INNER JOIN histories ON users.id = histories.user_id WHERE users.id > 10 ORDER BY users.updated_at" > myquery.dump

$ cat myquery.dump
INSERT INTO `users_with_history` VALUES (foo),(bar)…
```

## DETAIL
With `-add-delete-table` OPTION, empty the table before INSERT.
If `-t` OPTION is not provided, table name is parsed from SQL query (next string to first `FROM`).
```
$ myquerydump -add-delete-table mydatabase "SELECT * FROM users ORDER BY users.updated_at"
DELETE FROM `users`;
INSERT INTO `users` VALUES (…)
```
## NOT SUPPORT
* `CREATE TABLE`: Cannot specify the table schema for records selected by any query.
* database dump: Using `mysqldump` is better. `myquerydump` only support dumping with single SQL query.

## PERFORMANCE
Much slower than `mysqldump` 🙄, though still practical.

```
mysql> DESC users;
+------------+-------------+------+-----+---------+----------------+
| Field      | Type        | Null | Key | Default | Extra          |
+------------+-------------+------+-----+---------+----------------+
| id         | int(11)     | NO   | PRI | NULL    | auto_increment |
| first_name | varchar(64) | YES  |     | NULL    |                |
| last_name  | varchar(64) | YES  |     | NULL    |                |
| age        | int(11)     | YES  | MUL | NULL    |                |
| country    | varchar(64) | YES  |     | NULL    |                |
+------------+-------------+------+-----+---------+----------------+
mysql> SELECT COUNT(*) FROM users;
+----------+
| count(*) |
+----------+
|  3000000 |
+----------+

$ mysqldump mydatabase users > users.dump
=> 4 sec

$ myquerydump mydatabase "SELECT * FROM users" > users.dump
=> 78 sec
```

## ToDo
* `--fields-terminated-by=name` OPTION
* `--fields-enclosed-by=name` OPTION
* `--lines-terminated-by=name` OPTION
* `-S, --socket=name` OPTION
* `--ssl-mode=name` OPTION
* `-V, --version` OPTION
* Install from Homebrew

## LICENSE

[MIT](https://github.com/showwin/myquerydump/blob/master/LICENSE)
