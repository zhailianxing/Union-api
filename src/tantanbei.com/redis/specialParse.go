package redis

import (
	"errors"
	"time"
)

//This function parses return int specially
//-2^63 ( -9,223,372,036,854,775,808) ~ 2^63-1 ( +9,223,372,036,854,775,807 )
func (c *poolConn) ParseRedisProtocolInt64() (int, error) {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return 0, err
		}
	}

	fristChar, err := c.ReadOneBuffer()
	if err != nil {
		return 0, err
	}

	switch fristChar {
	case '$':
		//$ is bulk reply
		//A "$" byte followed by the number of bytes composing the string (a prefixed length),
		// terminated by CRLF ("\r\n").
		//The actual string data.
		//A final CRLF.

		num, err := c.ReadInt()
		if err != nil {
			return 0, err
		}
		if num == -1 {
			//Null Bulk String
			//not return an empty string, but a nil object.

			return 0, ERR_REPLY_IS_NIL
		}
		result, err := c.ReadBuffer(num)
		if err != nil {
			return 0, err
		}
		if len(result) > 20 {
			return 0, ERR_OVERFLOWS
		} else {
			//but the int is returned by this bulk reply
			return MyBstoI64(result)
		}

	case ':':
		//: is integer reply
		//terminated by CRLF ("\r\n").

		num, err := c.ReadInt()
		if err != nil {
			return 0, err
		}
		return num, nil
	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return 0, err
		} else {
			return 0, errors.New(string(result))
		}

	default:
		return 0, ERR_NOT_INT
	}
}

// -2,147,483,648 ~ +2,147,483,647
func (c *poolConn) ParseRedisProtocolInt32() (int32, error) {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return 0, err
		}
	}

	fristChar, err := c.ReadOneBuffer()
	if err != nil {
		return 0, err
	}

	switch fristChar {
	case '$':
		//$ is bulk reply
		//A "$" byte followed by the number of bytes composing the string (a prefixed length),
		// terminated by CRLF ("\r\n").
		//The actual string data.
		//A final CRLF.

		num, err := c.ReadInt()
		if err != nil {
			return 0, err
		}
		if num == -1 {
			//Null Bulk String
			//not return an empty string, but a nil object.

			return 0, ERR_REPLY_IS_NIL
		}
		result, err := c.ReadBuffer(num)
		if err != nil {
			return 0, err
		}
		if len(result) > 11 {
			return 0, ERR_OVERFLOWS
		} else {
			//but the int is returned by this bulk reply
			return MyBstoI32(result)
		}

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return 0, err
		} else {
			return 0, errors.New(string(result))
		}

	default:
		return 0, ERR_NOT_INT
	}
}

//  -32768 ~ +32767
func (c *poolConn) ParseRedisProtocolInt16() (int16, error) {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return 0, err
		}
	}

	fristChar, err := c.ReadOneBuffer()
	if err != nil {
		return 0, err
	}

	switch fristChar {
	case '$':
		//$ is bulk reply
		//A "$" byte followed by the number of bytes composing the string (a prefixed length),
		// terminated by CRLF ("\r\n").
		//The actual string data.
		//A final CRLF.

		num, err := c.ReadInt()
		if err != nil {
			return 0, err
		}
		if num == -1 {
			//Null Bulk String
			//not return an empty string, but a nil object.

			return 0, ERR_REPLY_IS_NIL
		}
		result, err := c.ReadBuffer(num)
		if err != nil {
			return 0, err
		}
		if len(result) > 6 {
			return 0, ERR_OVERFLOWS
		} else {
			//but the int is returned by this bulk reply
			return MyBstoI16(result)
		}

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return 0, err
		} else {
			return 0, errors.New(string(result))
		}

	default:
		return 0, ERR_NOT_INT
	}
}

//  -128 ~ 127
func (c *poolConn) ParseRedisProtocolInt8() (int8, error) {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return 0, err
		}
	}

	fristChar, err := c.ReadOneBuffer()
	if err != nil {
		return 0, err
	}

	switch fristChar {
	case '$':
		//$ is bulk reply
		//A "$" byte followed by the number of bytes composing the string (a prefixed length),
		// terminated by CRLF ("\r\n").
		//The actual string data.
		//A final CRLF.

		num, err := c.ReadInt()
		if err != nil {
			return 0, err
		}
		if num == -1 {
			//Null Bulk String
			//not return an empty string, but a nil object.

			return 0, ERR_REPLY_IS_NIL
		}
		result, err := c.ReadBuffer(num)
		if err != nil {
			return 0, err
		}
		if len(result) > 4 {
			return 0, ERR_OVERFLOWS
		} else {
			//but the int is returned by this bulk reply
			return MyBstoI8(result)
		}

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return 0, err
		} else {
			return 0, errors.New(string(result))
		}

	default:
		return 0, ERR_NOT_INT
	}
}

//This function is parses return bytes specially
func (c *poolConn) ParseRedisProtocolBytes() ([]byte, error) {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return nil, err
		}
	}

	fristChar, err := c.ReadOneBuffer()

	if err != nil {
		return nil, err
	}
	switch fristChar {
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
		//fmt.Println("needSize:", num)
		if num == -1 {
			//Null Bulk String
			//not return an empty string, but a nil object.

			return nil, ERR_REPLY_IS_NIL
		}
		return c.ReadBuffer(num)

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New(string(result))
		}

	default:
		return nil, ERR_NOT_BYTES
	}
}

//This function is parses return OK specially
func (c *poolConn) ParseRedisProtocolOk() error {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return err
		}
	}

	fristChar, err := c.ReadOneBuffer()
	if err != nil {
		return err
	}
	switch fristChar {
	case '+':
		//+ is status reply
		//terminated by CRLF ("\r\n").

		ok, err := c.ReadOK()
		if err != nil {
			return err
		}
		if ok == true {
			return nil
		} else {
			return ERR_STATUS_NOT_OK
		}

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return err
		} else {
			return errors.New(string(result))
		}

	default:
		return ERR_NOT_STATUS
	}
}

//This function is parses return PONG specially
func (c *poolConn) ParseRedisProtocolPong() error {
	if c.readTimeout != 0 {
		if err := c.realConn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			c.err = err
			return err
		}
	}

	fristChar, err := c.ReadOneBuffer()
	if err != nil {
		return err
	}
	switch fristChar {
	case '+':
		//+ is status reply
		//terminated by CRLF ("\r\n").

		ok, err := c.ReadPong()
		if err != nil {
			return err
		}
		if ok == true {
			return nil
		} else {
			return ERR_STATUS_NOT_PONG
		}

	case '-':
		//- is error reply
		//terminated by CRLF ("\r\n").

		result, err := c.ReadPartSafe()
		if err != nil {
			return err
		} else {
			return errors.New(string(result))
		}

	default:
		return ERR_NOT_STATUS
	}
}
