package redis

import (
	"errors"
)

const (
	//Defind the RedisResult's type
	ERROR  = -1
	EMPTY  = 0
	INT    = 1
	STRING = 2
	BYTES  = 3
	ARRAY  = 4
	OK     = 5
	PONG   = 6

	//Command name
	FLUSHALL  = "FLUSHALL"
	GET       = "GET"
	SET       = "SET"
	DEL       = "DEL"
	EXPIRE    = "EXPIRE"
	SETEX     = "SETEX"
	AUTH      = "AUTH"
	TYPE      = "TYPE"
	KEYS      = "KEYS"
	RANDOMKEY = "RANDOMKEY"
	RENAME    = "RENAME"
	RENAMENX  = "RENAMENX"
	DBSIZE    = "DBSIZE"
	TTL       = "TTL"
	GETSET    = "GETSET"
	MGET      = "MGET"
	SETNX     = "SETNX"
	MSET      = "MSET"
	MSETNX    = "MSETNX"
	INCR      = "INCR"
	INCRBY    = "INCRBY"
	DECR      = "DECR"
	DECRBY    = "DECRBY"
	APPEND    = "APPEND"
	GETRANGE  = "GETRANGE"
	RPUSH     = "RPUSH"
	LPUSH     = "LPUSH"
	LLEN      = "LLEN"
	LRANGE    = "LRANGE"
	LTRIM     = "LTRIM"
	LINDEX    = "LINDEX"
	LSET      = "LSET"
	LREM      = "LREM"
	LPOP      = "LPOP"
	RPOP      = "RPOP"
	SADD      = "SADD"
	SREM      = "SREM"
	SPOP      = "SPOP"
	SMOVE     = "SMOVE"
	SCARD     = "SCARD"
	SISMEMBER = "SISMEMBER"
	SINTER    = "SINTER"
	SUNION    = "SUNION"
	ZREM      = "ZREM"
	ZRANK     = "ZRANK"
	ZREVRANK  = "ZREVRANK"
	ZADD      = "ZADD"
	PING      = "PING"
)

var ERR_NOT_INT = errors.New("It is not INT!")
var ERR_NOT_BYTES = errors.New("It is not BYTES!")
var ERR_REPLY_IS_NIL = errors.New("It is nil object!")
var ERR_STATUS_NOT_OK = errors.New("Status is not OK!")
var ERR_STATUS_NOT_PONG = errors.New("Status is not PONG!")
var ERR_NOT_STATUS = errors.New("It is not status reply!")
var ERR_REDIS_RESULT_IS_NIL = errors.New("The RedisResult is nil!")
var ERR_TIME_OUT = errors.New("Time out!")
var ERR_OVERFLOWS = errors.New("Constant overflows!")
