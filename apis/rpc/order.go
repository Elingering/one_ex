package rpc

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	protos "one/proto"
)

type OrderRpc struct {
}

func (o *OrderRpc) CreateOrder(c context.Context) error {
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
	logrus.Info(c)
	logrus.Info(&req)
	msg, err := RpcClient().CreateOrder(context.Background(), &req)
	logrus.Info(msg)
	logrus.Info(err)
	return err
}
