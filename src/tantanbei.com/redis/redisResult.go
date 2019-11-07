package redis

import (
	"fmt"
)

//when initial, the RedisResult is empty
type RedisResult struct {
	resultType int8

	resultBytes  []byte
	resultInt    int
	resultString string
	resultError  error //this error is redis reply
	resultArray  []*RedisResult
	resultOk     bool
	resultPong   bool
}

//Get the RedisResult type id
//  ERROR  = -1
//  EMPTY  = 0
//  INT    = 1
//  STRING = 2
//  BYTES  = 3
//  ARRAY  = 4
//  OK     = 5
//  PONG   = 6
func (r *RedisResult) GetTypeId() int8 {
	return r.resultType
}

//Get the RedisResult type
func (r *RedisResult) GetType() string {
	switch r.GetTypeId() {
	case ERROR:
		return "ERROR"
	case EMPTY:
		return "EMPTY"
	case INT:
		return "INT"
	case STRING:
		return "STRING"
	case BYTES:
		return "BYTES"
	case ARRAY:
		return "ARRAY"
	case OK:
		return "OK"
	case PONG:
		return "PONG"
	default:
		return ""
	}
}

//Set the RedisResult is Pong
func (r *RedisResult) SetPong() {
	r.resultType = PONG
	r.resultPong = true
}

//Paser if the RedisResult is PONG
func (r *RedisResult) IsPong() error {
	if r == nil {
		return ERR_REDIS_RESULT_IS_NIL
	}
	switch r.resultType {
	case PONG:
		if r.resultPong == true {
			return nil
		} else {
			return ERR_STATUS_NOT_PONG
		}
	case ERROR:
		return r.resultError
	default:
		return fmt.Errorf("The RedisResult type is not PONG, but %s!\r\n", r.GetType())
	}
}

//Set the RedisResult is OK
func (r *RedisResult) SetOk() {
	r.resultType = OK
	r.resultOk = true
}

//Paser if the RedisResult is OK
func (r *RedisResult) IsOk() error {
	if r == nil {
		return ERR_REDIS_RESULT_IS_NIL
	}
	switch r.resultType {
	case OK:
		if r.resultOk == true {
			return nil
		} else {
			return ERR_STATUS_NOT_OK
		}
	case ERROR:
		return r.resultError
	default:
		return fmt.Errorf("The RedisResult type is not OK, but %s!\r\n", r.GetType())
	}
}

//Set the RedisResult's resultError, resultType = -1
func (r *RedisResult) SetError(e error) {
	r.resultType = ERROR
	r.resultError = e
}

//Get the RedisResult's resultError
func (r *RedisResult) GetError() error {
	if r == nil {
		return ERR_REDIS_RESULT_IS_NIL
	}
	switch r.resultType {
	case ERROR:
		return r.resultError
	default:
		return fmt.Errorf("The RedisResult type is not ERROR, but %s!\r\n", r.GetType())
	}
}

//Set the RedisResult's resultInt, resultType = 1
func (r *RedisResult) SetInt(i int) {
	r.resultType = INT
	r.resultInt = i
}

//Get the RedisResult's resultInt
func (r *RedisResult) GetInt() (int, error) {
	if r == nil {
		return 0, ERR_REDIS_RESULT_IS_NIL
	}
	switch r.resultType {
	case ERROR:
		return 0, r.GetError()
	case INT:
		return r.resultInt, nil
	default:
		return 0, fmt.Errorf("The RedisResult type is not INT, but %s!\r\n", r.GetType())
	}
}

//Set the RedisResult's resultString, resultType = 2
func (r *RedisResult) SetString(str string) {
	r.resultType = STRING
	r.resultString = str
}

//Get the RedisResult's resultString
func (r *RedisResult) GetString() (string, error) {
	if r == nil {
		return "", ERR_REDIS_RESULT_IS_NIL
	}
	switch r.resultType {
	case ERROR:
		return "", r.GetError()
	case STRING:
		return r.resultString, nil
	default:
		return "", fmt.Errorf("The RedisResult type is not STRING, but %s!\r\n", r.GetType())
	}

}

//Set the RedisResult's resultBytes, resultType = 3
func (r *RedisResult) SetBytes(b []byte) {
	r.resultType = BYTES
	r.resultBytes = b
}

//Get the RedisResult's resultBytes
//if resultBytes has the content, it will be returned, and dont care the resultType
func (r *RedisResult) GetBytes() ([]byte, error) {
	if r == nil {
		return nil, ERR_REDIS_RESULT_IS_NIL
	}
	if r.resultBytes != nil {
		return r.resultBytes, nil
	}
	switch r.resultType {
	case ERROR:
		return nil, r.GetError()
	case BYTES:
		return r.resultBytes, nil
	default:
		return nil, fmt.Errorf("The RedisResult type is not BYTES, but %s!\r\n", r.GetType())
	}
}

//Set the RedisResult's resultArray, resultType = 4
func (r *RedisResult) SetArray(rr []*RedisResult) {
	r.resultType = ARRAY
	r.resultArray = rr
}

//Get the RedisResult's resultArray
//return the RedisResult array
//need more parser each ResisResults to know the concrete results
func (r *RedisResult) GetArray() ([]*RedisResult, error) {
	if r == nil {
		return nil, ERR_REDIS_RESULT_IS_NIL
	}
	switch r.resultType {
	case ERROR:
		return nil, r.GetError()
	case ARRAY:
		return r.resultArray, nil
	default:
		return nil, fmt.Errorf("The RedisResult tpye is not ARRAY, but %s!\r\n", r.GetType())
	}
}
