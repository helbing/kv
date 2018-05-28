## kv

kv是一个轻量级的类memcached的小型类库，和memcached一样key和value均为string类型

## Usage

### 初始化

```go
cache := kv.New()
```

如果在初始化没有指定占用的内存大小，那么默认为1048576个字节，如果超过就会触发lru

带有内存大小的初始化

```go
cache := kv.New("1GB")
```

### 添加缓存

```go
cache.Set(key, value)
```

如果没有指定过期时间，那么默认就是永久缓存

带有过期时间(单位：毫秒)的缓存

```go
cache.Set(key, value, 1000)
```

### 获取缓存

```go
value, ok := cache.Get(key)

if ok {
    fmt.Println(value)
}
```

### 删除缓存

```go
cache.Del(key)
```
