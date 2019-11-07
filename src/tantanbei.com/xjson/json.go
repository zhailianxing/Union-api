package xjson

import "tantanbei.com/json"

func Encode(v interface{}) (bs []byte) {

	var err error

	if bs, err = json.Marshal(v); err != nil {
		panic(err)
	}
	return
}

func Decode(data []byte, v interface{}) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}
