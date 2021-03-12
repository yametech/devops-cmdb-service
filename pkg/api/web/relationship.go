package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
)

type RelationshipApi struct {
	relationshipService *service.RelationshipService
}

func (r *RelationshipApi) addModelRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.relationshipService.AddModelRelation(string(rawData), "")
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) deleteModelRelation(ctx *gin.Context) {
	result, err := r.relationshipService.DeleteModelRelation(ctx.Param("uid"))
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) getModelRelationList(ctx *gin.Context) {
	result := r.relationshipService.GetModelRelationList(ctx.Param("uid"))
	Success(ctx, result)
}

func (r *RelationshipApi) updateModelRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.relationshipService.UpdateModelRelation(string(rawData), "")
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) getModelRelationUsageCount(ctx *gin.Context) {
	result, err := r.relationshipService.GetResourceRelationsByModelRelationUid(ctx.Param("uid"))
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
