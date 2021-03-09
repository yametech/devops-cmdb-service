package web

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
)

type ResourceApi struct {
	Server
	resourceService *service.ResourceService
}

// 获取资源实例字段列表
func (r *ResourceApi) GetModelAttributeList(ctx *gin.Context) {
	//b, _ := ioutil.ReadAll(ctx.Request.Body)
	//fmt.Println(string(b))
	result := &[]common.ModelAttributeVisibleVO{}
	utils.SimpleConvert(result, r.resourceService.GetModelAttributeList(ctx.Param("modelUid")))
	Success(ctx, result)
}

func (r *ResourceApi) ConfigModelAttribute(ctx *gin.Context) {
	result := &[]common.ModelAttributeVisibleVO{}
	json.Unmarshal([]byte(ctx.Query("columns")), result)
	r.resourceService.SetModelAttribute(ctx.Query("modelUid"), result)
	Success(ctx, nil)
}

// 获取模型列表
func (r *ResourceApi) GetModelList(ctx *gin.Context) {
	result := &[]common.SimpleModelVO{}
	utils.SimpleConvert(result, (&service.ModeService{}).List())
	Success(ctx, result)
}

// 获取模型实例列表
func (r *ResourceApi) GetModelInsList(ctx *gin.Context) {
	result := &[]common.SimpleModelVO{}
	utils.SimpleConvert(result, (&service.ModeService{}).List())
	Success(ctx, result)
}
