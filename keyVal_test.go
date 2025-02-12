package main

import (
	"testing"
)

func TestKeyVal(t *testing.T) {
	KeyVal := NewKeyVal()
	key := "mykey"
	val := "myval"
	err := KeyVal.Set([]byte(key), []byte(val))
	if err != nil {
		t.Fatal(err)
	}
	val2, ok := KeyVal.Get([]byte(key))
	if !ok {
		t.Fatal(err)
	}
	if string(val2) != val {
		t.Fatal(err)
	}
}
