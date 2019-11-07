package redis

import (
	"net"
	"time"
	"bufio"
)

type poolConn struct {
	realConn     net.Conn
	err          error
	cmdbuf       *bufio.Writer //send command buffer
	buffer       *ConnBuffer   //reply buffer
	writeTimeout time.Duration
	readTimeout  time.Duration
	mustRead     bool //mustRead==true: realBuffer has no usable data, need to read tcp
}
