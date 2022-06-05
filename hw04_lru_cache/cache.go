package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	i, ok := l.items[key]
	el := cacheItem{key: key, value: value}
	if ok {
		i.Value = el
		l.queue.MoveToFront(i)
	} else {
		i = l.queue.PushFront(el)
		l.items[key] = i
		if l.queue.Len() > l.capacity {
			toRemove := l.queue.Back()
			val, match := toRemove.Value.(cacheItem)
			if match {
				l.queue.Remove(toRemove)
				delete(l.items, val.key)
			}
		}
	}
	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	i, ok := l.items[key]
	if !ok {
		return nil, false
	}
	val, match := i.Value.(cacheItem)
	if !match {
		return nil, false
	}
	l.queue.MoveToFront(i)
	return val.value, true
}

func (l *lruCache) Clear() {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
