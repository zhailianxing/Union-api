package share

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var AdvDb *sql.DB
var TestDb *sql.DB

func init() {
	var err error

	if AdvDb, err = sql.Open("mysql", "master:hyTTv587@tcp(rm-2zeb0rhm0o605a7o3.mysql.rds.aliyuncs.com:3306)/union"); err != nil {
		panic(err)
	}
	if TestDb, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/union"); err != nil {
		panic(err)
	}
}
