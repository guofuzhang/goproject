package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", addOrder)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	//从redis里面读取数据,如果读取到了,就进入下单环节
	var inckey = "inc-count"
	var orderList = "inc-orderlist"
	var total int64 = 1000
	client := getRedis()
	defer client.Close()
	defer r.Body.Close()
	var res = client.IncrBy(inckey, 1)
	val := res.Val()
	if res.Err() != nil {
		fmt.Print(res.Err())
		return
	}
	fmt.Println("我的值现在是", val)
	//return
	if val <= total {
		//抢到后把用户的id 存入 另外一个队列,用于创建订单
		r.ParseForm()
		uid := r.FormValue("uid")
		client.LPush(orderList, uid)
		msg := fmt.Sprintf("我抢到了,我是第%d抢到的 我的用户id是 %v \n", val, uid)
		_, _ = w.Write([]byte(msg))
		fmt.Print(msg)
	} else {
		msg := "我啥子都没得抢到\n"
		_, _ = w.Write([]byte(msg))
	}
	return
}

func getRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	return client
}
