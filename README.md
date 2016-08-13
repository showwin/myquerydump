# myquerydump
`myquerydump` is like `mysqldump`, though with the very useful option `-q` (query).
`myquerydump` can dump with any complecated `SELECT` query.

## USAGE
```
$ myquerydump -u myuser -p mypassword -q "SELECT * FROM users INNER JOIN histories ON users.id = histories.user_id WHERE users.id > 10 ORDER BY users.updated_at" -t users_with_history mydatabase > users_with_history.dump

$ cat users_with_history.dump
> INSERT INTO `users_with_history` VALUES (foo),(bar)…
```
