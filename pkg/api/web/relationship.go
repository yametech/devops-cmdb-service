package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/service"
)

type RelationshipApi struct {
	resourceService *service.RelationshipService
}

func (r *RelationshipApi) AddModelRelation(ctx *gin.Context) {

}

func (r *RelationshipApi) GetModelRelationList(ctx *gin.Context) {

}
