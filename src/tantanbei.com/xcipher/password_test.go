package xcipher

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	result := HashPassword("diegov587")
	t.Log(result)
	if len(result) != 128 {
		t.Fatal("hash the password error")
	}
}
