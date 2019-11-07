package redis

import (
	"errors"
)

func (c *poolConn) HealthCheck() (e error) {
	if c.realConn == nil {
		return errors.New("Test failed...null connection...")
	}
	e = c.Error()

	return
}

func (c *poolConn) Destroy() {
	c.realConn = nil
	return
}

func (c *poolConn) Error() error {
	return c.err
}

func (c *poolConn) Init(capacity int) {
	if capacity < 13 {
		panic("the capacity is too small, At least 14!")
	}
	c.buffer.realBuffer = make([]byte, capacity)
	c.buffer.size = 0
	c.buffer.index = -1
	c.buffer.capacity = capacity
}
