package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/allegro/bigcache"
)

// Cache ...
type Cache interface {
	Set(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
	Delete(key interface{}) error
}

// bigCache ...
type bigCache struct {
	cache    *bigcache.BigCache
	DBGetter DBGettFunc // 缓存不存在返回的函数
}

type DBGettFunc func() interface{}

// Set 增加
func (c *bigCache) Set(key, value interface{}) error {
	// 断言key属于字符串类型
	keyString, ok := key.(string)
	if !ok {
		return errors.New("a cache key must be a string")
	}

	// 序列化value成字节
	valueBytes, err := serialize(value)
	if err != nil {
		return err
	}

	return c.cache.Set(keyString, valueBytes)
}

// Get 获取
func (c *bigCache) Get(key interface{}) (interface{}, error) {
	// 断言key属于字符串类型
	keyString, ok := key.(string)
	if !ok {
		return nil, errors.New("a cache key must be a string")
	}

	// 读取被存在字节序列里的value
	valueBytes, err := c.cache.Get(keyString)
	if err != nil {
		return nil, err
	}

	// 反序列化字节切片里的value
	value, err := deserialize(valueBytes)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Delete 删除
func (c *bigCache) Delete(key interface{}) error {
	// 断言key属于字符串类型
	keyString, ok := key.(string)
	if !ok {
		return errors.New("a cache key must be a string")
	}

	// 删除
	err := c.cache.Delete(keyString)
	if err != nil {
		return err
	}

	return nil
}

func serialize(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deserialize(valueBytes []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// New 新建
func New(expire time.Duration) (*bigCache, error) {

	var c bigCache
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(expire * time.Minute))
	if err != nil {
		return nil, errors.New("New bigcache is faild!")
	}
	c.cache = cache
	return &c, nil
}

func (c *bigCache) GetCache(key interface{}) error {
	obj := c.DBGetter()

	// 目前只要Get错误就重新设置缓存
	_, err := c.Get(key)
	if err != nil {
		err = c.Set(key, obj)
		if err != nil {
			return err
		}
	}

	return nil
}
