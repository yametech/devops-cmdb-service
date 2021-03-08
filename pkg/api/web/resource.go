package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"io/ioutil"
)

type ResourceApi struct {
	Server
	resourceService *service.ResourceService
}

func (r ResourceApi) GetModelAttributeList(ctx *gin.Context) {
	b, _ := ioutil.ReadAll(ctx.Request.Body)
	fmt.Println(string(b))
	fmt.Println(ctx.Param("modelUid"))
	result := &[]common.ModelAttributeVisibleVO{}
	utils.SimpleConvert(result, r.resourceService.GetModelAttributeList(ctx.Query("modelUid")))
	Success(ctx, result)
}

func (r ResourceApi) ConfigModelAttribute(ctx *gin.Context) {
	result := &[]common.ModelAttributeVisibleVO{}
	json.Unmarshal([]byte(ctx.Query("columns")), result)
	r.resourceService.SetModelAttribute(ctx.Query("modelUid"), result)
	Success(ctx, nil)
}
