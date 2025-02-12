package main

import "sync"

type KeyVal struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewKeyVal() *KeyVal {
	return &KeyVal{
		data: map[string][]byte{},
	}
}

func (KeyVal *KeyVal) Set(key, val []byte) error {
	KeyVal.mu.Lock()
	defer KeyVal.mu.Unlock()
	KeyVal.data[string(key)] = []byte(val)
	return nil
}

func (KeyVal *KeyVal) Get(key []byte) ([]byte, bool) {
	KeyVal.mu.RLock()
	defer KeyVal.mu.RUnlock()
	val, ok := KeyVal.data[string(key)]
	return val, ok
}
