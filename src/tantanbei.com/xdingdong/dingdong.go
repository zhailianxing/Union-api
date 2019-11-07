package xdingdong

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"errors"
	log "github.com/alecthomas/log4go"
	"tantanbei.com/xjson"
)

type Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result string `json:"result"`
}

func Send(mobile string, code string) error {
	if len(code) < 4 || len(code) > 6 {
		return ERROR_INVALID
	}

	for _, element := range code {
		if element < '0' || element > '9' {
			return ERROR_INVALID
		}
	}

	content := fmt.Sprintf(CONTENT, code)

	data := url.Values{"apikey": {APIKEY}, "mobile": {mobile}, "content": {content}}

	return httpsPostForm(URL, data)
}

func httpsPostForm(url string, data url.Values) error {
	resp, err := http.PostForm(url, data)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	response := &Response{}
	xjson.Decode(body, &response)
	if response.Code == 1 {
		return nil
	} else {
		log.Error(response)
		return ERROR_SEND_FAILED
	}
}

func SendTZ(mobile string, content string) error {
	if len(mobile) != 11 || content == "" {
		return errors.New("bad data")
	}

	data := url.Values{"apikey": {APIKEY}, "mobile": {mobile}, "content": {content}}

	return httpsPostForm(URL_TZ, data)
}
