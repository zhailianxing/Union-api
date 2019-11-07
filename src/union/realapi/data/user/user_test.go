package user

import "testing"

func TestAddNewUser(t *testing.T) {
	err := AddNewUser("18117541072", "paipable2016", "paipable")
	if err != nil {
		t.Fatal("err", err)
	}
}
