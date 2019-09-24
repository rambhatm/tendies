package main

import (
	"fmt"
	"testing"
)

func TestAuth(t *testing.T) {
	a := NewAuthData("testUser", "correctpassword")

	if match, _ := a.IsMatch("wrongpassword"); match {
		t.Errorf("Password comparison/hashing error; should not match")
	}

	if match, cookie := a.IsMatch("correctpassword"); !match {
		t.Errorf("Password comparison/hashing error; should match")
	} else {
		fmt.Printf("%+v", cookie)
	}
}
