package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
)

type RelationshipApi struct {
	relationshipService *service.RelationshipService
}

func (r *RelationshipApi) router(e *gin.Engine) {
	groupRoute := e.Group(common.WEB_API_GROUP)

	groupRoute.POST("/model-relation-list", r.getModelRelationList)
	groupRoute.POST("/model-relation-usage", r.getModelRelationUsageCount)
	groupRoute.POST("/add-model-relation", r.addModelRelation)
	groupRoute.POST("/delete-model-relation", r.deleteModelRelation)
	groupRoute.POST("/update-model-relation", r.updateModelRelation)

	groupRoute.GET("/resource-relation/:uuid", r.getResourceRelationList)
	groupRoute.POST("/resource-relation", r.addResourceRelation)
	groupRoute.DELETE("/resource-relation", r.deleteResourceRelation)
}

func (r *RelationshipApi) addModelRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.relationshipService.AddModelRelation(string(rawData), "")
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) deleteModelRelation(ctx *gin.Context) {
	data, _ := ctx.GetRawData()
	param, _ := utils.Stream2Json(data)
	result, err := r.relationshipService.DeleteModelRelation((*param)["uid"])
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) getModelRelationList(ctx *gin.Context) {
	data, _ := ctx.GetRawData()
	param, _ := utils.Stream2Json(data)
	result := r.relationshipService.GetModelRelationList((*param)["uid"])
	Success(ctx, result)
}

func (r *RelationshipApi) updateModelRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.relationshipService.UpdateModelRelation(string(rawData), "")
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) getModelRelationUsageCount(ctx *gin.Context) {
	result, err := r.relationshipService.GetResourceRelationsByModelRelationUid(ctx.Query("uid"))
	ResultHandle(ctx, len(result), err)
}

func (r *RelationshipApi) getResourceRelationList(ctx *gin.Context) {
	result, err := r.relationshipService.GetResourceRelationList(ctx.Param("uuid"))
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) addResourceRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	paramMap, err := utils.Stream2Json(rawData)
	if err != nil {
		ResultHandle(ctx, paramMap, err)
		return
	}

	result, err := r.relationshipService.AddResourceRelation((*paramMap)["source_uuid"], (*paramMap)["target_uuid"], (*paramMap)["uid"])
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) deleteResourceRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	paramMap, err := utils.Stream2Json(rawData)
	if err != nil {
		ResultHandle(ctx, paramMap, err)
		return
	}
	result, err := r.relationshipService.DeleteResourceRelation((*paramMap)["source_uuid"], (*paramMap)["target_uuid"], (*paramMap)["uid"])
	ResultHandle(ctx, result, err)
}
