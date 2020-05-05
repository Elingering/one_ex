package base

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"one/infra"
	"one/models"
	"one/proto"
	"reflect"
	"sync"
)

var rpcServer *grpc.Server
var once sync.Once

func RpcServer() *grpc.Server {
	return rpcServer
}
func RpcRegister(ri interface{}) {
	typ := reflect.TypeOf(ri)
	log.Infof("RPC Register: %s", typ.String())
	//RpcServer().Register(ri)
	protos.RegisterIOrderServiceServer(RpcServer(), models.NewOrderService())
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *grpc.Server
}

func (s *GoRPCStarter) Init(ctx infra.StarterContext) {
	once.Do(func() {
		s.server = grpc.NewServer()
	})
	rpcServer = s.server
}

func (s *GoRPCStarter) Start(ctx infra.StarterContext) {
	port := ctx.Props().GetDefault("app.rpc.port", "8082")
	//监听网络端口
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic(err)
	}
	log.Info("tcp port listened for rpc:", port)
	//处理网络连接和请求
	go s.server.Serve(listener)
}
