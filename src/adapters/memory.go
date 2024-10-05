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
	Key  CacheKey
	Data CacheInternalObj
	TTL  TTL
}

func NewCacheObject(key CacheKey, data CacheInternalObj, ttl TTL) *CacheObject {
	return &CacheObject{
		Key:  key,
		Data: data,
		TTL:  ttl,
	}
}

type MemoryCacheProvider struct {
	CacheProvider[CacheInternalObj]
	store       map[CacheKey]*CacheObject
	storesAsObj bool
	mu          sync.Mutex
}

func (mcp *MemoryCacheProvider) GetStoresAsObj() bool {
	return mcp.storesAsObj
}

func NewMemoryCacheProvider() {
	mcp := &MemoryCacheProvider{
		store:       make(map[CacheKey]*CacheObject),
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
	return nil, errors.New("method not implemented")
}

func (mcp *MemoryCacheProvider) Name() string {
	return "memory"
}

func (mcp *MemoryCacheProvider) Expire(key CacheKey, newTTLFromNow TTL) (int, error) {
	now := time.Now().Unix()
	object := mcp.store[key]
	if object != nil && (int64(object.TTL) > now) {
		object.TTL = TTL(now) + newTTLFromNow
		return 1, nil
	}
	return 0, nil
}

func (mcp *MemoryCacheProvider) Del(keys ...CacheKey) (int, error) {
	for _, key := range keys {
		delete(mcp.store, key)
	}
	return len(keys), nil
}

func (mcp *MemoryCacheProvider) Set(key CacheKey, data CacheInternalObj, ttl TTL) any {
	mcp.store[key] = NewCacheObject(key, data, TTL(time.Now().Unix())+ttl)
	return data
}

func (mcp *MemoryCacheProvider) Get(key CacheKey) (any, error) {
	object := mcp.store[key]
	if object != nil {
		if object.TTL < TTL(time.Now().Unix()) {
			_, err := mcp.Del(key)
			if err != nil {
				return nil, err
			}
			return nil, errors.New("data expired")
		}
		return object.Data, nil
	} else {
		return nil, errors.New("data unavailable")
	}
}

func (mcp *MemoryCacheProvider) Dump() (map[CacheKey]*CacheObject, error) {
	return mcp.store, nil
}

func (mcp *MemoryCacheProvider) FlushDB() {
	clear(mcp.store)
}

func (mcp *MemoryCacheProvider) MGet(keys ...CacheKey) ([]interface{}, error) {
	dataCh := make(chan any)
	go func() {
		for _, key := range keys {
			object, _ := mcp.Get(key)
			dataCh <- object
		}
		close(dataCh)
	}()
	var objects []interface{}
	isArrayOfNils := true
	for object := range dataCh {
		if isArrayOfNils && object != nil {
			isArrayOfNils = false
		}
		objects = append(objects, object)
	}
	if isArrayOfNils {
		return nil, errors.New("no objects returned")
	}
	return objects, nil
}

func (mcp *MemoryCacheProvider) MSet(kvPairs ...struct {
	CacheKey
	CacheInternalObj
}) (string, error) {
	if kvPairs == nil {
		return "", errors.New("no data supplied")
	}
	for _, kvPair := range kvPairs {
		//hardcoded for 5 minutes
		mcp.Set(kvPair.CacheKey, kvPair.CacheInternalObj, 300)
	}

	return "OK", nil
}
