package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"io/ioutil"
)

type RelationshipApi struct {
	relationshipService *service.RelationshipService
}

func (r *RelationshipApi) AddModelRelation(ctx *gin.Context) {
	// relationshipUid, constraint, sourceUid, targetUid, comment, operator string
	b, _ := ioutil.ReadAll(ctx.Request.Body)
	result, err := r.relationshipService.AddModelRelation(string(b), "")

	if err != nil {
		Error(ctx, err.Error())
	} else {
		Success(ctx, result)
	}
}

func (r *RelationshipApi) GetModelRelationList(ctx *gin.Context) {
	result := r.relationshipService.GetModelRelationList(ctx.Query("uid"))
	Success(ctx, result)
}
