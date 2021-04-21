package db

import "testing"

func TestUser(t *testing.T) {
	u := User{1, "fens", "grs"}

	if u.Login() {
		t.Log("pass,u=", u)
	} else {
		t.Fatal("fail, u=", u)
	}
}
