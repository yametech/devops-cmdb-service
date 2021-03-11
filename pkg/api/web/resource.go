package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"strconv"
)

type ResourceApi struct {
	//Server
	resourceService *service.ResourceService
}

// 获取资源实例字段列表
func (r *ResourceApi) getModelAttribute(ctx *gin.Context) {
	result := &[]common.ModelAttributeVisibleVO{}
	utils.SimpleConvert(result, r.resourceService.GetModelAttributeList(ctx.Param("uid")))
	Success(ctx, result)
}

// 模型字段预览显示设置
func (r *ResourceApi) configModelAttribute(ctx *gin.Context) {
	b, _ := ctx.GetRawData()
	req := &common.ConfigModelAttributeVO{}
	json.Unmarshal(b, req)
	r.resourceService.SetModelAttribute(req.Uid, req.Columns)
	Success(ctx, nil)
}

// 获取模型菜单
func (r *ResourceApi) getModelMenu(ctx *gin.Context) {
	result := &[]common.ModelMenuVO{}
	utils.SimpleConvert(result, r.resourceService.GetAllModeList())
	Success(ctx, result)
}

func (r *ResourceApi) addResource(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.resourceService.AddResource(string(rawData), "")
	if err != nil {
		panic(err)
	}
	Success(ctx, result)
}

// 获取模型实例列表
func (r *ResourceApi) getResourcePageList(ctx *gin.Context) {
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	pageNumberStr := ctx.DefaultQuery("page_number", "1")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		Error(ctx, "参数有误")
		return
	}
	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil || pageNumber <= 0 {
		Error(ctx, "参数有误")
		return
	}

	Success(ctx, r.resourceService.GetResourcePageList(ctx.Query("modelUid"), pageNumber, pageSize))
}

func (r *ResourceApi) getResourceDetail(ctx *gin.Context) {
	result, err := r.resourceService.GetResourceDetail(ctx.Param("uuid"))
	if err != nil {
		result = nil
		fmt.Printf("找不到记录,uuid:%v, msg:%v\n", ctx.Param("uuid"), err)
	}
	Success(ctx, result)
}

func (r *ResourceApi) deleteResource(ctx *gin.Context) {
	err := r.resourceService.DeleteResource(ctx.Param("uuid"))
	if err != nil {
		fmt.Printf("找不到记录,uuid:%v, msg:%v\n", ctx.Param("uuid"), err)
		Error(ctx, err.Error())
	} else {
		Success(ctx, "删除成功")
	}
}

func (r *ResourceApi) updateResourceAttribute(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	dataMap := map[string]string{}
	_ = json.Unmarshal(rawData, &dataMap)
	err := r.resourceService.UpdateResourceAttribute(ctx.Param("uuid"), dataMap["attributeInsValue"], "")
	if err != nil {
		fmt.Printf("UpdateResourceAttribute更新失败, msg:%v\n", err)
		Error(ctx, err.Error())
	} else {
		Success(ctx, "更新成功")
	}
}
