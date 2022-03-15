package main

import (
	"acurd.com/m/proto/gen/go/goods"
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	Port = ":9988"
)

type Server struct {
}

var server Server

func (receiver *Server) mustEmbedUnimplementedGoodsInfoServer() {
	panic("implement me")
}

//添加商品
func (receiver *Server) AddGoods(ctx context.Context, req *goods.Goods) (resp *goods.GoodsId) {
	goodsId := goods.GoodsId{}
	goodsId.Value = GetUUID()
	return &goodsId
}

//获取商品
func (receiver *Server) GetGoods(ctx context.Context, req *goods.Goods) (resp *goods.GoodsId) {
	goodsId := goods.GoodsId{}
	goodsId.Value = GetUUID()
	return &goodsId
}

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func main() {
	listener, err := net.Listen("tcp", Port)
	if err != nil {
		log.Println("net listen err ", err)
		return
	}
	s := grpc.NewServer()
	goods.RegisterGoodsInfoServer(s, server)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}

}
