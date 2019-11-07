package api

import (
	"github.com/gocraft/web"
	"io/ioutil"
)

const (
	//0:mean failed, 1:mean succeed
	HEADER_SUCCEED = "succeed"
	FAILED         = "0"
	SUCCEED        = "1"

	HEADER_TOKEN = "t-tokenid"

	HEADER_ERROR         = "t-error-code"
	ERR_CODE_AUTH_FAILED = "0"
	ERR_CODE_FAIL_DATA   = "1"

	//one day
	ONE_DAY  = 60 * 60 * 24
	ONE_HOUR = 60 * 60
	//session expire time is one hour
	SESSION_EXPIRE_TIME = ONE_HOUR
)

var (
	EMPTY_BODY []byte = []byte("")
)

func Fail(w web.ResponseWriter) {
	w.Header().Add(HEADER_SUCCEED, FAILED)
	w.Write(EMPTY_BODY)
}

func Success(w web.ResponseWriter) {
	w.Header().Add(HEADER_SUCCEED, SUCCEED)
}

func FailWithErrCode(w web.ResponseWriter, errCode string) {
	w.Header().Add(HEADER_SUCCEED, FAILED)
	w.Header().Add(HEADER_ERROR, errCode)
	w.Write(EMPTY_BODY)
}

func GetFormFile(req *web.Request, key string) ([]byte, error) {
	file, _, err := req.FormFile(key)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}
