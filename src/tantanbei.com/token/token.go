package token

/*IMPORTANT..DO NOT TOUCH THIS FILE UNLESS YOU ARE 100% SURE YOU KNOW WHAT YOU ARE DOING !!!!
.
.
.
.
.
.
.
.
.
.
.
DON'T MODIY!!!!!!
.
.
.
.
.
.
.
.
.
.
.
*/

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"sync"
)

const (
	//DON"T CHANGE!!!!!!
	TOKEN_SIZE    = 64
	TOKEN_ID_SIZE = 32
)

//IMPORTANT..DO NOT CHANGE THIS SECRET!!!!!!
var TOKEN_HASH_SECRET = []byte("j4083ut3&*u430JFOWEN")
var TOKEN_HASH_GEN_SIZE = len(TOKEN_HASH_SECRET) + 32

var TOKEN_HASH_BYTES_POOL = sync.Pool{
	New: func() interface{} {
		return make([]byte, TOKEN_HASH_GEN_SIZE)
	},
}

//generate token with userid as part of forgery protection
func Generate() []byte {
	t := make([]byte, TOKEN_SIZE)

	//32-64 is pure random...
	rand.Read(t[32:])

	hash := createHash(t[32:])

	//	logx.D("hash:", base64.URLEncoding.EncodeToString(hash))

	//copy hash protection to first 32bytes
	copy(t[0:32], hash)

	//	logx.D("token:", base64.URLEncoding.EncodeToString(t))

	return t
}

//only generate the 32 id part...
func GenerateId() []byte {
	t := make([]byte, TOKEN_ID_SIZE)

	//32-64 is pure random...
	rand.Read(t)

	return t
}

//validate the format of the token including forgery protection
func Validate(token []byte) (isValid bool) {
	//fast path
	if len(token) != TOKEN_SIZE {
		return
	}

	hash := createHash(token[32:])

	//now compare our hash with input
	if bytes.Compare(token[0:32], hash) == 0 {
		isValid = true
	}

	return
}

//hash fuction generates an anti-forgery key that is 32 bytes..
func createHash(tokenId []byte) []byte {
	u := TOKEN_HASH_BYTES_POOL.Get().([]byte)
	defer TOKEN_HASH_BYTES_POOL.Put(u)

	//let's generate real hash...
	copy(u, TOKEN_HASH_SECRET)
	copy(u[len(TOKEN_HASH_SECRET):], tokenId)

	//	logx.D("u:", base64.URLEncoding.EncodeToString(u), ", len(U):", len(u))

	hash := sha256.Sum256(u)

	return hash[0:32]
}

//pass in a token already with [32:64] with real id, and let's regenerate the hash part..
func RecreateHash(token []byte) {
	copy(token, createHash(token[32:64]))
}
