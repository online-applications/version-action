package utils

import (
	"testing"
	"os"
)

func TestGetEnvPositive(t *testing.T) {
	os.Setenv("TestEnv", "123")
	got := GetEnv("TestEnv")
    want := "123"

    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
	
	got_two := GetEnv("NotExist")
    want_two := ""

    if got_two != want_two {
        t.Errorf("got %q, wanted %q", got, want)
    }
	
}