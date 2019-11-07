package token

import (
	"bytes"
	"crypto/rand"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(2)
}

func TestMe(t *testing.T) {
	token := Generate()
	t.Log(token)

	if len(token) != TOKEN_SIZE {
		t.Fatal("invalid token size", len(token))
	}

	token2 := Generate()
	t.Log(token2)

	if len(token2) != len(token) {
		t.Fatal("tokens should be the same length")
	}

	if bytes.Compare(token, token2) == 0 {
		t.Fatal("Each token should be unique")
	}

	valid := Validate(token)

	if !valid {
		t.Fatal("token should be a valid one")
	}

	valid = Validate(token2)

	if !valid {
		t.Fatal("token should NOT be a valid one")
	}

	//random...
	fakeToken := make([]byte, 64)
	rand.Read(fakeToken)

	if Validate(fakeToken) {
		t.Fatal("faketoken should NOT be a valid one")
	}

}
