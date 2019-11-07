package pool

type GenericPoolItem struct {
}

//close connection
//if netConn, then this would equal to Close()
func (self *GenericPoolItem) Destroy() {

}

//get last error
func (self *GenericPoolItem) Error() error {
	return nil
}

//test to see if connection is alive..could as simple
//as testing if a realconn is = nil
//nil == health
func (self *GenericPoolItem) HealthCheck() error {
	return nil
}
