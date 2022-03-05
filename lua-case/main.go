package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

// SCRIPT LOAD "redis.call('SET',KEYS[1],KEYS[2]);redis.call('EXPIRE',KEYS[1],KEYS[3]);return 1;"
//EVALSHA  962e128406ab892b4fb7943255a7e1cf431c7a83 3 name jimmy 60
var script string = `
redis.call('SET',KEYS[1],KEYS[2])
redis.call('EXPIRE',KEYS[1],KEYS[3])
return 1
`
var luaHash string

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	luaHash, _ = Client.ScriptLoad(script).Result() //返回的脚本会产生一个sha1哈希值,下次用的时候可以直接使用这个值
	fmt.Println(luaHash)
	return
}

func useLuaHash() {
	n, err := Client.EvalSha(luaHash, []string{"name", "jimy", "60"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("结果", n, err)
}

func main() {
	useLuaHash()
}
