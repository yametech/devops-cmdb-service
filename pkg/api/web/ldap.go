package web

import (
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
)

type LdapApi struct {
	service.LdapService
}

func (r *LdapApi) router(e *gin.Engine) {
	groupRoute := e.Group(common.WebApiGroup)
	groupRoute.GET("/ldap/user-list", r.getLdapUserList)
}

func (r *LdapApi) getLdapUserList(ctx *gin.Context) {
	pageList, err := r.LdapService.GetLdapUserList()
	if err != nil {
		Error(ctx, err.Error())
		return
	}

	Success(ctx, pageList)
}
