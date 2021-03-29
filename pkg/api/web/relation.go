package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
)

type RelationshipApi struct {
	relationService *service.RelationService
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
	vo := &common.AddModelRelationVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	result, err := r.relationService.AddModelRelation(vo, ctx.GetHeader("x-wrapper-username"))
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) deleteModelRelation(ctx *gin.Context) {
	data, _ := ctx.GetRawData()
	param, _ := utils.Stream2Json(data)
	result, err := r.relationService.DeleteModelRelation((*param)["uid"])
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) getModelRelationList(ctx *gin.Context) {
	data, _ := ctx.GetRawData()
	param, _ := utils.Stream2Json(data)
	result := r.relationService.GetModelRelationList((*param)["uid"])
	Success(ctx, result)
}

func (r *RelationshipApi) updateModelRelation(ctx *gin.Context) {
	vo := &common.UpdateModelRelationVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	result, err := r.relationService.UpdateModelRelation(vo, ctx.GetHeader("x-wrapper-username"))
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) getModelRelationUsageCount(ctx *gin.Context) {
	result, err := r.relationService.GetResourceRelationsByModelRelationUid(ctx.Query("uid"))
	ResultHandle(ctx, len(result), err)
}

func (r *RelationshipApi) getResourceRelationList(ctx *gin.Context) {
	result, err := r.relationService.GetResourceRelationList(ctx.Param("uuid"))
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) addResourceRelation(ctx *gin.Context) {
	vo := &common.ResourceRelationVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		ResultHandle(ctx, nil, err)
		return
	}

	result, err := r.relationService.AddResourceRelation(vo.SourceUUID, vo.TargetUUID, vo.Uid)
	ResultHandle(ctx, result, err)
}

func (r *RelationshipApi) deleteResourceRelation(ctx *gin.Context) {
	vo := &common.ResourceRelationVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		ResultHandle(ctx, nil, err)
		return
	}
	result, err := r.relationService.DeleteResourceRelation(vo.SourceUUID, vo.TargetUUID, vo.Uid)
	ResultHandle(ctx, result, err)
}
