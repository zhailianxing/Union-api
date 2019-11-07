package xdb

import (
	"database/sql"
)

//the arg must be string
func SelectOne(db *sql.DB, sql string, arg ...interface{}) (*sql.Rows, bool) {
	result, err := db.Query(sql, arg...)
	if err != nil {
		return nil, false
	}

	if result.Next() {
		return result, true
	} else {
		return nil, false
	}
}
