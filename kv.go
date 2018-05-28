package kv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// 默认是占用1MB内存
	defaultMaxSize int64 = 1024 * 1024
	// 默认key的大小
	defaultMaxKeyLen int64 = 64
	// 默认value的大小
	defaultMaxValueLen int64 = 1024
	// 默认缓存时长
	defaultExpire time.Duration = 0
)

var (
	// ErrKeyInValid is the key len gt defaultMaxKeyLen
	ErrKeyInValid = errors.New(fmt.Errorf("The key len is great than %d", defaultMaxKeyLen).Error())
	// ErrValueInvalid is the value is greate than max size
	ErrValueInvalid = errors.New("The value size is great than max size")
)

// cache 支持过期时间和最大内存大小的内存缓存库
type cache interface {
	// size 是一个字符。支持以下参数：1KB，1MB，2MB，1GB 等
	SetMaxMemory(size string) bool
	// 设置一个缓冲项，并且在expire时间之后过期
	Set(key string, val string, expire ...time.Duration) (bool, error)
	// 获取一个值
	Get(key string) (string, bool)
	// 删除一个值
	Del(key string) bool
	// 检测一个值是否存在
	Exists(key string) bool
	// 清除所有值
	Flush() bool
	// 返回所有key的数量
	Keys() int64
}

// Cache 对Cache的实现
type Cache struct {
	cache
	// 内存大小
	maxMemory int64
	// 当前内存大小
	currentMemory int64
	// 读写锁
	lock sync.RWMutex
	// lru
	lruList *LruList
}

// New 初始化
func New(args ...string) *Cache {

	var memory = defaultMaxSize

	if len(args) == 1 {
		realSize, err := parseSizeStr(args[0])

		if err == nil {
			memory = realSize
		}
	}

	lruList := NewLruList()

	return &Cache{
		maxMemory:     memory,
		currentMemory: 0,
		lruList:       lruList,
	}
}

// SetMaxMemory 设置最大内存
func (c *Cache) SetMaxMemory(size string) bool {

	realSize, err := parseSizeStr(size)

	if err != nil {
		return false
	}

	c.maxMemory = realSize

	return true
}

// Set 设置键值对
func (c *Cache) Set(key string, val string, args ...time.Duration) (bool, error) {

	var expire time.Duration
	var expireTime int64

	if len(args) == 1 {
		expire = args[0]
	}

	if expire != defaultExpire {
		expireTime = c.getCurrentTimestamp() + expire.Nanoseconds()
	}

	keySize := c.getSize(key)

	if keySize > defaultMaxKeyLen {
		return false, ErrKeyInValid
	}

	valSize := c.getSize(val)

	if valSize > defaultMaxValueLen {
		return false, ErrValueInvalid
	}

	c.currentMemory = c.currentMemory + keySize + valSize

	c.lock.Lock()

	c.lruList.Set(key, &Item{
		value:      val,
		expireTime: expireTime,
	})

	c.lock.Unlock()

	// 触发lru
	if c.currentMemory >= c.maxMemory {
		c.lru()
	}

	return true, nil
}

// Get 取值
func (c *Cache) Get(key string) (string, bool) {

	keySize := c.getSize(key)

	if keySize > defaultMaxKeyLen {
		return "", false
	}

	c.lock.Lock()

	item, ok := c.lruList.items[key]

	if !ok {

		c.lock.Unlock()

		return "", false
	}

	if item.expireTime != 0 && item.expireTime < c.getCurrentTimestamp() {

		c.lock.Unlock()

		return "", false
	}

	c.lock.Unlock()

	return item.value, true
}

// Del 删除值
func (c *Cache) Del(key string) bool {

	c.lock.Lock()

	exists := c.lruList.Exists(key)

	if exists {
		c.lruList.Del(key)

		keySize := c.getSize(key)
		valSize := c.getSize(c.lruList.items[key].value)
		c.currentMemory = c.currentMemory - keySize - valSize
	}

	c.lock.Unlock()

	return exists
}

// Exists 判断key是否存在
func (c *Cache) Exists(key string) bool {

	c.lock.Lock()

	exists := c.lruList.Exists(key)

	c.lock.Unlock()

	return exists
}

// Flush 清空所有key
func (c *Cache) Flush() bool {

	c.lock.Lock()

	c.lruList = NewLruList()
	c.currentMemory = 0

	c.lock.Unlock()

	return true
}

// Keys 返回key的数量
func (c *Cache) Keys() int64 {

	c.lock.Lock()

	n := c.lruList.Size()

	c.lock.Unlock()

	return int64(n)
}

// 转换size字符串为对应的size大小。1KB，1MB，2MB，1GB 等
func parseSizeStr(sizeStr string) (int64, error) {

	sizeType := string([]byte(strings.ToUpper(sizeStr)[len(sizeStr)-2:]))
	sizeNum := string([]byte(sizeStr[0 : len(sizeStr)-2]))

	num, err := strconv.ParseInt(sizeNum, 10, 64)

	if err != nil {
		return 0, err
	}

	switch sizeType {
	case "KB":
		return num * 1024, nil
	case "MB":
		return num * 1024 * 1024, nil
	case "GB":
		return num * 1024 * 1024 * 1024, nil
	}

	return defaultMaxSize, nil
}

// get current timstamp
func (c *Cache) getCurrentTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// get variable size
func (c *Cache) getSize(val string) int64 {
	return int64(len(val))
}

// removeTailNode 触发lru，删除尾部节点
func (c *Cache) lru() {

	for c.currentMemory >= c.maxMemory {

		c.lock.Lock()

		node := c.lruList.RemoveTailNode()

		if node != nil {
			keySize := c.getSize(node.key)
			valSize := c.getSize(node.data.value)

			c.currentMemory = c.currentMemory - keySize - valSize
		}

		c.lock.Unlock()
	}

}
