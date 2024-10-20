package src

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"reflect"
	"time"
)

type CacheKeyGenerator struct {
	functionName string
}

func (ckg *CacheKeyGenerator) NewCacheKeyGenerator(functionName string) *CacheKeyGenerator {
	return &CacheKeyGenerator{functionName: functionName}
}

func (ckg *CacheKeyGenerator) helperFlatten(o any) any {
	if o == nil {
		return "__nil__"
	}
	switch v := o.(type) {
	case string:
		return v
	case int, float64, bool:
		return fmt.Sprintf("%v", v)
	case time.Time:
		return v.Unix()
	case []interface{}:
		return ckg.flattenObject(v)
	default:
		val := reflect.ValueOf(o)
		if val.Kind() == reflect.Map {
			var result []interface{}
			for _, key := range val.MapKeys() {
				result = append(result, ckg.flattenObject(val.MapIndex(key).Interface()))
			}
			return result
		}
		if val.Kind() == reflect.Struct {
			var result []interface{}
			for i := 0; i < val.NumField(); i++ {
				result = append(result, ckg.flattenObject(val.Field(i).Interface()))
			}
			return result
		}
	}
	return fmt.Sprintf("%+v", o)
}

func (ckg *CacheKeyGenerator) flattenObject(o ...any) []interface{} {
	var result []interface{}
	for _, item := range o {
		switch v := item.(type) {
		case []interface{}:
			result = append(result, ckg.flattenObject(v))
			break
		default:
			result = append(result, ckg.helperFlatten(v))
		}
	}
	return result
}

func (ckg *CacheKeyGenerator) For(args ...any) (string, error) {
	var cacheKey string
	flattenedString := append([]interface{}{ckg.functionName}, ckg.flattenObject(args)...)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(flattenedString)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	cacheKey = string(hash.Sum(network.Bytes()))
	return cacheKey, nil
}
