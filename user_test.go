package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	u := NewUser("ram", "")

	err := InsertUserToDB(u)
	if err != nil {
		t.Errorf("insert error %s", err)
	}

}
