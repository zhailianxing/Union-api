package xdingdong

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

//func TestSend(t *testing.T) {
//	err := Send("18117541072", "123456")
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func TestSendTZ(t *testing.T) {
	db, err := sql.Open("mysql", "root:tantan@tcp(127.0.0.1:3306)/chexiang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query("SELECT phone FROM user")
	if err != nil {
		panic(err)
	}

	for result.Next() {
		phone := ""
		result.Scan(&phone)
		err := SendTZ(phone, "【拍牌宝】值此新春佳节，拍牌宝团队为您及家人献上最诚挚的祝福：愿您在新的一年里，阖家欢乐，平安健康！退订回T")
		if err != nil {
			fmt.Println(err, phone)
			continue
		}
	}
}
