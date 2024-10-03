package adapters

import (
	"errors"
	"sync"
	"time"
)

type CacheInternalObj struct {
	Data   any
	Reject bool
}

type CacheObject struct {
	Key  string
	Data CacheInternalObj
	TTL  TTL
}

func NewCacheObject(key string, data CacheInternalObj, ttl TTL) *CacheObject {
	return &CacheObject{
		Key:  key,
		Data: data,
		TTL:  ttl,
	}
}

type MemoryCacheProvider struct {
	CacheProvider[CacheInternalObj]
	store       map[string]*CacheObject
	storesAsObj bool
	mu          sync.Mutex
}

func (m *MemoryCacheProvider) GetStoresAsObj() bool {
	return m.storesAsObj
}

func NewMemoryCacheProvider() {
	mcp := &MemoryCacheProvider{
		store:       make(map[string]*CacheObject),
		storesAsObj: true,
	}
	go mcp.cleanupExpiredItems()
}

func (mcp *MemoryCacheProvider) cleanupExpiredItems() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mcp.mu.Lock()
		for key, obj := range mcp.store {
			if int64(obj.TTL) < time.Now().Unix() {
				delete(mcp.store, key) // Remove expired item
			}
		}
		mcp.mu.Unlock()
	}
}

func (mcp *MemoryCacheProvider) Pipeline() (Pipeline[CacheInternalObj], error) {
	return nil, errors.New("Method not implemented.")
}

func (mcp *MemoryCacheProvider) Name() string {
	return "memory"
}

func (mcp *MemoryCacheProvider) Expire(key string, newTTLFromNow TTL) (int, error) {
	now := time.Now().Unix()
	object := mcp.store[key]
	if object != nil && (int64(object.TTL) > now) {
		object.TTL = TTL(now) + newTTLFromNow
		return 1, nil
	}
	return 0, nil
}

//TODO: add more functions
