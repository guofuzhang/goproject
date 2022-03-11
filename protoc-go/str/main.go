package main

import (
	person "acurd.com/m/proto/gen/go"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	var p1 person.Person
	p1.Name = "小名"
	p1.Age = 18
	fmt.Println(&p1)
	b, err := proto.Marshal(&p1) //这个就是我们传输的内容,一个二进制流
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%X\n", b)

	//前面的值假设是我们服务端接到了,开始进行解码
	var p2 person.Person
	err = proto.Unmarshal(b, &p2)
	fmt.Println(&p2)
	b, err = json.Marshal(&p2) //转化为json
	fmt.Printf("%s\n", b)

}
