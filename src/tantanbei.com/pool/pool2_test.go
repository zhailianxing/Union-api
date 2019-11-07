package pool

//import (
//	"fmt"
//	"net"
//	"runtime"
//	"testing"
//	"time"

//	"lbxtech.com/errorx"
//	"lbxtech.com/time/timer"
//)

//func init() {
//	runtime.GOMAXPROCS(2)
//}

//type poolTestConn struct {
//	realConn net.Conn
//	err      error
//}

////if connection has error, return error
//func (c *poolTestConn) HealthCheck() (e error) {
//	if c.realConn == nil {
//		return errorx.New("Test failed...null connection...")
//	}

//	e = c.Error()

//	return
//}

//func (c *poolTestConn) Destroy() {
//	return
//}

//func (c *poolTestConn) Error() error {
//	return c.err
//}

//func (c *poolTestConn) Do(cmd string) {
//	c.realConn.Write([]byte(cmd))
//}

//func TestMe(t *testing.T) {
//	p := NewPool(2, 2,
//		func(args ...interface{}) (c PoolItem, e error) {
//			con := poolTestConn{}
//			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
//			c = &con
//			return
//		})
//	if p.Size() != 0 {
//		t.Fatal("should be zero size pool")
//	}

//	loopCount := 8000 * 60

//	var cn = make(chan int, loopCount)
//	for i := 0; i < loopCount; i++ {
//		go func(a int) {

//			c, e := p.Get(5 * time.Second)
//			if e == nil {

//				<-timer.After(4 * time.Second)

//				p.Put(c)
//			}

//			//			fmt.Println("after put", a)
//			cn <- 1
//		}(i)

//	}
//	res := 0
//	fmt.Println("in wait res")
//	for {
//		i := <-cn

//		res += i
//		if res >= loopCount {
//			break
//		}
//		//		fmt.Println("index of i :", res)
//	}
//	fmt.Println("after  wait res:", res)
//}
