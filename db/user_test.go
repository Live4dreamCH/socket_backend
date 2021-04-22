package db

import "testing"

func TestUser(t *testing.T) {
	u := User{1, "fens", "grs"}

	if u.Login() {
		t.Log("Login pass,u=", u)
	} else {
		t.Fatal("Login fail, u=", u)
	}

	if u.HasFriend(1) {
		t.Log("HasFriend pass, u=", u)
	} else {
		t.Fatal("HasFriend fail, u=", u)
	}

	if suss, err := u.AddFriend(2); suss && err == nil {
		t.Log("AddFriend pass, u=", u)
	} else {
		t.Fatal("AddFriend fail, u=", u, "err=", err)
	}
}
