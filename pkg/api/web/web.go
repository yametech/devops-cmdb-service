package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
)

type Server struct {
	api.IApiServer
}

func NewServer(apiServer api.IApiServer) *Server {
	apiServer.GINEngine()
	server := &Server{
		apiServer,
	}
	groupRoute := apiServer.GINEngine().Group(common.WEB_API_GROUP)
	groupRoute.GET("/member", server.ListMemberApi)
	groupRoute.POST("/model-group", server.ListModelGroup)

	groupRoute.GET("/model_group", server.getAllGroup)
	groupRoute.GET("/model_group/:uid", server.getGroup)
	groupRoute.POST("/model_group", server.createGroup)
	groupRoute.PUT("/model_group/:uid", server.putGroup)
	groupRoute.DELETE("/model_group/:uid", server.deleteGroup)

	groupRoute.GET("/model", server.getAllModel)
	groupRoute.GET("/model/:uid", server.getModel)
	groupRoute.POST("/model", server.createGroup)
	groupRoute.PUT("/model/:uid", server.putGroup)
	groupRoute.DELETE("/model/:uid", server.deleteGroup)

	// resource
	resource := &ResourceApi{*server, &service.ResourceService{}}
	groupRoute.POST("resource/model-attribute-list", resource.GetModelAttributeList)
	groupRoute.POST("resource/model-list", resource.GetModelList)

	return server
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, &common.ApiResponseVO{Data: data, Code: 200})
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(200, &common.ApiResponseVO{Msg: msg, Code: 400})
}

func (s *Server) ListMemberApi(ctx *gin.Context) {
	ctx.JSON(200, map[string]interface{}{"abc": "123"})
}

func (s *Server) ListModelGroup(ctx *gin.Context) {

	fmt.Printf("uid=%v\n", ctx.Param("uid"))
	fmt.Printf("uid=%v\n", ctx.Query("uid"))
	mgs := &service.ModeGroupService{}
	ctx.JSON(200, mgs.ListByUid(ctx.Query("uid")))
}

func (s *Server) ListModel(ctx *gin.Context) {
	ms := &service.ModeService{}
	ctx.JSON(200, ms.List())
}
