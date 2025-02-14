package main

import (
	"container/list"
	"sync"
	"time"
)

type KeyVal struct {
	mu           sync.RWMutex
	data         map[string][]byte
	lruList      *list.List
	keyToElement map[string]*list.Element
	memoryLimit  int64
	currentSize  int64
}

type lruItem struct {
	key   string
	value []byte
}

// memoryLimit should be in MB
func NewKeyVal(memoryLimit int64) *KeyVal {
	memoryLimit = memoryLimit * 1024 * 1024 //Size in mb
	kv := &KeyVal{
		data:         make(map[string][]byte),
		lruList:      list.New(),
		keyToElement: make(map[string]*list.Element),
		memoryLimit:  memoryLimit,
		currentSize:  0,
	}
	go kv.startEviction() // Start the background eviction goroutine
	return kv
}
func (kv *KeyVal) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	keyStr := string(key)

	if existingElem, exists := kv.keyToElement[keyStr]; exists {
		existingItem := existingElem.Value.(*lruItem)
		kv.currentSize -= int64(len(existingItem.value))
		existingItem.value = val
		kv.currentSize += int64(len(val))
		kv.lruList.MoveToFront(existingElem)
	} else {
		item := &lruItem{key: keyStr, value: val}
		element := kv.lruList.PushFront(item)
		kv.keyToElement[keyStr] = element
		kv.data[keyStr] = val
		kv.currentSize += int64(len(val))
	}
	kv.evictIfNeeded()
	return nil
}

func (kv *KeyVal) Get(key []byte) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	keyStr := string(key)
	if element, exists := kv.keyToElement[keyStr]; exists {
		kv.lruList.MoveToFront(element)
		return kv.data[keyStr], true
	}
	return nil, false
}

func (kv *KeyVal) evictIfNeeded() {
	for kv.currentSize > kv.memoryLimit {
		lastElem := kv.lruList.Back()
		if lastElem == nil {
			break
		}
		item := lastElem.Value.(*lruItem)
		delete(kv.data, item.key)
		kv.lruList.Remove(lastElem)
		delete(kv.keyToElement, item.key)
		kv.currentSize -= int64(len(item.value))
	}
}

func (kv *KeyVal) startEviction() {
	for {
		time.Sleep(10 * time.Second)
		kv.mu.Lock()
		kv.evictIfNeeded()
		kv.mu.Unlock()
	}
}
