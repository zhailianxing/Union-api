package xstring

func CheckIpAndPort(addr string) (ip, port string) {
	for i, s := range addr {
		if s == rune(':') {
			ip = addr[:i]
			port = addr[i+1:]
			return
		}
	}

	return
}
