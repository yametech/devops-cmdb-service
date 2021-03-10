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
	*service.ModelService
	*service.AttributeService
}

func NewServer(apiServer api.IApiServer) *Server {
	apiServer.GINEngine()
	server := &Server{
		apiServer,
		&service.ModelService{},
		&service.AttributeService{},
	}
	groupRoute := apiServer.GINEngine().Group(common.WEB_API_GROUP)
	groupRoute.GET("/member", server.ListMemberApi)
	//groupRoute.POST("/model-group", server.ListModelGroup)

	groupRoute.GET("/model-group", server.getAllGroup)
	groupRoute.GET("/model-group/:uid", server.getGroup)
	groupRoute.POST("/model-group", server.createGroup)
	groupRoute.PUT("/model-group/:uid", server.putGroup)
	groupRoute.DELETE("/model-group/:uid", server.deleteGroup)

	groupRoute.GET("/model", server.getAllModel)
	groupRoute.GET("/model/:uid", server.getModel)
	groupRoute.POST("/model", server.createModel)
	groupRoute.PUT("/model/:uid", server.putModel)
	groupRoute.DELETE("/model/:uid", server.deleteModel)

	groupRoute.GET("/attribute-group", server.getAllAttributeGroup)
	groupRoute.GET("/attribute-group/:uid", server.getAttributeGroup)
	groupRoute.POST("/attribute-group", server.createAttributeGroup)
	groupRoute.PUT("/attribute-group/:uid", server.putAttributeGroup)
	groupRoute.DELETE("/attribute-group/:uid", server.deleteAttributeGroup)

	groupRoute.GET("/attribute", server.getAllAttribute)
	groupRoute.GET("/attribute/:uid", server.getAttribute)
	groupRoute.POST("/attribute", server.createAttribute)
	groupRoute.PUT("/attribute/:uid", server.putAttribute)
	groupRoute.DELETE("/attribute/:uid", server.deleteAttribute)

	// resource
	resource := &ResourceApi{&service.ResourceService{}}
	groupRoute.POST("/resource/model-attribute-list", resource.GetModelAttributeList)
	groupRoute.POST("/resource/model-list", resource.GetModelList)
	groupRoute.POST("/resource/resource-list/:modelUid", resource.GetResourceList)
	groupRoute.POST("/resource/resource-detail/:uuid", resource.GetResourceDetail)

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
