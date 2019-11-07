package pool

import (
	"errors"
	"fmt"
	"strconv"

	"sync"
	"sync/atomic"
	"time"

	l4g "github.com/alecthomas/log4go"
)

// errorx.ERR_FULL is returned from a pool connection method (Do, Send,
// Receive, Flush, Err) when the maximum number of database connections in the
// pool has been reached.

// Pool maintains a pool of stuff.

//API usage:
//c := Pool.Get(timeout duration)
//defer Pool.Put(c)
//if error, set c.err to non-nil and HealthCheck() should return non-nil
//---- if while u you using c and detect it as bad, just call c.SetError() to non-null and Put() will know it
//it should not return it to the pool
type Pool struct {

	//descripe the pool connect info
	TAG string

	// Dial is an application supplied function for creating new connections.
	Create func(args ...interface{}) (PoolItem, error)

	// Maximum number of idle/pooled + pending/dialing connections in the pool.
	MaxAlive int64

	//pending active connections, connects()/dials()
	//NOTE: pending count is only used for statistics gathering and not use for any logic/decision making...
	PendingCount int64

	//maxidle is the actual pool size...can be different than max active...
	MaxIdlePoolSize int64

	//this value is used for predicated timeout...should be set to max throughput get/put per second
	MaxThroughputRate int64

	//GET/PUT counter, we record the number of PUT counters which we can use
	//to optimze further..value roughly equals to how fast a pool user is doing work...get/put equals 1 complete task
	PutCounter int64
	//we have a dedicated goroutine that will calculate put rate per second...
	PutRatePerSecond int64

	//connection in play, even if api called Pool.Get(), it would not be decremented
	//only decrement if connection is closed, has error, or dropped
	AliveCount     int64
	AliveCountLock sync.Mutex

	//OpeningEvent means either
	//a) we have a new item that just entered pool
	//or
	//b) pool has a new slot open that can hold a new item
	//slice of ready-to-use idle pool items...
	//LIFO..last in, first out...
	IdlePool     []PoolItem
	IdlePoolSize int64
	//protects the above 2 vars
	IdlePoolLock sync.Mutex

	//if item is released via Release(), aliveCount--, and we have a opening so waiters can be notified and new pool item can be created
	OpeningEvent chan struct{}

	//retry delay for new pool item creating go rnoutine that is a for loop that only exit on success
	CreateRetryDelay time.Duration

	WaitingCount int64
}

// NewPool is a convenience function for initializing a pool.
func NewPool(tag string, maxAlive int64, maxIdle int64, creater func(args ...interface{}) (PoolItem, error)) *Pool {
	if maxIdle > maxAlive {
		panic("Len of IdlePool more than MaxActive")
	}

	p := &Pool{
		TAG: tag,

		Create:          creater,
		MaxAlive:        maxAlive,
		MaxIdlePoolSize: maxIdle,

		IdlePool:          make([]PoolItem, maxIdle),
		OpeningEvent:      make(chan struct{}, maxIdle),
		MaxThroughputRate: 15000,
		CreateRetryDelay:  time.Millisecond * 1500,
	}

	return p
}

func (p *Pool) StartGatherStats() {
	//loop 1second to get put stats
	go p.GatherStats()
}

func (p *Pool) GatherStats() {
	lastPutCounter := int64(0)

	for {
		<-time.After(time.Millisecond * 500)

		//tmp var
		t := atomic.LoadInt64(&p.PutCounter)

		//caculate rate
		//check for zero exceptions
		if lastPutCounter == 0 || t == 0 {
			atomic.StoreInt64(&p.PutRatePerSecond, 0)

			//we record the last seond put rate
			lastPutCounter = t
			continue
		}

		//store new rate
		atomic.StoreInt64(&p.PutRatePerSecond, t-lastPutCounter)

		//we record the last seond put rate
		lastPutCounter = t
	}
}

func (p *Pool) PrintDebug() {
	//fmt.Println("PutRate:", atomic.LoadInt64(&p.PutRatePerSecond)*2, "PutCounter:", atomic.LoadInt64(&p.PutCounter), "TimerOnlyWait:", atomic.LoadInt64(&p.TimerOnlyWaitingCount), "WaitingCount", atomic.LoadInt64(&p.WaitingCount), "PendingCount:", atomic.LoadInt64(&p.PendingCount), "AliveCount:", atomic.LoadInt64(&p.AliveCount), "MaxActive:", p.MaxAlive)

	fmt.Println("PutRate:", atomic.LoadInt64(&p.PutRatePerSecond)*2, "PutCounter:", atomic.LoadInt64(&p.PutCounter), "WaitingCount", atomic.LoadInt64(&p.WaitingCount), "PendingCount:", atomic.LoadInt64(&p.PendingCount), "AliveCount:", atomic.LoadInt64(&p.AliveCount), "MaxActive:", p.MaxAlive)

	if atomic.LoadInt64(&p.AliveCount) > p.MaxAlive {
		//fmt.Println("xxxxxxxxxxxxxxxx p.AliveCount > max active")
		panic("alive count  > max alive")
	}
}

//safe but locking code..100% accurate
func (p *Pool) Size() int64 {
	p.AliveCountLock.Lock()
	defer p.AliveCountLock.Unlock()

	return p.AliveCount
}

//unsafe can be used for stats gathering where 100% accuracy is not required, no locks
func (p *Pool) UnsafeSize() int64 {
	return p.AliveCount
}

//release is called if you want to drop the pool item due an error or other reason
//it is like put() except it is just dropped...not returned to the pool
func (p *Pool) Release(c PoolItem) {
	//call c destroy in go routine
	if c != nil {
		go p.destoryItem(c)
	}

	p.AliveCountLock.Lock()
	//	fmt.Println("Release called, alivecount--")
	//	logx.D("---------------  Release called... old alive count:", p.AliveCount)

	//fire OpeningEvent event
	p.AliveCount--
	//logx.D("AliveCount--", p.AliveCount)

	//check make sure there is no bug if pool API is mis-used...this cannot happen internally
	//but will happen if you call pool.Release() by mistake outside of this package...
	if p.AliveCount < 0 {
		p.AliveCountLock.Unlock()
		panic("API user bug! AliveCount is below 0")
	}

	p.AliveCountLock.Unlock()

	//release call also generates an opening event because there is now a potential space for new poolitem creation
	if atomic.LoadInt64(&p.WaitingCount) > 0 {

		select {
		case p.OpeningEvent <- struct{}{}:
			//			break
			//	case p.IdlePool <- nil:
			//		//			break
		default:
			//do nothing..., no one is waiting...
			//panic("no one is listening to event...but should")
			//		fmt.Println("xxxxxxxxxx unexpected release event happened...")
		}
	}
}

//TODO: we need a GC like thread to do c.Destory..instead of spawning a go routine for each item..
//Destroy might be doing i/o and take a non-deterministic time....use go routine to protect
func (p *Pool) destoryItem(c PoolItem) {
	defer func() {
		if e := recover(); e != nil {
			l4g.Error("HandlePanic() recover failed, error:", e)
			//			l4g.Stack()
		}
	}()

	c.Destroy()
}

//FILO queue...
//attempto speed up get by using a L1 cache
func (p *Pool) getFromIdle() (c PoolItem) {

	for {

		p.IdlePoolLock.Lock()

		//WE MUST RE-CHECK HERE!!!!!..since between RUnloCk and Lock..value might be changed
		if p.IdlePoolSize == 0 {
			p.IdlePoolLock.Unlock()
			return
		}

		//finally safe to take from pool
		p.IdlePoolSize--

		//we need to take from front of the pool for FILO to distribute the workload instead of FIFO where 1 pool item can be hammered
		c = p.IdlePool[0]

		//shift pool items to left by 1...
		copy(p.IdlePool, p.IdlePool[1:])

		//nil the last position
		//we need to release the reference to pool item...otherwise, it would leak a pool item
		p.IdlePool[p.IdlePoolSize] = nil

		p.IdlePoolLock.Unlock()

		if c == nil {
			//fmt.Println(p.IdlePoolSize, p.IdlePool)
			panic("c from non-empty idlepool is null? wtf?" + strconv.FormatInt(p.IdlePoolSize, 10))
		}

		//WARNING: all pool item health checks must be completed in deterministic time
		if c.HealthCheck() == nil {
			//got healthy c
			return
		}

		//we need to release a unhealthy pool item from pool
		p.Release(c)

		//see if we have more to get...return to top of loop
	}

	return
}

// Get gets a connection. The application must close the returned connection.
// The connection acquires an underlying connection on the first call to the
// connection Do, Send, Receive, Flush or Err methods. An application can force
// the connection to acquire an underlying connection without executing a Redis
// command by calling the Err method.
// isReuse boolean (is the connection retured a re-cycled item or a new item) gives more detail to API caller...mysql driver can better manage error situations with this info
func (p *Pool) Get(timeout time.Duration) (c PoolItem, e error) {
	//user wants to wait forever..
	if timeout == 0 {
		panic("Please don't pass a zero time duration to Pool.Get()....")
	}

	//	fmt.Println("in get timed...")
	startTime := time.Now().UnixNano()
	timeoutNano := timeout.Nanoseconds()

	for {
		//got good item from idle pool
		if c = p.getFromIdle(); c != nil {
			//			logx.D("After Get(): Alive:", p.AliveCount, "IdleSize:", p.IdlePoolSize)
			return
		}

		//we have reached time limit...
		if time.Now().UnixNano()-startTime >= timeoutNano {
			e = errors.New("time out")
			return
		}

		//if pool is empty and AliveCounter is less than max...
		//AND pendingcount is zero...we can then create...otherwise a single Get() can
		//create more than 1 entry...
		p.AliveCountLock.Lock()

		if p.AliveCount < p.MaxAlive {
			//creating a new conn is counted as both alive and pending..
			p.AliveCount++

			atomic.AddInt64(&p.PendingCount, 1)

			//			fmt.Println("aliveCount++, pendingCount++")
			//			p.PrintDebug()

			//			logx.Stack()
			go p.createNewEntry()
		}

		p.AliveCountLock.Unlock()

		atomic.AddInt64(&p.WaitingCount, 1)

		nsTimeLeft := timeoutNano - (time.Now().UnixNano() - startTime)

		if nsTimeLeft <= 0 {
			//timedout..
			continue
		}

		//our soft-clock has only 100ms accuracy
		//minimum spin time is 100ms
		if nsTimeLeft < 100000000 {
			nsTimeLeft = 100000000
		}

		//wait for either one...
		select {
		case <-p.OpeningEvent:
		case <-time.After(time.Duration(nsTimeLeft) * time.Nanosecond):
		}

		atomic.AddInt64(&p.WaitingCount, -1)
	}
}

func (p *Pool) createNewEntry() (c PoolItem, e error) {
	//	logx.D("in create new...")
	for {

		c, e = p.Create()
		l4g.Debug(p.TAG, e)

		//put into cache and fire event
		if e == nil {
			//fmt.Println("dial success")
			atomic.AddInt64(&p.PendingCount, -1)
			//put into cache
			p.Put(c)
			//			fmt.Println("create finished")
			//			p.PrintDebug()
			return

		}

		//go and retry on failed create
		//wait 1s before retry to prevent infinite loop
		<-time.After(p.CreateRetryDelay)
	}

	return
}

// Close releases the resources used by the pool.
func (p *Pool) Close() error {
	//reset

	p.AliveCountLock.Lock()

	atomic.StoreInt64(&p.PendingCount, 0)
	p.AliveCount = 0

	p.AliveCountLock.Unlock()

	p.IdlePoolLock.Lock()
	p.IdlePoolSize = 0
	p.IdlePool = make([]PoolItem, p.MaxIdlePoolSize)
	p.IdlePoolLock.Unlock()

	close(p.OpeningEvent)
	p.OpeningEvent = make(chan struct{}, p.MaxIdlePoolSize)

	return nil
}

//return item to pool...
func (p *Pool) Put(c PoolItem) (e error) {
	//record healthy C put...
	atomic.AddInt64(&p.PutCounter, 1)

	//safety..
	if c == nil {
		//logx.D("put c == nil, release")
		p.Release(nil)
		return
	}

	//unhealhty c
	//WARNING: all pool item health checks must be completed in determinic time
	if c.HealthCheck() != nil {
		//logx.D("Not healthy put")

		p.Release(c)
		return
	}

	p.IdlePoolLock.Lock()

	//WE MUST RE-CHECK HERE!!!!!..since between RUnloCk and Lock..
	//there might have been a idle pool ADD event...
	//otherwise, we might miss a p.Release() call...
	if p.IdlePoolSize >= p.MaxIdlePoolSize {
		//overflow
		p.Release(c)

		//logx.D("idlePool put overflow..",p.IdlePoolSize)
	} else {
		//finally safe to put item back into idle pool...
		p.IdlePool[p.IdlePoolSize] = c
		p.IdlePoolSize++

	}
	p.IdlePoolLock.Unlock()

	if atomic.LoadInt64(&p.WaitingCount) > 0 {

		select {
		case p.OpeningEvent <- struct{}{}:
		default:

			//do nothing..., no one is waiting...
			//			fmt.Println("opening event...no one is waiting")
			//			panic("no one is listening to event...but should")
			//			fmt.Println("XXXXXXXXXXXXXX")
			//			p.PrintDebug()

			//TODO this should not happen, but we must make sure c is recycled
			//put item back into idle pool...
			//		p.IdlePool[p.IdlePoolSize] = c
			//		p.IdlePoolSize++
			//		fmt.Println("xxxxxxxxxxxxxx Unexpected put event happened...")
		}
	}

	return
}
