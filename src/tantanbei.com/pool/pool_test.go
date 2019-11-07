package pool

import (
	"errors"
	"fmt"
	"net"
	//	"runtime"
	"strconv"
	"testing"
	"time"

	"sync/atomic"
)

//func init() {
//	runtime.GOMAXPROCS(2)
//}

type poolTestConn struct {
	realConn net.Conn
	err      error
}

//if connection has error, return error
func (c *poolTestConn) HealthCheck() (e error) {
	if c.realConn == nil {
		return errors.New("Test failed...null connection...")
	}

	e = c.Error()

	return
}

func (c *poolTestConn) Destroy() {
	return
}

func (c *poolTestConn) Error() error {
	return c.err
}

func (c *poolTestConn) Do(cmd string) {
	c.realConn.Write([]byte(cmd))
}

//test Get() timeout
func TestOverTime(t *testing.T) {
	p := NewPool("tcp "+"127.0.0.1:6379", 2, 2,
		func(args ...interface{}) (c PoolItem, e error) {
			con := poolTestConn{}
			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
			c = &con
			return
		})

	//debug only
	p.StartGatherStats()

	if p.Size() != 0 {
		t.Fatal("should be zero pool")
	}

	var cReal *poolTestConn

	c1, e := p.Get(time.Second * 2)
	//	p.PrintDebug()

	if e != nil {
		t.Fatal(e)
	}

	cReal = c1.(*poolTestConn)
	cReal.Do("PING")

	c2, err := p.Get(time.Second * 2)
	//	p.PrintDebug()

	if err != nil {
		t.Fatal(err)
	}
	cReal = c2.(*poolTestConn)
	cReal.Do("PING")

	if p.Size() != 2 {
		p.PrintDebug()
		t.Fatal("should be size 2 pool")
	}

	go func() {
		//over time
		c3, err := p.Get(time.Second * 2)
		if err == nil {
			t.Fatal(err)
		}
		return

		cReal = c3.(*poolTestConn)
		cReal.Do("PING")

		p.Put(c3)

		c4, er := p.Get(time.Second * 1)
		if er != nil {
			t.Fatal(er)
		}

		cReal = c4.(*poolTestConn)
		cReal.Do("PING")
	}()
	//after 5 s ,to release object
	<-time.After(time.Second * 5)
	p.Release(nil)

	if p.Size() != 1 {
		t.Fatal("alivecount not 1")
	}
}

//test put more than the len of Idlepool objects to Idlepool
func TestMaxIdlepool(t *testing.T) {
	p := NewPool("tcp "+"127.0.0.1:6379", 2, 2,
		func(args ...interface{}) (c PoolItem, e error) {
			con := poolTestConn{}
			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
			c = &con
			return
		})

	//debug only
	p.StartGatherStats()

	p.PrintDebug()

	if p.Size() != 0 {
		t.Fatal("should be zero pool")
	}

	var cReal *poolTestConn
	c1, e := p.Get(time.Second * 2)
	if e != nil {
		t.Fatal(e)
	}
	cReal = c1.(*poolTestConn)
	cReal.Do("PING")

	c2, err := p.Get(time.Second * 2)
	if err != nil {
		t.Fatal(err)

	}
	cReal = c2.(*poolTestConn)
	cReal.Do("PING")

	p.Put(c1)
	p.Put(c2)

	//	<-time.After(time.Second * 5)
	//	p.PrintDebug()
}

//let 100 goroutine to Get(),test performance
func TestMoreGoroutine(t *testing.T) {
	n := int64(2000000)
	a := int64(0)
	errN := int64(0)

	p := NewPool("tcp "+"127.0.0.1:6379", 10, 10,
		func(args ...interface{}) (c PoolItem, e error) {
			con := poolTestConn{}
			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
			c = &con
			return
		})

	if p.Size() != 0 {
		t.Fatal("should be zero size pool")
	}

	//debug only
	p.StartGatherStats()

	for i := int64(0); i < n; i++ {
		go func() {
			//			fmt.Println("enter go r")
			c, e := p.Get(time.Second * 5)
			//			fmt.Println("got c")
			if e != nil {
				//				fmt.Println("fatal error at ", a, e)
				//				panic(e)
				atomic.AddInt64(&errN, 1)
			} else {
				p.Put(c)
			}

			atomic.AddInt64(&a, 1)
			//			fmt.Println("put c"i
			//			fmt.Println(a)
		}()
	}

	//	fmt.Println("done...sleeping for 60")
	for i := int64(0); atomic.LoadInt64(&a) != n; i++ {
		fmt.Println("sec", i, "n", atomic.LoadInt64(&a), "err n", atomic.LoadInt64(&errN))
		p.PrintDebug()
		//		<-timer.After(time.Second * 1)
	}

	//	<-time.After(time.Second * 5)

	//	fmt.Println("end size of event chan:", len(p.OpeningEvent))
	//	p.PrintDebug()

}

//when object had been Release() , Get() can new a object
func TestRelease(t *testing.T) {
	ch := make(chan int)

	p := NewPool("tcp "+"127.0.0.1:6379", 2, 2,
		func(args ...interface{}) (c PoolItem, e error) {
			con := poolTestConn{}
			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
			c = &con
			return
		})

	if p.Size() != 0 {
		t.Fatal("should be zero size pool")
	}

	//debug only
	p.StartGatherStats()

	c1, e1 := p.Get(time.Second * 2)

	if e1 != nil {
		t.Fatal(e1)
	}
	_ = c1

	c2, e2 := p.Get(time.Second * 2)

	if e2 != nil {
		t.Fatal(e2)
	}
	_ = c2

	go func() {
		startime := time.Now().UnixNano() / 1000000

		c3, e3 := p.Get(time.Second * 10)
		if e3 != nil {
			t.Fatal(e3)
		}

		totaltime := time.Now().UnixNano()/1000000 - startime

		fmt.Println(totaltime)

		if totaltime > 5300 || totaltime < 4800 {

			panic("wrong time" + strconv.Itoa(int(totaltime)))
		}

		c3.(*poolTestConn).Do("PING")

		ch <- 1
	}()

	<-time.After(time.Second * 5)
	p.Release(nil)

	<-ch
}

func TestPoolReuse(t *testing.T) {
	p := NewPool("tcp "+"127.0.0.1:6379", 2, 2,
		func(args ...interface{}) (c PoolItem, e error) {
			con := poolTestConn{}
			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
			c = &con
			return
		})
	if p.Size() != 0 {
		t.Fatal("should be zero size pool")
	}

	//debug only
	p.StartGatherStats()

	for i := 0; i < 1000; i++ {
		var cReal *poolTestConn
		c1, e := p.Get(time.Second * 2)

		if e != nil {
			t.Fatal(e)
		}

		cReal = c1.(*poolTestConn)

		cReal.Do("PING")

		p.Put(c1)

		if i > 0 && p.Size() != 1 {
			t.Fatal("should have 1 item in l1 cache", p.Size())
		}

		c1, e = p.Get(time.Second * 2)
		if e != nil {
			t.Fatal(e)
		}

		cReal = c1.(*poolTestConn)

		cReal.Do("PING")

		p.Put(c1)

		c1, e = p.Get(time.Second * 2)
		if e != nil {
			t.Fatal(e)
		}

		cReal = c1.(*poolTestConn)

		cReal.Do("PING")

		p.Put(c1)

		if p.Size() != 1 {
			t.Fatal("should have 1 items in l1 cache", p.Size())
		}

	}

}

func TestPoolAfterIdleDeath(t *testing.T) {
	p := NewPool("tcp "+"127.0.0.1:6379", 2, 2,
		func(args ...interface{}) (c PoolItem, e error) {
			con := poolTestConn{}
			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
			c = &con
			return
		})

	if p.Size() != 0 {
		t.Fatal("should be zero size pool")
	}

	//debug only
	p.StartGatherStats()

	for i := 0; i < 2; i++ {
		if i > 0 {
			fmt.Println("about to sleep for a long time....")
			//sleep for a long time
			<-time.After(time.Second * 10)
		}

		var cReal *poolTestConn
		c1, e := p.Get(time.Second * 2)

		if e != nil {
			t.Fatal(e)
		}

		cReal = c1.(*poolTestConn)

		cReal.Do("PING")

		p.Put(c1)

		if i == 0 && p.Size() != 1 {
			t.Fatal("should have 1 item in l1 cache", p.Size())
		}

		c1, e = p.Get(time.Second * 2)
		if e != nil {
			t.Fatal(e)
		}

		cReal = c1.(*poolTestConn)

		cReal.Do("PING")

		p.Put(c1)

		if p.Size() != 1 {
			t.Fatal("should have 1 item in l1 cache", p.Size())
		}

	}

}

//func BenchmarkReUseNotTimed(b *testing.B) {
//	p := NewPool("tcp "+"127.0.0.1:6379",5, 5,
//		func(args ...interface{}) (c PoolItem, e error) {
//			con := poolTestConn{}
//			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
//			c = &con
//			return
//		})

//	if p.Size() != 0 {
//		b.Fatal("should be zero size pool")
//	}

//	c, e := p.Get(time.Second * 2)

//	if e != nil {
//		b.Fatal(e)
//	}

//	p.Put(c)

//	b.ResetTimer()

//	for i := 0; i < b.N; i++ {
//		c, e = p.Get(time.Second * 2)

//		if e != nil {
//			b.Fatal(e)
//		}

//		p.Put(c)
//	}

//}

//func BenchmarkReUseTimed(b *testing.B) {
//	p := NewPool("tcp "+"127.0.0.1:6379",5, 5,
//		func(args ...interface{}) (c PoolItem, e error) {
//			con := poolTestConn{}
//			con.realConn, e = net.Dial("tcp", "127.0.0.1:6379")
//			c = &con
//			return
//		})

//	if p.Size() != 0 {
//		b.Fatal("should be zero size pool")
//	}

//	c, e := p.Get(time.Second * 2)

//	if e != nil {
//		b.Fatal(e)
//	}

//	p.Put(c)

//	b.ResetTimer()

//	for i := 0; i < b.N; i++ {
//		c, e = p.Get(time.Second * 2)

//		if e != nil {
//			b.Fatal(e)
//		}

//		p.Put(c)
//	}

//}
