package charge

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	log "github.com/alecthomas/log4go"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"strconv"
	"tantanbei.com/xjson"
	"tantanbei.com/xjuhe"
	"time"
)

type ResultPacket struct {
	ErrorCode int         `json:"error_code"`
	Reason    interface{} `json:"reason"`
	Result    interface{} `json:"result"`
}

func TopUp(phoneNum string, cardNum string) bool {

	switch cardNum {
	case CARD_NUM_10:
	case CARD_NUM_20:
	case CARD_NUM_30:
	case CARD_NUM_50:
	case CARD_NUM_100:
	case CARD_NUM_300:
	default:
		return false
	}

	orderId := "HFCZ" + strconv.Itoa(int(time.Now().UnixNano()))

	md5Byte := md5.Sum([]byte(xjuhe.OPEN_ID + APP_KEY + phoneNum + cardNum + orderId))

	sign := hex.EncodeToString(md5Byte[:])

	log.Debug(sign)

	response, err := http.Get(TOP_UP_URL + "?phoneno=" + phoneNum + "&cardnum=" + cardNum + "&orderid=" + orderId + "&key=" + APP_KEY + "&sign=" + sign)
	if err != nil {
		log.Error(err)
		return false
	}

	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return false
	}

	result := &ResultPacket{}
	log.Debug(string(bs))
	xjson.Decode(bs, &result)

	//check the result data, it maybe a null object
	if result.Result == nil {
		result.Result = ""
	}

	//check the reason data, it maybe a null object
	if result.Reason == nil {
		result.Reason = ""
	}

	//add the charge log
	AddRecord(phoneNum, cardNum, orderId, fmt.Sprintf("%v", result.Reason), fmt.Sprintf("%v", result.Result), result.ErrorCode)

	if result.ErrorCode == 0 {
		return true
	} else {
		log.Error(result)
		return false
	}
}

func AddRecord(phoneNum, cardNum, orderId, reason, result string, errorCode int) {
	db, err := sql.Open("mysql", "root:tantan@tcp(127.0.0.1:3306)/chexiang")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec(ADD_RECORD_SQL, phoneNum, cardNum, orderId, errorCode, reason, result)
	if err != nil {
		panic(err)
	}
}
