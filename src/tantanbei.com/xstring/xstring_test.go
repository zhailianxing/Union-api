package xstring

import "testing"

func TestCheckIpAndPort(t *testing.T) {
	addr := "192.168.1.6:60344"
	ip, port := CheckIpAndPort(addr)
	if ip != "192.168.1.6" || port != "60344" {
		t.Fatal(ip, port)
	}
}
