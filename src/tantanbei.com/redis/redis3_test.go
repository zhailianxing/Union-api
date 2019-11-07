package redis

//import (
//	"fmt"
//	"io/ioutil"
//	"os"
//	"testing"
//	"time"
//)

////func TestBigBytesLoop(T *testing.T) {
////	r, _ := NewRedis("127.0.0.1:6379", "", 500*time.Second)
////	r.FlushAll()

////	loopCount := 500000

////	for i := 0; i < loopCount; i++ {
////		f, err := os.Open("new0.jpg")
////		if err != nil {
////			T.Error(err)
////		}
////		defer f.Close()
////		bytes, err := ioutil.ReadAll(f)
////		if err != nil {
////			T.Error(i, err)
////		}
////		r.Set("pic", bytes, 0)
////		fmt.Println("GET")

////		bs, err := r.GetBytes("pic")
////		if len(bs) != len(bytes) || err != nil || bs[len(bs)-1] != bytes[len(bytes)-1] {
////			panic(err)
////		}
////		//ioutil.WriteFile(fmt.Sprintf("%d.jpg", a), bytes, 0666)
////	}
////}

//func TestPictureGoroutine(T *testing.T) {
//	r, _ := NewRedis("127.0.0.1:6379", "", 500*time.Second)
//	r.FlushAll()

//	loopCount := 100000
//	var cn = make(chan int, loopCount)

//	for i := 0; i < loopCount; i++ {
//		go func(a int) {
//			f, err := os.Open("new0.jpg")
//			if err != nil {
//				T.Error(err)
//			}
//			bytes, err := ioutil.ReadAll(f)
//			if err != nil {
//				T.Error(err)
//			}
//			r.Set("pic", bytes, 0)
//			bs, err := r.GetBytes("pic")
//			if len(bs) != len(bytes) || err != nil {
//				T.Error("picture error!", err)
//				fmt.Println("len(bs):", len(bs), "len(bytes):", len(bytes))
//			}
//			//ioutil.WriteFile(fmt.Sprintf("%d.jpg", a), bytes, 0666)
//			cn <- 1
//		}(i)
//	}
//	res := 0
//	for {
//		i := <-cn
//		res += i
//		if res >= loopCount {
//			break
//		}
//	}
//}
