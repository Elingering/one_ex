package rpc

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"one/infra"
	"one/infra/base"
	protos "one/proto"
)

var rpcClient protos.IOrderServiceClient

//var once sync.Once

func RpcClient() protos.IOrderServiceClient {
	return rpcClient
}

type GoRpcApiStarter struct {
	infra.BaseStarter
}

func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(OrderRpc))
}

func (g *GoRpcApiStarter) Setup(ctx infra.StarterContext) {
	port := ctx.Props().GetDefault("app.rpc.port", "8082")
	//监听网络端口
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		logrus.Fatal("failed to connect : ", err)
	}
	logrus.Info("请求rpc端口，" + port)
	// 存根
	rpcClient = protos.NewIOrderServiceClient(conn)
}
