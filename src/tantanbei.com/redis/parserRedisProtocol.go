package redis

import (
	"errors"
	"time"
)

//Parser the Redis's protocol
//return RedisResult struct
func (c *poolConn) ParserRedisProtocol() (*RedisResult, error) {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return nil, err
		}
	}

	RedisResult := &RedisResult{}

	fristChar, err := c.ReadOneBuffer()

	if err != nil {
		return nil, err
	}
	switch fristChar {
	case '+':
		//+ is status reply
		//terminated by CRLF ("\r\n").

		ok, err := c.ReadOK()
		if err != nil {
			return nil, err
		}
		if ok == true {
			RedisResult.SetOk()
		} else {
			pong, err := c.ReadPong()
			if err != nil {
				return nil, err
			}
			if pong == true {
				RedisResult.SetPong()
			} else {
				result, err := c.ReadPartSafe()
				if err != nil {
					return nil, err
				}
				RedisResult.SetString(string(result))
			}
		}

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return nil, err
		}
		RedisResult.SetError(errors.New(string(result)))

	case ':':
		//: is integer reply
		//terminated by CRLF ("\r\n").

		num, err := c.ReadInt()
		if err != nil {
			return nil, err
		}
		RedisResult.SetInt(num)

	case '$':
		//$ is bulk reply
		//A "$" byte followed by the number of bytes composing the string (a prefixed length),
		// terminated by CRLF ("\r\n").
		//The actual string data.
		//A final CRLF.

		//this num means sizes of the bulk replay.
		num, err := c.ReadInt()
		if err != nil {
			return nil, err
		}
		if num == -1 {
			//Null Bulk String
			//return a null object and not an empty Array

			return nil, nil
		}
		result, err := c.ReadBuffer(num)
		if err != nil {
			return nil, err
		}

		//the bulk reply is a string data
		//but sometime, the bytes is also be used
		//the return has bytes and string, the type is string
		RedisResult.SetBytes(result)

		//but the int is returned by this bulk reply
		//number, err := strconv.Atoi(string(result))
		number, err := MyBstoI64(result)
		if err == nil {
			RedisResult.SetInt(number)
		} else {
			RedisResult.SetString(string(result))
		}

	case '*':
		//* is multi bulk reply
		//A * character as the first byte,
		// followed by the number of elements in the array as a decimal number,
		// followed by CRLF.
		//An additional RESP(Redis Serialization Protocol) type for every element of the Array.

		//this num means how many bulks in this array reply
		num, err := c.ReadInt()
		if err != nil {
			return nil, err
		}
		if num == -1 {
			//Null Array
			//return a null object and not an empty Array

			return nil, nil
		} else if num == 0 {
			// an empty Array

			return RedisResult, nil
		}
		result, err := c.ParserRedisProtocolStar(num)
		if err != nil {
			return nil, err
		}
		RedisResult.SetArray(result)
	}

	return RedisResult, nil
}

//Parser the array reply
//Arrays the first byte of the reply is "*"
func (c *poolConn) ParserRedisProtocolStar(num int) ([]*RedisResult, error) {
	redisResult := &RedisResult{}
	var redisResults []*RedisResult
	for i := 0; i < num; i++ {
		fristChar, err := c.ReadOneBuffer()
		if err != nil {
			return nil, err
		}
		switch fristChar {
		case '+':
			//+ is status reply
			//terminated by CRLF ("\r\n").

			ok, err := c.ReadOK()
			if err != nil {
				return nil, err
			}
			if ok == true {
				redisResult.SetOk()
				redisResults = append(redisResults, redisResult)
			} else {
				pong, err := c.ReadPong()
				if err != nil {
					return nil, err
				}
				if pong == true {
					redisResult.SetPong()
					redisResults = append(redisResults, redisResult)
				} else {
					result, err := c.ReadPartSafe()
					if err != nil {
						return nil, err
					}
					redisResult.SetString(string(result))
					redisResults = append(redisResults, redisResult)
				}
			}

		case '-':
			//- is error reply
			//terminated by CRLF ("\r\n").

			result, err := c.ReadPartSafe()
			if err != nil {
				return nil, err
			}
			redisResult.SetError(errors.New(string(result)))
			redisResults = append(redisResults, redisResult)

		case ':':
			//: is integer reply
			//terminated by CRLF ("\r\n").

			num, err := c.ReadInt()
			if err != nil {
				return nil, err
			}
			redisResult.SetInt(num)
			redisResults = append(redisResults, redisResult)

		case '$':
			//$ is bulk reply
			//A "$" byte followed by the number of bytes composing the string (a prefixed length),
			// terminated by CRLF ("\r\n").
			//The actual string data.
			//A final CRLF.

			num, err := c.ReadInt()
			if err != nil {
				return nil, err
			}
			if num == -1 {
				//Null Bulk String
				//not return an empty string, but a nil object.

				redisResults = append(redisResults, nil)
				continue
			}
			result, err := c.ReadBuffer(num)
			if err != nil {
				return nil, err
			}

			//the bulk reply is a string data
			//but sometime, the bytes is also be used
			//the return has bytes and string, the type is string
			redisResult.SetBytes(result)

			//but the int is returned by this bulk reply
			//number, err := strconv.Atoi(string(result))
			number, err := MyBstoI64(result)
			if err == nil {
				redisResult.SetInt(number)
				redisResults = append(redisResults, redisResult)
			} else {
				redisResult.SetString(string(result))
				redisResults = append(redisResults, redisResult)
			}
		}
	}

	return redisResults, nil
}

//Frist: Parser the redis protocol and return RedisResult struct
//Second: Parser if this RedisResult struct is OK type
//return: error!=nil means is OK.
func (c *poolConn) ParserReturnIsOk() error {
	result, err := c.ParserRedisProtocol()
	if err != nil {
		return err
	}
	return result.IsOk()
}

//Frist: Parser the redis protocol and return RedisResult struct
//Second: Parser if this RedisResult struct is PONG type
//return: error!=nil means is PONG.
func (c *poolConn) ParserReturnIsPong() error {
	result, err := c.ParserRedisProtocol()
	if err != nil {
		return err
	}
	return result.IsPong()
}

//Frist: Parser the redis protocol and return RedisResult struct
//Second: Parser if this RedisResult struct is STRING type
//return: string type of result.
func (c *poolConn) ParserReturnString() (string, error) {
	result, err := c.ParserRedisProtocol()
	if err != nil {
		return "", err
	}
	return result.GetString()
}

//Frist: Parser the redis protocol and return RedisResult struct
//Second: Parser if this RedisResult struct is INT type
//return: int type of result
func (c *poolConn) ParserReturnInt() (int, error) {
	result, err := c.ParserRedisProtocol()
	if err != nil {
		return 0, err
	}
	return result.GetInt()
}

//Frist: Parser the redis protocol and return RedisResult struct
//Second: Parser if this RedisResult struct is BYTES type
//return: []byte type of result
func (c *poolConn) ParserReturnBytes() ([]byte, error) {
	result, err := c.ParserRedisProtocol()
	if err != nil {
		return nil, err
	}
	return result.GetBytes()
}

//Frist: Parser the redis protocol and return RedisResult struct
//Second: Parser if this RedisResult struct is ARRAY type
//return: string type of each result
func (c *poolConn) ParserReturnArray() ([]string, error) {
	result, err := c.ParserRedisProtocol()
	if err != nil {
		return nil, err
	}

	as := make([]string, 0)
	results, err := result.GetArray()
	if err != nil {
		return nil, err
	}
	for _, a := range results {
		aa, err := a.GetString()
		if err != nil {
			as = append(as, err.Error())
		} else {
			as = append(as, aa)
		}
	}
	return as, nil
}
