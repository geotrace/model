package model

import "testing"

func TestPassword(t *testing.T) {
	passwd := NewPassword("test")
	if !passwd.Compare("test") {
		t.Fatal("bad compare password")
	}
}
