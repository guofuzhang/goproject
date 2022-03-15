package main

import (
	"acurd.com/m/proto/gen/go/goods"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
)

const (
	Adress = "localhost:9988"
)

type Server struct {
	goods.UnimplementedGoodsInfoServer
}

func (s Server) AddGoods(ctx context.Context, g *goods.Goods) (*goods.GoodsId, error) {
	//参数校验忽略
	db, err := sql.Open("mysql", "acurd:acurd@tcp(127.0.0.1:3306)/test")
	if err != nil {
		grpclog.Fatal(err)
	}
	defer db.Close()

	getUUID := GetUUID()
	name := g.Name
	desc := g.Desc
	sqlStr := "INSERT INTO an_goods (`id`,`name`,`desc`) VALUES (?,?,?)"
	_, err = db.Exec(sqlStr, getUUID, name, desc)
	if err != nil {
		grpclog.Fatal(err)
	}

	goodsId := goods.GoodsId{}
	goodsId.Value = getUUID

	log.Printf("add success ,the id is %s", getUUID)
	return &goodsId, nil
}

func (s Server) GetGoods(ctx context.Context, id *goods.GoodsId) (*goods.Goods, error) {

	db, err := sql.Open("mysql", "acurd:acurd@tcp(127.0.0.1:3306)/test")
	if err != nil {
		grpclog.Fatal(err)
	}
	defer db.Close()

	search_id := id.Value
	g := goods.Goods{}

	sqlStr := "SELECT `id`,`name`,`desc` FROM an_goods WHERE id=?"
	err = db.QueryRow(sqlStr, search_id).Scan(&g.Id, &g.Name, &g.Desc)

	return &g, nil
}

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

var server Server

func main() {
	//绑定监听端口
	listener, err := net.Listen("tcp", Adress)
	if err != nil {
		log.Println("net listen err ", err)
		return
	}
	//创建grpc服务
	s := grpc.NewServer()
	//开始grpc监听
	goods.RegisterGoodsInfoServer(s, &server)
	log.Println("start gRPC listen on Adress " + Adress)

	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}

}
