package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"sync"
)

const orderSet = "orderSet"    //用户id的集合
const goodsTotal = "goodTotal" //商品库存的key
const orderList = "orderList"  //订单队列
func createScript() *redis.Script {
	script := redis.NewScript(`
		local userId    = tostring(KEYS[1])
		local orderSet=tostring(KEYS[2])

-- 是否已经抢购到了,如果是返回
		local hasBuy = redis.call("sIsMember", orderSet, userId)
		if hasBuy ~= 0 then
		  return 0
		end

-- 库存的数量
		local goodsTotal=tonumber(ARGV[1])
		local total=redis.call("GET", goodsTotal)

-- 是否已经没有库存了,如果是返回
		if total <= 0 then
		  return 0
		end

-- 可以下单 
		local flag

-- 增加至订单队列
		local orderList=tostring(ARGV[2])
		flag = redis.call("LPUSH", orderList, userId)

-- 增加至用户集合
       flag = redis.call("SADD", orderSet, userId)

-- 库存数减1
		flag = redis.call("DECR", goodsTotal)
-- 返回当时缓存的数量
		return total
	`)
	return script
}

func evalScript(client *redis.Client, userId string, wg *sync.WaitGroup) {
	defer wg.Done()
	script := createScript()
	fmt.Printf("the script is %+v", script)
	return
	sha, err := script.Load(client.Context(), client).Result()
	if err != nil {
		log.Fatalln(err)
	}
	ret := client.EvalSha(client.Context(), sha, []string{
		userId,
		orderSet,
	}, []string{
		goodsTotal,
		orderList,
	})
	if result, err := ret.Result(); err != nil {
		log.Fatalf("Execute Redis fail: %v", err.Error())
	} else {
		fmt.Println("")
		fmt.Printf("userid: %s, result: %d", userId, result)
	}
}

func main() {
	http.HandleFunc("/", addOrder)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	client := getRedis()

	defer r.Body.Close()
	defer client.Close()

	r.ParseForm()
	uid := r.FormValue("uid")

	go evalScript(client, uid, &wg)
	wg.Wait()
}

func getRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	return client
}
