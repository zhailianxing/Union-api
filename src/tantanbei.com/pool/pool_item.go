package pool

type PoolItem interface {
	//close connection
	//if netConn, then this would equal to Close()
	Destroy()

	//get last error
	Error() error

	//test to see if connection is alive..could as simple
	//as testing if a realconn is = nil
	//nil == health
	HealthCheck() error
}
