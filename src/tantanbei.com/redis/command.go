package redis

import (
	"strconv"
	"tantanbei.com/strconvx"
	"time"
	"bufio"
	"errors"
)

//Send command
//frist : the name of cmd
//then : import the parameters. The order accords with the Redis standard
func (c *poolConn) SendCommand(cmd string, args ...string) (err error) {
	if cmd == "" {
		return errors.New("CMD can not a nil!")
	}

	defer func() {
		c.err = err
	}()

	if c.writeTimeout != 0 {
		if err = c.realConn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return
		}
	}

	if err = c.CommandBytes(cmd, args...); err != nil {
		return
	}

	return c.cmdbuf.Flush()
}

//Send command , if the size of arguments is large
func (c *poolConn) SendCommandLargePayload(cmd string, args ...[]byte) (err error) {
	if cmd == "" {
		return errors.New("CMD can not a nil!")
	}

	defer func() {
		c.err = err
	}()

	if c.writeTimeout != 0 {
		if err = c.realConn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return err
		}
	}

	//initial the c.cmdbuf
	if c.cmdbuf == nil {
		c.cmdbuf = bufio.NewWriter(c.realConn)
	} else {
		c.cmdbuf.Reset(c.realConn)
	}

	//Frist: Link the cmd name
	c.cmdbuf.WriteString("*")
	c.cmdbuf.WriteString(strconv.Itoa(len(args) + 1))
	c.cmdbuf.WriteString("\r\n")
	c.cmdbuf.WriteString("$")
	c.cmdbuf.WriteString(strconv.Itoa(len(cmd)))
	c.cmdbuf.WriteString("\r\n")
	c.cmdbuf.WriteString(cmd)
	c.cmdbuf.WriteString("\r\n")

	err = c.cmdbuf.Flush()
	if err != nil {
		return err
	}

	//Second: Link the arguments after the cmd
	for _, s := range args {
		if s == nil {
			return errors.New("Command parameters can not a nil!")
		}
		if len(s) >= 100 {
			c.cmdbuf.Reset(c.realConn)
			c.cmdbuf.WriteString("$")
			c.cmdbuf.WriteString(strconv.Itoa(len(s)))
			c.cmdbuf.WriteString("\r\n")

			err = c.cmdbuf.Flush()
			if err != nil {
				return err
			}

			_, err = c.realConn.Write(s)
			if err != nil {
				return err
			}

			c.cmdbuf.Reset(c.realConn)
			c.cmdbuf.WriteString("\r\n")

			err = c.cmdbuf.Flush()
			if err != nil {
				return err
			}
		} else {
			c.cmdbuf.Reset(c.realConn)
			c.cmdbuf.WriteString("$")
			c.cmdbuf.WriteString(strconv.Itoa(len(s)))
			c.cmdbuf.WriteString("\r\n")
			c.cmdbuf.Write(s)
			c.cmdbuf.WriteString("\r\n")

			err = c.cmdbuf.Flush()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//Parser the command parameters
func (c *poolConn) CommandBytes(cmd string, args ...string) error {
	//initial the c.cmdbuf
	if c.cmdbuf == nil {
		c.cmdbuf = bufio.NewWriter(c.realConn)
	} else {
		c.cmdbuf.Reset(c.realConn)
	}

	//Frist: Link the cmd name
	c.cmdbuf.WriteString("*")
	c.cmdbuf.WriteString(strconv.Itoa(len(args) + 1))
	c.cmdbuf.WriteString("\r\n")
	c.cmdbuf.WriteString("$")
	c.cmdbuf.WriteString(strconv.Itoa(len(cmd)))
	c.cmdbuf.WriteString("\r\n")
	c.cmdbuf.WriteString(cmd)
	c.cmdbuf.WriteString("\r\n")

	//Second: Link the arguments after the cmd
	for _, s := range args {
		if s == "" {
			return errors.New("Command parameters can not a nil!")
		}
		c.cmdbuf.WriteString("$")
		c.cmdbuf.WriteString(strconv.Itoa(len(s)))
		c.cmdbuf.WriteString("\r\n")
		c.cmdbuf.WriteString(s)
		c.cmdbuf.WriteString("\r\n")
	}
	return nil
}

func (c *poolConn) SendFlushAll() {
	c.SendCommand(FLUSHALL)
}

func (c *poolConn) SendGet(key string) {
	c.SendCommand(GET, key)
}

func (c *poolConn) SendSet(key string, val []byte) {
	if len(val) < 100 {
		c.SendCommand(SET, key, string(val))
	} else {
		c.SendCommandLargePayload(SET, []byte(key), val)
	}
}

func (c *poolConn) SendDel(key ...string) {
	c.SendCommand(DEL, key...)
}

func (c *poolConn) SendExpire(key string, sec int) {
	c.SendCommand(EXPIRE, key, strconv.Itoa((sec)))
}

func (c *poolConn) SendSetExpire(key string, val []byte, sec int) {
	if len(val) < 100 {
		c.SendCommand(SETEX, key, strconv.Itoa(sec), string(val))
	} else {
		c.SendCommandLargePayload(SETEX, []byte(key), strconvx.Itobs(sec), val)
	}
}

func (c *poolConn) SendAuth(password string) {
	c.SendCommand(AUTH, password)
}

func (c *poolConn) SendAuthAndReply(password string) {
	c.SendCommand(AUTH, password)
	err := c.ParseRedisProtocolOk()
	if err != nil {
		panic(err)
	}
}

func (c *poolConn) SendMGet(key ...string) {
	c.SendCommand(MGET, key...)
}

func (c *poolConn) SendIncr(key string) {
	c.SendCommand(INCR, key)
}

func (c *poolConn) SendDecr(key string) {
	c.SendCommand(DECR, key)
}

func (c *poolConn) SendType(key string) {
	c.SendCommand(TYPE, key)
}

func (c *poolConn) SendKeys(key string) {
	c.SendCommand(KEYS, key)
}

func (c *poolConn) SendRandomKey() {
	c.SendCommand(RANDOMKEY)
}

func (c *poolConn) SendRename(src string, dst string) {
	c.SendCommand(RENAME, src, dst)
}

func (c *poolConn) SendRenameNx(src string, dst string) {
	c.SendCommand(RENAMENX, src, dst)
}

func (c *poolConn) SendDbsize() {
	c.SendCommand(DBSIZE)
}

func (c *poolConn) SendTtl(key string) {
	c.SendCommand(TTL, key)
}

func (c *poolConn) SendGetSet(key string, val []byte) {
	c.SendCommand(GETSET, key, string(val))
}

func (c *poolConn) SendSetNX(key string, val []byte) {
	if len(val) < 100 {
		c.SendCommand(SETNX, key, string(val))
	} else {
		c.SendCommandLargePayload(SETNX, []byte(key), val)
	}

}

func (c *poolConn) SendMSet(mapping map[string][]byte) {
	args := make([]string, len(mapping)*2)
	i := 0
	for k, v := range mapping {
		args[i] = k
		args[i+1] = string(v)
		i += 2
	}
	c.SendCommand(MSET, args...)
}

func (c *poolConn) SendMsetnx(mapping map[string][]byte) {
	args := make([]string, len(mapping)*2)
	i := 0
	for k, v := range mapping {
		args[i] = k
		args[i+1] = string(v)
		i += 2
	}
	c.SendCommand(MSET, args...)
}

func (c *poolConn) SendIncrby(key string, num int) {
	c.SendCommand(INCRBY, key, strconv.Itoa(num))
}

func (c *poolConn) SendDecrby(key string, num int) {
	c.SendCommand(DECRBY, key, strconv.Itoa(num))
}

func (c *poolConn) SendAppend(key string, val []byte) {
	c.SendCommand(APPEND, key, string(val))
}

func (c *poolConn) SendGetrange(key string, start int, end int) {
	c.SendCommand(GETRANGE, key, strconv.Itoa(start), strconv.Itoa(end))
}

func (c *poolConn) SendRpush(key string, val []byte) {
	c.SendCommand(RPUSH, key, string(val))
}

func (c *poolConn) SendLpush(key string, val []byte) {
	c.SendCommand(LPUSH, key, string(val))
}

func (c *poolConn) SendLlen(key string) {
	c.SendCommand(LLEN, key)
}

func (c *poolConn) SendLrange(key string, start int, end int) {
	c.SendCommand(LRANGE, key, strconv.Itoa(start), strconv.Itoa(end))
}

func (c *poolConn) SendLtrim(key string, start int, end int) {
	c.SendCommand(LTRIM, key, strconv.Itoa(start), strconv.Itoa(end))
}

func (c *poolConn) SendLindex(key string, index int) {
	c.SendCommand(LINDEX, key, strconv.Itoa(index))
}

func (c *poolConn) SendLSet(key string, index int, value []byte) {
	c.SendCommand(LSET, key, strconv.Itoa(index), string(value))
}

func (c *poolConn) SendLrem(key string, count int, value []byte) {
	c.SendCommand(LREM, key, strconv.Itoa(count), string(value))
}

func (c *poolConn) SendLpop(key string) {
	c.SendCommand(LPOP, key)
}

func (c *poolConn) SendRpop(key string) {
	c.SendCommand(RPOP, key)
}

func (c *poolConn) SendSadd(key string, value []byte) {
	c.SendCommand(SADD, key, string(value))
}

func (c *poolConn) SendSrem(key string, value []byte) {
	c.SendCommand(SREM, key, string(value))
}

func (c *poolConn) SendSpop(key string) {
	c.SendCommand(SPOP, key)
}

func (c *poolConn) SendSmove(src string, dst string, val []byte) {
	c.SendCommand(SMOVE, src, dst, string(val))
}

func (c *poolConn) SendScard(key string) {
	c.SendCommand(SCARD, key)
}

func (c *poolConn) SendSismember(key string, value []byte) {
	c.SendCommand(SISMEMBER, key, string(value))
}

func (c *poolConn) SendSinter(keys ...string) {
	c.SendCommand(SINTER, keys...)
}

func (c *poolConn) SendSunion(keys ...string) {
	c.SendCommand(SUNION, keys...)
}

func (c *poolConn) SendZrem(key string, value []byte) {
	c.SendCommand(ZREM, key, string(value))
}

func (c *poolConn) SendZrank(key string, value []byte) {
	c.SendCommand(ZRANK, key, string(value))
}

func (c *poolConn) SendZrevrank(key string, value []byte) {
	c.SendCommand(ZREVRANK, key, string(value))
}

func (c *poolConn) SendZadd(key string, score string, member string) {
	c.SendCommand(ZADD, key, score, member)
}

func (c *poolConn) SendPing() {
	c.SendCommand(PING)
}
