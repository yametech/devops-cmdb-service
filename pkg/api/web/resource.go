package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"io/ioutil"
	"strconv"
)

type ResourceApi struct {
	//Server
	resourceService *service.ResourceService
}

// 获取资源实例字段列表
func (r *ResourceApi) GetModelAttributeList(ctx *gin.Context) {
	result := &[]common.ModelAttributeVisibleVO{}
	utils.SimpleConvert(result, r.resourceService.GetModelAttributeList(ctx.Param("modelUid")))
	Success(ctx, result)
}

// 模型字段预览显示设置
func (r *ResourceApi) ConfigModelAttribute(ctx *gin.Context) {
	b, _ := ioutil.ReadAll(ctx.Request.Body)
	req := &common.ConfigModelAttributeVO{}
	json.Unmarshal(b, req)
	r.resourceService.SetModelAttribute(req.ModelUid, req.Columns)
	Success(ctx, nil)
}

// 获取模型列表
func (r *ResourceApi) GetModelList(ctx *gin.Context) {
	result := &[]common.SimpleModelVO{}
	utils.SimpleConvert(result, r.resourceService.GetModeList())
	Success(ctx, result)
}

func (r *ResourceApi) AddResource(ctx *gin.Context) {
	b, _ := ioutil.ReadAll(ctx.Request.Body)
	result, err := r.resourceService.AddResource(string(b), "")
	if err != nil {
		panic(err)
	}
	Success(ctx, result)
}

// 获取模型实例列表
func (r *ResourceApi) GetResourcePageList(ctx *gin.Context) {
	currentPage, _ := strconv.Atoi(ctx.Query("currentPage"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	Success(ctx, r.resourceService.GetResourcePageList(ctx.Query("modelUid"), currentPage, pageSize))
}

func (r *ResourceApi) GetResourceDetail(ctx *gin.Context) {
	result, err := r.resourceService.GetResourceDetail(ctx.Param("uuid"))
	if err != nil {
		result = nil
		fmt.Printf("找不到记录,uuid:%v, msg:%v\n", ctx.Param("uuid"), err)
	}
	Success(ctx, result)
}

func (r *ResourceApi) DeleteResource(ctx *gin.Context) {
	err := r.resourceService.DeleteResource(ctx.Param("uuid"))
	if err != nil {
		fmt.Printf("找不到记录,uuid:%v, msg:%v\n", ctx.Param("uuid"), err)
		Error(ctx, err.Error())
	} else {
		Success(ctx, "删除成功")
	}
}

func (r *ResourceApi) UpdateResourceAttribute(ctx *gin.Context) {
	err := r.resourceService.UpdateResourceAttribute(ctx.Query("uuid"), ctx.Query("attributeInsValue"), "")
	if err != nil {
		//b, _ := ioutil.ReadAll(ctx.Request.Body)
		fmt.Printf("UpdateResourceAttribute更新失败,RequestURI:%v, msg:%v\n", ctx.Request.RequestURI, err)
		Error(ctx, err.Error())
	} else {
		Success(ctx, "更新成功")
	}
}
