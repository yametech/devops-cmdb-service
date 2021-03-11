package web

import (
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

	groupRoute.GET("/model-group", server.getAllGroup)
	groupRoute.GET("/model-group/:uuid", server.getGroup)
	groupRoute.POST("/model-group", server.createGroup)
	groupRoute.PUT("/model-group/:uuid", server.putGroup)
	groupRoute.DELETE("/model-group/:uuid", server.deleteGroup)

	groupRoute.GET("/model", server.getAllModel)
	groupRoute.GET("/model/:uuid", server.getModel)
	groupRoute.POST("/model", server.createModel)
	groupRoute.PUT("/model/:uuid", server.putModel)
	groupRoute.DELETE("/model/:uuid", server.deleteModel)

	groupRoute.GET("/attribute-group", server.getAllAttributeGroup)
	groupRoute.GET("/attribute-group/:uuid", server.getAttributeGroup)
	groupRoute.POST("/attribute-group", server.createAttributeGroup)
	groupRoute.PUT("/attribute-group/:uuid", server.putAttributeGroup)
	groupRoute.DELETE("/attribute-group/:uuid", server.deleteAttributeGroup)

	groupRoute.GET("/attribute", server.getAllAttribute)
	groupRoute.GET("/attribute/:uuid", server.getAttribute)
	groupRoute.POST("/attribute", server.createAttribute)
	groupRoute.PUT("/attribute/:uuid", server.putAttribute)
	groupRoute.DELETE("/attribute/:uuid", server.deleteAttribute)

	// resource
	resource := &ResourceApi{&service.ResourceService{}}
	groupRoute.GET("/model-menu", resource.getModelMenu)
	groupRoute.GET("/model-attribute/:uid", resource.getModelAttribute)
	groupRoute.PUT("/model-attribute/:uid", resource.configModelAttribute)
	groupRoute.GET("/resource", resource.getResourcePageList)
	groupRoute.GET("/resource/:uuid", resource.getResourceDetail)
	groupRoute.POST("/resource", resource.addResource)
	groupRoute.DELETE("/resource/:uuid", resource.deleteResource)
	groupRoute.PUT("/resource-attribute/:uuid", resource.updateResourceAttribute)

	// relationship
	relationship := &RelationshipApi{&service.RelationshipService{}}
	groupRoute.GET("/model-relation/:uid", relationship.getModelRelationList)
	groupRoute.POST("/model-relation", relationship.addModelRelation)
	groupRoute.DELETE("/model-relation/:uid", relationship.deleteModelRelation)
	groupRoute.PUT("/model-relation/:uid", relationship.updateModelRelation)

	return server
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, &common.ApiResponseVO{Data: data, Code: 200})
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(200, &common.ApiResponseVO{Msg: msg, Code: 400})
}
