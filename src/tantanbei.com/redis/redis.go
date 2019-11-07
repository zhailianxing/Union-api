package redis

import (
	"errors"
	"net"
	"tantanbei.com/pool"
	"time"
)

type Redis struct {
	redisPool   *pool.Pool
	ConnTimeout time.Duration
	Password    string
}

//Initial the Redis
//ConnTimeout : set net.Dial and pool.Get timeout
//WriteTimeout and ReadTimeout : set write and read tcp timeout
func NewRedis(address string, password string, bufferCapacity int, ConnTimeout time.Duration, WriteTimeout time.Duration, ReadTimeout time.Duration) (*Redis, error) {
	r := new(Redis)
	r.ConnTimeout = ConnTimeout
	r.Password = password
	r.redisPool = pool.NewPool("redis", 10, 10,
		func(args ...interface{}) (c pool.PoolItem, e error) {
			buffer := ConnBuffer{nil, 0, 0, 0}
			con := poolConn{nil, nil, nil, &buffer, WriteTimeout, ReadTimeout, true}
			con.Init(bufferCapacity)
			con.realConn, e = net.DialTimeout("tcp", "127.0.0.1:6379", ConnTimeout)
			if password != "" {
				con.SendAuthAndReply(password)
				//config set requirepass 123
			}

			c = &con
			return
		})
	if r.redisPool.Size() != 0 {
		return nil, errors.New("should be zero size pool")
	}
	return r, nil
}

//Remove all keys from all databases
func (r *Redis) FlushAll() error {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendFlushAll()

	return pCon.ParseRedisProtocolOk()
}

//Get the value(string) of a key
func (r *Redis) GetString(key string) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGet(key)

	return pCon.ParserReturnString()
}

//Get the value(int64) of a key
func (r *Redis) GetInt64(key string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGet(key)

	return pCon.ParseRedisProtocolInt64()
}

//Get the value(int32) of a key
func (r *Redis) GetInt32(key string) (int32, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGet(key)

	return pCon.ParseRedisProtocolInt32()
}

//Get the value(int16) of a key
func (r *Redis) GetInt16(key string) (int16, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGet(key)

	return pCon.ParseRedisProtocolInt16()
}

//Get the value(int8) of a key
func (r *Redis) GetInt8(key string) (int8, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGet(key)

	return pCon.ParseRedisProtocolInt8()
}

//Get the value(bytes) of a key
func (r *Redis) GetBytes(key string) ([]byte, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return nil, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGet(key)

	return pCon.ParseRedisProtocolBytes()
}

//Set a string value of a key
//After the sec(second) has expired, the key will automatically be deleted;
//if sec <= 0, the key cannot be deleted.
func (r *Redis) Set(key string, val []byte, sec int) error {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	if sec <= 0 {

		//Send command
		pCon.SendSet(key, val)
	} else {

		//Send command
		pCon.SendSetExpire(key, val, sec)
	}

	return pCon.ParseRedisProtocolOk()
}

//Authenticate to the server
func (r *Redis) Auth(password string) error {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendAuth(password)

	return pCon.ParseRedisProtocolOk()
}

//Delete a key
func (r *Redis) Del(key ...string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendDel(key...)

	return pCon.ParseRedisProtocolInt64()
}

//Set a key's time to live in seconds
func (r *Redis) Expire(key string, sec int) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendExpire(key, sec)

	return pCon.ParserReturnInt()
}

//Get the valuse of all the given keys
func (r *Redis) MGet(key ...string) ([]string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return nil, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendMGet(key...)

	return pCon.ParserReturnArray()
}

//Increment the integer value of a key by one
func (r *Redis) Incr(key string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendIncr(key)

	return pCon.ParserReturnInt()
}

//Decrement the integer value of a key by the given number
func (r *Redis) Decr(key string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendDecr(key)

	return pCon.ParserReturnInt()
}

//Determine the type stored at key
func (r *Redis) Type(key string) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendType(key)

	return pCon.ParserReturnString()
}

//Find all keys matching the given pattern
func (r *Redis) Keys(key string) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendKeys(key)

	return pCon.ParserReturnString()
}

//Return arandom key from the keyspace
func (r *Redis) RandomKey() (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendRandomKey()

	return pCon.ParserReturnString()
}

//Rename a key
//if careNX == true, only if the new key does not exist
func (r *Redis) Rename(src string, dst string, careNX bool) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)
	if careNX != true {

		//Send command
		pCon.SendRename(src, dst)

		if pCon.ParseRedisProtocolOk() != nil {
			return 0, pCon.ParseRedisProtocolOk()
		} else {
			return 1, nil
		}
	} else {

		//Send command
		pCon.SendRenameNx(src, dst)

		return pCon.ParseRedisProtocolInt64()
	}
}

//Return the number of keys in the selected database
func (r *Redis) Dbsize() (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendDbsize()

	return pCon.ParseRedisProtocolInt64()
}

//Set the value of a key, only if the key does not exist
func (r *Redis) SetNX(key string, val []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSetNX(key, val)

	return pCon.ParserReturnInt()
}

//Get the time to live for a key
func (r *Redis) Ttl(key string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendTtl(key)

	return pCon.ParserReturnInt()
}

//Set the string value of a key and return its old value
func (r *Redis) GetSet(key string, val []byte) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGetSet(key, val)

	return pCon.ParserReturnString()
}

//Set multiple keys to multiple values
func (r *Redis) MSet(mapping map[string][]byte, careNX bool) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)
	if careNX != true {

		//Send command
		pCon.SendMSet(mapping)

		if pCon.ParserReturnIsOk() != nil {
			return 0, pCon.ParserReturnIsOk()
		} else {
			return len(mapping), pCon.ParserReturnIsOk()
		}
	} else {

		//Send command
		pCon.SendMsetnx(mapping)

		return pCon.ParserReturnInt()
	}
}

//Increment the integer value of a key by the given amount
func (r *Redis) Incrby(key string, num int) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendIncrby(key, num)

	return pCon.ParserReturnInt()
}

//Decrement the integer value of a key by the given number
func (r *Redis) Decrby(key string, num int) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendDecrby(key, num)

	return pCon.ParserReturnInt()
}

//Append a value to a key
func (r *Redis) Append(key string, val []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendAppend(key, val)

	return pCon.ParserReturnInt()
}

//Get a substring of the string strored at a key
//it is called SUBSTR in Redis versions <= 2.0.
func (r *Redis) Getrange(key string, start int, end int) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendGetrange(key, start, end)

	return pCon.ParserReturnString()
}

//Prepend(or Append) one or multiple values to a list
func (r *Redis) Push(key string, val []byte, isRight bool) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)
	if isRight == true {

		//Send command
		pCon.SendRpush(key, val)
	} else {

		//Send command
		pCon.SendLpush(key, val)
	}
	return pCon.ParserReturnInt()
}

//Get the length of a list
func (r *Redis) Llen(key string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendLlen(key)

	return pCon.ParserReturnInt()
}

//Get a range of elements from a list
func (r *Redis) Lrange(key string, start int, end int) ([]string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return nil, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendLrange(key, start, end)

	return pCon.ParserReturnArray()
}

//Trim a list to the specified range
func (r *Redis) Ltrim(key string, start int, end int) error {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendLtrim(key, start, end)

	return pCon.ParserReturnIsOk()
}

//Get an element from a list by its index
func (r *Redis) Lindex(key string, index int) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendLindex(key, index)

	return pCon.ParserReturnString()
}

//Set the value of an element in a list by its index
func (r *Redis) LSet(key string, index int, value []byte) error {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendLSet(key, index, value)

	return pCon.ParserReturnIsOk()
}

//Remove elements from a list
func (r *Redis) Lrem(key string, count int, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendLrem(key, count, value)

	return pCon.ParserReturnInt()
}

//Remove and get the last(or first) element in a list
func (r *Redis) Pop(key string, isRight bool) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)
	if isRight == true {

		//Send command
		pCon.SendRpop(key)
	} else {

		//Send command
		pCon.SendLpop(key)
	}
	return pCon.ParserReturnString()
}

//Add one or more members to a set
func (r *Redis) Sadd(key string, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSadd(key, value)

	return pCon.ParserReturnInt()
}

//Remove one or more members from a set
func (r *Redis) Srem(key string, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSrem(key, value)

	return pCon.ParserReturnInt()
}

//Remove and return one or multiple random members from a set
func (r *Redis) Spop(key string) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSpop(key)

	return pCon.ParserReturnString()
}

//Move a member from one set to another
func (r *Redis) Smove(src string, dst string, val []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)

	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSmove(src, dst, val)

	return pCon.ParserReturnInt()
}

//Get the number of members in a set
func (r *Redis) Scard(key string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendScard(key)

	return pCon.ParserReturnInt()
}

//Determine if a given value is a member of a set
func (r *Redis) Sismember(key string, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSismember(key, value)

	return pCon.ParserReturnInt()
}

//Intersect multiple sets
func (r *Redis) Sinter(keys ...string) (string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return "", err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSinter(keys...)

	return pCon.ParserReturnString()
}

//Add multiple sets
func (r *Redis) Sunion(keys ...string) ([]string, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return nil, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendSunion(keys...)

	return pCon.ParserReturnArray()
}

//Remove one or more members from a sorted set
func (r *Redis) Zrem(key string, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendZrem(key, value)

	return pCon.ParserReturnInt()
}

//Determine the index of a member in a sorted set
func (r *Redis) Zrank(key string, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendZrank(key, value)

	return pCon.ParserReturnInt()
}

//Determine the index of a member in a sorted set, with scores ordered from high to low
func (r *Redis) Zrevrank(key string, value []byte) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendZrevrank(key, value)

	return pCon.ParserReturnInt()
}

func (r *Redis) Zadd(key string, score string, member string) (int, error) {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return 0, err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	pCon.SendZadd(key, score, member)

	return pCon.ParserReturnInt()
}

//Ping the server
func (r *Redis) Ping() error {
	c, err := r.redisPool.Get(r.ConnTimeout)
	if err != nil {
		return err
	}
	defer r.redisPool.Put(c)

	pCon := c.(*poolConn)

	//Send command
	pCon.SendPing()

	return pCon.ParseRedisProtocolPong()
}
