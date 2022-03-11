package main

import (
	person "acurd.com/m/proto/gen/go"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	var user person.User
	user.Id = 42
	fmt.Println(&user)
	b, err := proto.Marshal(&user) //这个就是我们传输的内容,一个二进制流
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%O\n", b)

	//前面的值假设是我们服务端接到了,开始进行解码
	var user2 person.User
	err = proto.Unmarshal(b, &user2)
	fmt.Println(&user2)
	b, err = json.Marshal(&user2) //转化为json
	fmt.Printf("%s\n", b)
}
