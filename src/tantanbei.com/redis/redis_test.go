package redis

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

func dial(T *testing.T) *Redis {
	r, err := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	if err != nil {
		T.Error("New redis error!")
		return nil
	}
	return r
}

func TestGetString(T *testing.T) {
	r := dial(T)
	r.Del("tan")
	s, _ := r.GetString("tan")
	if s != "" {
		T.Error("GetString error")
	}
}

func BenchmarkGetString(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.GetString("tan")
	}
}

func TestGetBytes(T *testing.T) {
	r := dial(T)
	r.Set("tan", []byte("tan"), 0)
	b, _ := r.GetBytes("tan")
	if b[0] != 't' {
		T.Error("GetBytes error")
	}
}

func TestGetBigBytes(T *testing.T) {
	r := dial(T)

	f, _ := os.Open("new0.jpg")
	bytes, _ := ioutil.ReadAll(f)
	r.Set("pic", bytes, 0)
	b, _ := r.GetBytes("pic")
	if b[100] != bytes[100] {
		T.Error("bigbytes get error")
	}
}

func BenchmarkGetBytes(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("tan", []byte("tan"), 0)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.GetBytes("tan")
	}
}

func BenchmarkGetPicture(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	f, _ := os.Open("new0.jpg")
	bytes, _ := ioutil.ReadAll(f)
	r.Set("pic", bytes, 0)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := r.GetBytes("pic")
		if err != nil {
			panic(err)
		}
	}
}

func TestGetInt(T *testing.T) {
	r := dial(T)
	r.Set("abc", []byte("123"), 0)
	i, _ := r.GetInt64("abc")
	if i != 123 {
		T.Error("GetInt error")
	}
}

func BenchmarkGetInt64(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("123"), 0)
	for i := 0; i < b.N; i++ {
		n, _ := r.GetInt64("abc")
		if n != 123 {
			b.Error("GetInt error")
		}
	}
}

func BenchmarkGetInt32(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("123"), 0)
	for i := 0; i < b.N; i++ {
		n, _ := r.GetInt32("abc")
		if n != 123 {
			b.Error("GetInt error")
		}
	}
}

func BenchmarkGetInt16(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("123"), 0)
	for i := 0; i < b.N; i++ {
		n, _ := r.GetInt16("abc")
		if n != 123 {
			b.Error("GetInt error")
		}
	}
}

func BenchmarkGetInt8(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("123"), 0)
	for i := 0; i < b.N; i++ {
		n, _ := r.GetInt8("abc")
		if n != 123 {
			b.Error("GetInt error")
		}
	}
}

func TestSet(T *testing.T) {
	r := dial(T)
	err := r.Set("aaa", []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), 0)
	if err == nil {
		a, _ := r.GetString("aaa")
		if a != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
			T.Error("GetString get data error")
		}
	}
}

func BenchmarkSet(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	for i := 0; i < b.N; i++ {
		r.Set("tan", []byte("tan"), 0)
	}
}

func BenchmarkSetLarge(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	data := make([]byte, 99)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Set("aaa", data, 0)
	}
}

func BenchmarkSetPicture(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	f, _ := os.Open("new0.jpg")
	bytes, _ := ioutil.ReadAll(f)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Set("pic", bytes, 0)
	}
}

func TestSetExpire(T *testing.T) {
	r := dial(T)

	r.Set("de", []byte("de"), 3)
	de, _ := r.GetString("de")
	if de != "de" {
		T.Error("SetEx get data error")
	}
	<-time.After(4 * time.Second)
	de, err := r.GetString("de")
	if de != "" || err == nil {
		T.Error("SetEx get data error")
	}
}

func TestDel(T *testing.T) {
	r := dial(T)
	r.Set("id11", []byte("id11"), 0)
	r.Set("id22", []byte("id22"), 0)
	i, _ := r.Del("id11", "id22", "idid", "idip", "iddd")
	if i != 2 {
		T.Error("Del error")
	}
}

func BenchmarkDel(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	for i := 0; i < b.N; i++ {
		r.Del("id1")
	}
}

func TestMGet(T *testing.T) {
	r := dial(T)
	r.Set("id11", []byte("id11"), 0)
	r.Set("id22", []byte("id22"), 0)
	as, _ := r.MGet("id11", "id22", "id33")
	if len(as) != 3 {
		T.Error("MGet error")
	}
}

func BenchmarkMGet(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	for i := 0; i < b.N; i++ {
		r.MGet("id11", "id22", "id33")
	}
}

func TestGo(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 500*time.Second, 5*time.Second, 5*time.Second)

	loopCount := 100
	var cn = make(chan int, loopCount)

	for i := 0; i < loopCount; i++ {
		go func(a int) {
			s := fmt.Sprintf("id%d", a)
			r.Set(s, []byte(s), 0)
			ss, err := r.GetString(s)
			if ss != s {
				T.Error("Go error", err, a)
				//panic("go error")
			}
			cn <- 1
		}(i)
	}
	res := 0
	for {
		i := <-cn
		res += i
		if res >= loopCount {
			break
		}
	}
}

func TestIncr(T *testing.T) {
	r := dial(T)
	r.Set("abc", []byte("10"), 0)
	i, _ := r.Incr("abc")
	if i != 11 {
		T.Error("Incr error")
	}
	i, _ = r.Decr("abc")
	if i != 10 {
		T.Error("Incr error")
	}
}

func BenchmarkIncr(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("0"), 0)
	for i := 0; i < b.N; i++ {
		r.Incr("abc")
	}
}

func TestRename(T *testing.T) {
	r := dial(T)
	r.Set("abc", []byte("10"), 0)
	r.Set("aaa", []byte("10"), 0)
	r.Del("abc1")
	i, _ := r.Rename("abc", "abc1", true)
	if i != 1 {
		T.Error("rename error")
	}
	i, _ = r.Rename("aaa", "abc1", true)
	if i != 0 {
		T.Error("rename2 error")
	}
}

func BenchmarkRename(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("0"), 0)
	for i := 0; i < b.N; i++ {
		r.Rename("abc", "abc1", true)
	}
}

func TestDbsize(T *testing.T) {
	r := dial(T)
	r.Set("abc", []byte("10"), 0)
	r.Set("aaa", []byte("10"), 0)
	i, _ := r.Dbsize()
	if i < 2 {
		T.Error("dbsize error!", i)
	}
}

func BenchmarkDbsize(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	for i := 0; i < b.N; i++ {
		r.Dbsize()
	}
}

func TestGetSet(T *testing.T) {
	r := dial(T)
	r.Del("id")
	s, _ := r.GetSet("id", []byte("id"))
	if s != "" {
		T.Error("getset error!")
	}
	s, _ = r.GetSet("id", []byte("id1"))
	if s != "id" {
		T.Error("getset2 error!")
	}
	s, err := r.GetString("id")
	if err != nil {
		fmt.Println(err)
	}
	if s != "id1" {
		T.Error("getset3 error!")
	}
}

func BenchmarkGetSet(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	for i := 0; i < b.N; i++ {
		r.GetSet("id", []byte("id"))
	}
}

func TestIncrby(T *testing.T) {
	r := dial(T)
	r.Set("abc", []byte("10"), 0)
	i, _ := r.Incrby("abc", 10)
	if i != 20 {
		T.Error("Incr error")
	}
	i, _ = r.Decrby("abc", 10)
	if i != 10 {
		T.Error("Incr error")
	}
}

func BenchmarkIncrby(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("abc", []byte("0"), 0)
	for i := 0; i < b.N; i++ {
		r.Incrby("abc", 100)
	}
}

func TestAppend(T *testing.T) {
	r := dial(T)
	r.Set("id", []byte("id"), 0)
	r.Append("id", []byte("0105"))
	s, err := r.GetString("id")
	if err != nil {
		fmt.Println(err)
	}
	if s != "id0105" {
		T.Error("Append error!")
	}
}

func BenchmarkAppend(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("id", []byte("0"), 0)
	for i := 0; i < b.N; i++ {
		r.Append("id", []byte("0105"))
	}
}

func TestPush(T *testing.T) {
	r := dial(T)
	r.Push("list", []byte("one"), true)
	s, _ := r.Lrange("list", 0, -1)
	if s[0] != "one" {
		T.Error("Rpush error!")
	}
	s, err := r.Lrange("my", 0, -1)
	if err == nil {
		T.Error("Rpush2 error!")
	}
}

func BenchmarkPush(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("id", []byte("0"), 0)
	for i := 0; i < b.N; i++ {
		r.Push("list", []byte("two"), true)
	}
}

func TestPop(T *testing.T) {
	r := dial(T)
	r.Sadd("set", []byte("one"))
	s, _ := r.Spop("set")
	if s != "one" {
		T.Error("pop error!")
	}
	s, _ = r.Spop("set")
	if s != "" {
		T.Error("pop2 error")
	}
}

func BenchmarkPop(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	r.Set("id", []byte("0"), 0)
	for i := 0; i < b.N; i++ {
		r.Spop("set")
	}
}

func TestPing(T *testing.T) {
	r := dial(T)
	err := r.Ping()
	if err != nil {
		T.Error(err)
	}
}

func BenchmarkPing(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)
	for i := 0; i < b.N; i++ {
		r.Ping()
	}
}

func TestPictureGoroutine(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 500*time.Second, 5*time.Second, 5*time.Second)

	loopCount := 100
	var cn = make(chan int, loopCount)

	for i := 0; i < loopCount; i++ {
		go func(a int) {
			f, err := os.Open("new0.jpg")
			if err != nil {
				T.Error(err)
			}
			bytes, err := ioutil.ReadAll(f)
			if err != nil {
				T.Error(err)
			}
			r.Set("pic", bytes, 0)
			bs, err := r.GetBytes("pic")
			if len(bs) != len(bytes) || err != nil {
				T.Error("picture error!", err)
			}
			//ioutil.WriteFile(fmt.Sprintf("%d.jpg", a), bytes, 0666)
			cn <- 1
		}(i)
	}
	res := 0
	for {
		i := <-cn
		res += i
		if res >= loopCount {
			break
		}
	}
}

func BenchmarkLargeToSmall(b *testing.B) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	f, _ := os.Open("new0.jpg")
	bytes, _ := ioutil.ReadAll(f)
	for i := 0; i < 100; i++ {
		//go func() {
		r.Set("pic", bytes, 0)
		//}()
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		//r.Set("pic", bytes, 0)
		r.Set("tan", []byte("tan"), 0)
		//r.GetBytes("tan")
	}
}

func TestLargeToSmall(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	f, _ := os.Open("new0.jpg")
	bytes, _ := ioutil.ReadAll(f)
	for i := 0; i < 10; i++ {
		err := r.Set("pic", bytes, 0)
		if err != nil {
			T.Error(err)
		}
	}

	for i := 0; i < 10; i++ {
		err := r.Set("tan", []byte("tan"), 0)
		if err != nil {
			T.Error(err)
		}
		tan, err := r.GetBytes("tan")
		if tan[0] != 't' || err != nil {
			T.Error(err)
		}
	}
	for i := 0; i < 10; i++ {
		err := r.Set("pic", bytes, 0)
		if err != nil {
			T.Error(err)
		}
	}
}

func TestForLoop(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	loopCount := 1

	for i := 0; i < loopCount; i++ {
		r.Set("tan", []byte("tan"), 0)
		r.GetBytes("tan")
	}
}

func TestBigBytesLoop(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	loopCount := 100000

	for i := 0; i < loopCount; i++ {
		f, err := os.Open("111.jpg")
		if err != nil {
			T.Error(i, err)
		}
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			T.Error(i, err)
		}
		r.Set("pic", bytes, 0)
		bs, err := r.GetBytes("pic")
		if len(bs) != len(bytes) || err != nil || bs[len(bs)-1] != bytes[len(bytes)-1] {
			T.Error("picture error!", err, i)
		}
		//err = ioutil.WriteFile(fmt.Sprintf("%d.jpg", i), bs, 0666)
		if err != nil {
			T.Error("create picture error!")
		}
	}
}

func TestTxt(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	f, err := os.Open("a.txt")
	if err != nil {
		T.Error(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		T.Error(err)
	}
	r.Set("txt", bytes, 0)
	bs, err := r.GetBytes("txt")
	if len(bs) != len(bytes) || err != nil || bs[len(bs)-1] != bytes[len(bytes)-1] {
		T.Error("txt error!", err)
	}
	ioutil.WriteFile("b.txt", bytes, 0666)
}

func TestZadd(T *testing.T) {
	r, _ := NewRedis("127.0.0.1:6379", "", 1024, 5*time.Second, 5*time.Second, 5*time.Second)

	for i := 0; i < 10000; i += 100 {
		_, err := r.Zadd("paipai", strconv.Itoa(79400+i), strconv.Itoa(int(time.Now().UnixNano())))
		if err != nil {
			T.Fatal(err)
		}
	}
}

// go 1.5
//BenchmarkGetString-4   	  300000	     26509 ns/op	     352 B/op	       4 allocs/op
//BenchmarkGetBytes-4    	  300000	     27103 ns/op	     256 B/op	       4 allocs/op
//BenchmarkGetPicture-4  	  200000	     37334 ns/op	   16880 B/op	       4 allocs/op
//BenchmarkGetInt-4      	  300000	     26401 ns/op	     240 B/op	       3 allocs/op
//BenchmarkSet-4         	  300000	     27869 ns/op	     352 B/op	       6 allocs/op
//BenchmarkSetPicture-4  	  200000	     40477 ns/op	     336 B/op	       5 allocs/op
//BenchmarkDel-4         	  300000	     26298 ns/op	     240 B/op	       3 allocs/op
//BenchmarkMGet-4        	  200000	     29497 ns/op	     792 B/op	      13 allocs/op
//BenchmarkIncr-4        	  300000	     26860 ns/op	     352 B/op	       4 allocs/op
//BenchmarkRename-4      	  300000	     28562 ns/op	     416 B/op	       9 allocs/op
//BenchmarkDbsize-4      	  300000	     25593 ns/op	     160 B/op	       2 allocs/op
//BenchmarkGetSet-4      	  300000	     28233 ns/op	     496 B/op	       9 allocs/op
//BenchmarkIncrby-4      	  300000	     28050 ns/op	     448 B/op	       6 allocs/op
//BenchmarkAppend-4      	  300000	     27424 ns/op	     448 B/op	       6 allocs/op
//BenchmarkPush-4        	  300000	     27793 ns/op	     448 B/op	       6 allocs/op
//BenchmarkPop-4         	  300000	     26347 ns/op	     352 B/op	       4 allocs/op
//BenchmarkPing-4        	  300000	     25728 ns/op	     160 B/op	       2 allocs/op
//BenchmarkLargeToSmall-4	  300000	     27589 ns/op	     352 B/op	       6 allocs/op

// go 1.4
//BenchmarkGetString	  500000	     16623 ns/op	     368 B/op	       5 allocs/op
//BenchmarkGetBytes	  500000	     16662 ns/op	     264 B/op	       5 allocs/op
//BenchmarkGetPicture	  200000	     34345 ns/op	   16896 B/op	       6 allocs/op
//BenchmarkGetInt	  500000	     16497 ns/op	     256 B/op	       4 allocs/op
//BenchmarkSet	  500000	     17444 ns/op	     352 B/op	       7 allocs/op
//BenchmarkSetPicture	  200000	     32442 ns/op	     344 B/op	       7 allocs/op
//BenchmarkDel	  500000	     16199 ns/op	     256 B/op	       4 allocs/op
//BenchmarkMGet	  500000	     18509 ns/op	     808 B/op	      14 allocs/op
//BenchmarkIncr	  500000	     16831 ns/op	     368 B/op	       5 allocs/op
//BenchmarkRename	  500000	     17429 ns/op	     424 B/op	      10 allocs/op
//BenchmarkDbsize	  500000	     15829 ns/op	     176 B/op	       3 allocs/op
//BenchmarkGetSet	  500000	     17752 ns/op	     488 B/op	      11 allocs/op
//BenchmarkIncrby	  500000	     17610 ns/op	     456 B/op	       7 allocs/op
//BenchmarkAppend	  500000	     17108 ns/op	     464 B/op	       8 allocs/op
//BenchmarkPush	  500000	     17423 ns/op	     464 B/op	       8 allocs/op
//BenchmarkPop	  500000	     16613 ns/op	     368 B/op	       5 allocs/op
//BenchmarkPing	  500000	     15635 ns/op	     176 B/op	       3 allocs/op
//BenchmarkLargeToSmall	  500000	     17640 ns/op	     352 B/op	       7 allocs/op
