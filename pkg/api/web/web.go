package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
)

type Server struct {
	api.IApiServer
}

func NewServer(apiServer api.IApiServer) *Server {
	server := &Server{
		apiServer,
	}
	groupRoute := server.GINEngine().Group(common.WEB_API_GROUP)
	groupRoute.GET("/member", server.ListMemberApi)

	return server
}

func (s *Server) ListMemberApi(ctx *gin.Context) {
	ctx.JSON(200, map[string]interface{}{"abc": "123"})
}
