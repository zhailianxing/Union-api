package security

import (
	"union/realapi/share"
	"fmt"
)

func addIp(ip string) (count int) {
	var err error
	if count, err = share.Redis.Incr(ip); err != nil {
		panic(err)
	}

	i, err := share.Redis.Expire(ip, 86400)
	fmt.Println(i, err)
	return
}

func CheckIpValid(ip string, maxCount int) bool {
	//do not check this error, the key maybe not exist
	count, _ := share.Redis.GetInt64(ip)
	if count >= maxCount {
		return false
	}

	addIp(ip)
	return true
}
