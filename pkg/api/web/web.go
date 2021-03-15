package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"net/http"
)

type Server struct {
	api.IApiServer
	*service.ModelService
	*service.AttributeService
}

// 自定义中间件
func filterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 打印日志
		//start := time.Now()
		rawData, _ := c.GetRawData()
		fmt.Printf("%v,  \"%v\",   %v  \n", c.Request.Method, c.Request.RequestURI, string(rawData))

		// 用户信息存储
		//c.Set("user", "中间件")
		c.Set("body", string(rawData))
		c.Next()
		//end := time.Now()
		//latency := end.Sub(start)
		//fmt.Printf("%v,  \"%v\", 耗时：%v \n", c.Request.Method, c.Request.RequestURI, latency.String())
	}
}

func NewServer(apiServer api.IApiServer) *Server {
	//apiServer.GINEngine()
	ms, as := api.NewService()
	server := &Server{
		apiServer,
		ms,
		as,
	}
	//apiServer.GINEngine().Use(filterMiddleware())

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

	groupRoute.POST("/model/relationship-list", server.getAllRelationship)
	groupRoute.POST("/model/relationship-add", server.createRelationship)
	groupRoute.POST("/model/relationship-update", server.updateRelationship)
	groupRoute.POST("/model/relationship-delete", server.deleteRelationship)

	// resource
	resource := &ResourceApi{&service.ResourceService{}}
	resource.router(apiServer.GINEngine())

	// relationship
	relationship := &RelationshipApi{&service.RelationshipService{}}
	relationship.router(apiServer.GINEngine())

	return server
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &common.ApiResponseVO{Data: data, Code: 200})
}

func Error(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, &common.ApiResponseVO{Msg: msg, Code: 400})
}

func ResultHandle(ctx *gin.Context, result interface{}, err error) {
	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, result)
	}
}
