package xmysql

import "database/sql"

var Db *sql.DB

func init() {
	var err error

	if Db, err = sql.Open("mysql", "root:tantan@tcp(127.0.0.1:3306)/chexiang"); err != nil {
		panic(err)
	}
}
