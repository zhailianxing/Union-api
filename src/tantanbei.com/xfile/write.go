package xfile

import (
	"io/ioutil"

	"tantanbei.com/xjson"
)

func WriteFile(filename string, data []byte) {
	err := ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		panic(err)
	}
}

func WriteJsonFile(filename string, v interface{}) {
	bs := xjson.Encode(v)

	WriteFile(filename, bs)
}
