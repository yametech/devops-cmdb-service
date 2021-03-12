package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/service"
)

type RelationshipApi struct {
	relationshipService *service.RelationshipService
}

func (r *RelationshipApi) addModelRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.relationshipService.AddModelRelation(string(rawData), "")

	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, result)
	}
}

func (r *RelationshipApi) deleteModelRelation(ctx *gin.Context) {
	result, err := r.relationshipService.DeleteModelRelation(ctx.Param("uid"))

	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, result)
	}
}

func (r *RelationshipApi) getModelRelationList(ctx *gin.Context) {
	result := r.relationshipService.GetModelRelationList(ctx.Param("uid"))
	Success(ctx, result)
}

func (r *RelationshipApi) updateModelRelation(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.relationshipService.UpdateModelRelation(string(rawData), "")
	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, result)
	}
}
