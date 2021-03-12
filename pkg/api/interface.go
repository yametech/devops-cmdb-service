package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type IApiServer interface {
	Run() error
	Stop() error
	GINEngine() *gin.Engine
}

var _ IApiServer = &BaseServer{}

type BaseServer struct {
	addrs []string
	e     *gin.Engine
}

func (b *BaseServer) GINEngine() *gin.Engine {
	return b.e
}

func NewBaseServer(addr string) IApiServer {
	baseServer := &BaseServer{e: gin.Default(), addrs: []string{addr}}
	return baseServer
}

func NewService() (*service.ModelService, *service.AttributeService) {
	session := store.GetSession(false)
	ms := &service.ModelService{Session: session}
	as := &service.AttributeService{Session: session}
	return ms, as
}

func (b *BaseServer) Run() error {
	return b.e.Run(b.addrs...)
}

func (b *BaseServer) Stop() error {
	panic("implement me")
}
