package redis

import (
	"errors"
)

//if c.readTimeout <= 0, just wait until the result be returned
func (c *poolConn) ReadTcpBlock() (err error) {
	c.buffer.size, err = c.realConn.Read(c.buffer.realBuffer)
	if err != nil {
		c.err = err
		return err
	}
	return nil
}

//if c.readTimeout <= 0, just wait until the result be returned
func (c *poolConn) ReadTcpBigBlockLink(bigBuffer []byte) ([]byte, error) {
	readSize := c.buffer.size - c.buffer.index

	//Make sure that the bigBuffer is completely filled
	for len(bigBuffer) != readSize {
		subReadSize, err := c.realConn.Read(bigBuffer[readSize:])
		readSize = readSize + subReadSize
		if err != nil {
			c.err = err
			return nil, err
		}
	}
	return bigBuffer, nil
}

//use in ReadUnsafeBuffer() link the new buffer
//if c.readTimeout <= 0, just wait until the result be returned
func (c *poolConn) ReadTcpBlockLink(sizeRemain int) error {
	readSize, err := c.realConn.Read(c.buffer.realBuffer[sizeRemain:])
	if err != nil {
		c.err = err
		return err
	} else {
		c.buffer.size = readSize + sizeRemain
	}
	return nil
}

//frist the buffer.index point the frist result position
//end the buffer.index point the frist of next reply
//return the []byte by size, you can convert to string or others
func (c *poolConn) ReadBuffer(size int) ([]byte, error) {
	if c.mustRead == true {
		err := c.ReadTcpBlock()
		if err != nil {
			c.err = err
			return nil, err
		}
		c.buffer.index = 0
		c.mustRead = false
	}

	//if size < c.buffer.size-c.buffer.index, normal stitching
	//if c.buffer.size-c.buffer.index < size < c.buffer.capacity-c.buffer.size+c.buffer.index, move usable data in buffer to front
	//if size > c.buffer.capacity, directly read the specified size
	if size+2 <= c.buffer.size-c.buffer.index {

		if c.buffer.realBuffer[c.buffer.index+size] == '\r' && c.buffer.realBuffer[c.buffer.index+size+1] == '\n' {
			cpy_index := c.buffer.index
			c.buffer.index = c.buffer.index + size + 2
			if c.buffer.index >= c.buffer.size {
				c.mustRead = true
			}
			return c.buffer.realBuffer[cpy_index: cpy_index+size], nil
		} else {
			return nil, errors.New("ReadBuffer is read wrong!")
		}
	} else if size+2 <= c.buffer.capacity-c.buffer.size+c.buffer.index {
		c.ReadUnsafeBuffer()
		if c.buffer.realBuffer[c.buffer.index+size] == '\r' && c.buffer.realBuffer[c.buffer.index+size+1] == '\n' {
			c.buffer.index = c.buffer.index + size + 2
			if c.buffer.index >= c.buffer.size {
				c.mustRead = true
			}
			return c.buffer.realBuffer[0:size], nil
		} else {
			return nil, errors.New("ReadBuffer is read wrong!")
		}

	} else {
		var err error
		bigBuffer := make([]byte, size+2)
		copy(bigBuffer, c.buffer.realBuffer[c.buffer.index:])

		//Make the results right , when the BigSize < buffer.capacity
		if len(bigBuffer) > c.buffer.size-c.buffer.index {
			bigBuffer, err = c.ReadTcpBigBlockLink(bigBuffer)
			if err != nil {
				return nil, err
			}
		}

		//judge weather the bigBuffer is right
		if bigBuffer[size] == '\r' && bigBuffer[size+1] == '\n' {
			c.buffer.index = c.buffer.index + size + 2
			if c.buffer.index >= c.buffer.size {
				c.mustRead = true
			}
			return bigBuffer[:size], nil
		} else {
			return nil, errors.New("bigBuffer is read wrong!")
		}
	}
}

//just read one byte from buffer.realBuffer
//and the buffer.index += 1
func (c *poolConn) ReadOneBuffer() (byte, error) {
	if c.mustRead == true {
		err := c.ReadTcpBlock()
		if err != nil {
			return 0, err
		}
		c.buffer.index = 0
		c.mustRead = false
	}
	c.buffer.index++
	return c.buffer.realBuffer[c.buffer.index-1], nil
}

//in it the buffer.index point after the '-' ,until cut by "\r\n"
//after , the buffer.index point the frist of next reply
func (c *poolConn) ReadPartSafe() ([]byte, error) {
	i := 0
	sign := 0
	result := make([]byte, 0)
	for sign != 2 {
		b, err := c.ReadOneBuffer()
		if err != nil {
			return nil, err
		}

		if c.buffer.index >= c.buffer.size {
			c.mustRead = true
		}

		//judge the end is "\r\n"
		if sign == 0 {
			if b == '\r' {
				sign++
			}
		} else if sign == 1 {
			if b == '\n' {
				sign++
			} else {
				sign = 0
			}
		}
		result = append(result, b)
		i++
	}
	return result[0: len(result)-2], nil
}

//do not copy, just do by c.buffer.realBuffer
func (c *poolConn) ReadUnsafeBuffer() error {

	//judge whether the buffer can be moved
	if c.buffer.index < c.buffer.size/2 {
		return errors.New("The bytes need to move is too long!")
	}

	j := 0
	for i := c.buffer.index; i < c.buffer.size; i, j = i+1, j+1 {
		c.buffer.realBuffer[j] = c.buffer.realBuffer[i]
	}
	c.buffer.index = 0
	return c.ReadTcpBlockLink(j)
}

//in it the buffer.index point the ':' '*' just by realBuffer
//after this the buffer.index point the frist of next reply
func (c *poolConn) ReadPart() ([]byte, error) {
	for {
		for i := c.buffer.index; i < c.buffer.size; i++ {
			if c.buffer.realBuffer[i-1] == '\r' && c.buffer.realBuffer[i] == '\n' {
				index := c.buffer.index
				c.buffer.index = i + 1
				if c.buffer.index >= c.buffer.size {
					c.mustRead = true
				}
				return c.buffer.realBuffer[index: i-1], nil
			}
		}
		err := c.ReadUnsafeBuffer()
		if err != nil {
			return nil, err
		}
	}
}

//in it the buffer.index point the ':'
//after this the buffer.index point the frist of next reply
//return int from the string
func (c *poolConn) ReadInt() (int, error) {
	for {
		for i := c.buffer.index; i < c.buffer.size; i++ {
			if c.buffer.realBuffer[i-1] == '\r' && c.buffer.realBuffer[i] == '\n' {
				index := c.buffer.index
				c.buffer.index = i + 1
				if c.buffer.index >= c.buffer.size {
					c.mustRead = true
				}
				return MyBstoI64(c.buffer.realBuffer[index: i-1])
			}
		}
		err := c.ReadUnsafeBuffer()
		if err != nil {
			return 0, err
		}
	}
}

//in it the buffer.index point after the '+'
//after this the buffer.index point the frist of next reply
func (c *poolConn) ReadOK() (bool, error) {
	for {
		if c.buffer.size-c.buffer.index >= 4 {
			if c.buffer.realBuffer[c.buffer.index] == 'O' && c.buffer.realBuffer[c.buffer.index+1] == 'K' && c.buffer.realBuffer[c.buffer.index+2] == '\r' && c.buffer.realBuffer[c.buffer.index+3] == '\n' {
				c.buffer.index = c.buffer.index + 4
				if c.buffer.index >= c.buffer.size {
					c.mustRead = true
				}
				return true, nil
			} else {
				return false, nil
			}
		} else {
			err := c.ReadUnsafeBuffer()
			if err != nil {
				return false, err
			}
		}
	}
}

//in it the buffer.index point after the '+'
//after this the buffer.index point the frist of next reply
func (c *poolConn) ReadPong() (bool, error) {
	for {
		if c.buffer.size-c.buffer.index >= 6 {
			if c.buffer.realBuffer[c.buffer.index] == 'P' && c.buffer.realBuffer[c.buffer.index+1] == 'O' && c.buffer.realBuffer[c.buffer.index+2] == 'N' && c.buffer.realBuffer[c.buffer.index+3] == 'G' && c.buffer.realBuffer[c.buffer.index+4] == '\r' && c.buffer.realBuffer[c.buffer.index+5] == '\n' {
				c.buffer.index = c.buffer.index + 6
				if c.buffer.index >= c.buffer.size {
					c.mustRead = true
				}
				return true, nil
			} else {
				return false, nil
			}
		} else {
			err := c.ReadUnsafeBuffer()
			if err != nil {
				return false, err
			}
		}
	}
}
