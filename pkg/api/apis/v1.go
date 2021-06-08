package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
)

type Server struct {
	api.IApiServer
	*service.V1
}

func NewServer(apiServer api.IApiServer) *Server {
	common.ApiConfig()

	server := &Server{
		apiServer,
		&service.V1{},
	}

	groupRoute := apiServer.GINEngine().Group(common.NormalApiGroup + "/v1")
	groupRoute.GET("/app-tree", server.GetAppTree)

	return server
}

func (s *Server) GetAppTree(ctx *gin.Context) {
	result, _ := s.V1.GetAppTree()
	api.RequestOK(ctx, result)
}
