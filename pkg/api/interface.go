package api

import (
	"github.com/gin-gonic/gin"
)

type IApiServer interface {
	Run() error
	Stop() error
	GINEngine() *gin.Engine
}

//var _ IApiServer = &BaseServer{}

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

func (b *BaseServer) Run() error {
	return b.e.Run(b.addrs...)
}

func (b *BaseServer) Stop() error {
	panic("implement me")
}
