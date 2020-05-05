package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	protos "one/proto"
	"time"
)

func main() {
	CreateOrder()
}

func CreateOrder() error {
	var req protos.OrderReq
	var item protos.OrderReqOrderItem
	var items []*protos.OrderReqOrderItem
	req.Address = 1
	req.Extra = "hoooooowooooo"
	item.Id = 1
	item.Amount = 2
	item.ProductSkuId = 3
	item.ProductId = 2
	point := &item
	items = append(items, point)
	req.Item = items
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//监听网络端口
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		logrus.Fatal("failed to connect : ", err)
	}
	logrus.Info("请求rpc接口")
	// 存根
	rpcClient := protos.NewIOrderServiceClient(conn)

	logrus.Info(rpcClient)
	msg, err := rpcClient.CreateOrder(ctx, &req) // Init 方法传入为 nil
	logrus.Info(msg)
	logrus.Info(err)
	return err
}
