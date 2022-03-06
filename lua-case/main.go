package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const orderSet = "orderSet"     //用户id的集合
const goodsTotal = "goodsTotal" //商品库存的key
const orderList = "orderList"   //订单队列
func createScript() *redis.Script {
	str, err := ioutil.ReadFile("./lua-case/script.lua")
	if err != nil {
		fmt.Println("Script read error", err)
		log.Println(err)
	}
	scriptStr := fmt.Sprintf("%s", str)
	script := redis.NewScript(scriptStr)
	return script
}

func evalScript(client *redis.Client, userId string, wg *sync.WaitGroup) {
	defer wg.Done()
	script := createScript()
	//fmt.Printf("%+v",script)
	//return
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
