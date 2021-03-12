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
	ms, as := api.NewService()
	server := &Server{
		apiServer,
		ms,
		as,
	}
	groupRoute := apiServer.GINEngine().Group(common.WEB_API_GROUP)

	groupRoute.POST("/model/model-group-list", server.getAllGroup)
	groupRoute.POST("/model/model-group-detail", server.getGroup)
	groupRoute.POST("/model/model-group-add", server.createGroup)
	groupRoute.POST("/model/model-group-update", server.putGroup)
	groupRoute.POST("/model/model-group-delete", server.deleteGroup)

	groupRoute.POST("/model/model-list", server.getAllModel)
	groupRoute.POST("/model/model-detail", server.getModel)
	groupRoute.POST("/model/model-add", server.createModel)
	groupRoute.POST("/model/model-update", server.putModel)
	groupRoute.POST("/model/model-delete", server.deleteModel)

	groupRoute.POST("/model/attribute-group-list", server.getAllAttributeGroup)
	groupRoute.POST("/model/attribute-group-detail", server.getAttributeGroup)
	groupRoute.POST("/model/attribute-group-add", server.createAttributeGroup)
	groupRoute.POST("/model/attribute-group-update", server.putAttributeGroup)
	groupRoute.POST("/model/attribute-group-delete", server.deleteAttributeGroup)

	groupRoute.POST("/model/attribute-list", server.getAllAttribute)
	groupRoute.POST("/model/attribute-detail", server.getAttribute)
	groupRoute.POST("/model/attribute-add", server.createAttribute)
	groupRoute.POST("/model/attribute-update", server.putAttribute)
	groupRoute.POST("/model/attribute-delete", server.deleteAttribute)

	// resource
	resource := &ResourceApi{&service.ResourceService{}}
	groupRoute.GET("/model-menu", resource.getModelMenu)
	groupRoute.GET("/model-attribute/:uid", resource.getModelAttribute)
	groupRoute.PUT("/model-attribute/:uid", resource.configModelAttribute)
	groupRoute.GET("/model-info/:uid", resource.getModelInfoForIns)
	groupRoute.GET("/resource", resource.getResourceListPage)
	groupRoute.GET("/resource/:uuid", resource.getResourceDetail)
	groupRoute.POST("/resource", resource.addResource)
	groupRoute.DELETE("/resource/:uuids", resource.deleteResource)
	groupRoute.PUT("/resource-attribute/:uuid", resource.updateResourceAttribute)

	// relationship
	relationship := &RelationshipApi{&service.RelationshipService{}}
	groupRoute.GET("/model-relation/:uid", relationship.getModelRelationList)
	groupRoute.GET("/model-relation/:uid/usage", relationship.getModelRelationUsageCount)
	groupRoute.POST("/model-relation", relationship.addModelRelation)
	groupRoute.DELETE("/model-relation/:uid", relationship.deleteModelRelation)
	groupRoute.PUT("/model-relation/:uid", relationship.updateModelRelation)

	groupRoute.GET("/resource-relation/:uuid", relationship.getResourceRelationList)
	groupRoute.POST("/resource-relation", relationship.addResourceRelation)
	groupRoute.DELETE("/resource-relation", relationship.deleteResourceRelation)

	return server
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, &common.ApiResponseVO{Data: data, Code: 200})
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(200, &common.ApiResponseVO{Msg: msg, Code: 400})
}

func ResultHandle(ctx *gin.Context, result interface{}, err error) {
	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, result)
	}
}
