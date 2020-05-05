package models

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	protos "one/proto"
)

//对外提供的工厂函数
func NewOrderService() *OrderService {
	return &OrderService{}
}

//**************************************************************
// 接口实现，接口定义是在proto生成的.pb.go文件中
//**************************************************************

// 接口实现对象，属性成员根据而业务自定义
type OrderService struct {
}

// Get接口方法实现
func (o *OrderService) CreateOrder(ctx context.Context, req *protos.OrderReq) (*protos.OkMsg, error) {
	logrus.Infof("%+v", req)
	return &protos.OkMsg{Status: "ok"}, nil
}

func (o *OrderService) UpdateOrder(ctx context.Context, req *protos.OrderUpdate) (*protos.OkMsg, error) {
	logrus.Infof("%+v", req)
	return &protos.OkMsg{Status: "ok"}, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, req *protos.Order) (*protos.OkMsg, error) {
	logrus.Infof("%+v", req)
	return &protos.OkMsg{Status: "ok"}, nil
}

func (o *OrderService) OrderDetail(ctx context.Context, req *protos.Order) (*protos.OrderInfo, error) {
	logrus.Infof("%+v", req)
	return &protos.OrderInfo{Status: "ok"}, nil
}

//// GetList接口方法实现
//func (o *OrderService) GetList(req *protos.UserReq, stream protos.IOrderService_GetListServer) error {
//	fmt.Println(*req)
//	// 流式返回多条数据
//	for i := 0; i < 5; i++ {
//		stream.Send(&protos.User{Id: int32(i), Name: "我是" + strconv.Itoa(i)})
//	}
//	return nil
//}
//
//// WaitGet接口方法实现
//func (o *OrderService) WaitGet(reqStream protos.IOrderService_WaitGetServer) error {
//	for { // 接收流式请求并返回单一对象
//		userReq, err := reqStream.Recv()
//		if err != io.EOF {
//			fmt.Println("流请求~", *userReq)
//		} else {
//			return reqStream.SendAndClose(&protos.User{Id: 100, Name: "shuai"})
//		}
//	}
//}
//
////双向流：请求流和响应流异步
//func (o *OrderService) LoopGet(reqStream protos.IOrderService_LoopGetServer) error {
//	for {
//		userReq, err := reqStream.Recv()
//		if err == io.EOF { //请求结束
//			return nil
//		}
//		if err != nil {
//			return err
//		}
//		if err = reqStream.Send(&protos.User{Id: userReq.Id, Name: "shuai"}); err != nil {
//			return err
//		}
//	}
//}
