package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"strings"
)

type ResourceApi struct {
	//Server
	resourceService *service.ResourceService
}

func (r *ResourceApi) router(e *gin.Engine) {
	groupRoute := e.Group(common.WEB_API_GROUP)
	groupRoute.GET("/model-menu", r.getModelMenu)
	groupRoute.GET("/model-attribute/:uid", r.getModelAttribute)
	groupRoute.PUT("/model-attribute/:uid", r.configModelAttribute)
	groupRoute.GET("/model-info/:uid", r.getModelInfoForIns)
	groupRoute.GET("/resource", r.getResourceListPage)
	groupRoute.GET("/resource/:uuid", r.getResourceDetail)
	groupRoute.POST("/resource", r.addResource)
	groupRoute.DELETE("/resource/:uuids", r.deleteResource)
	groupRoute.PUT("/resource-attribute/:uuid", r.updateResourceAttribute)
}

// 获取资源实例字段列表
func (r *ResourceApi) getModelAttribute(ctx *gin.Context) {
	result := &[]common.ModelAttributeVisibleVO{}
	utils.SimpleConvert(result, r.resourceService.GetModelAttributeList(ctx.Param("uid")))
	Success(ctx, result)
}

// 模型字段预览显示设置
func (r *ResourceApi) configModelAttribute(ctx *gin.Context) {
	req := &common.ConfigModelAttributeVO{}
	ctx.ShouldBindJSON(req)
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
func (r *ResourceApi) getResourceListPage(ctx *gin.Context) {
	vo := &common.ResourceListPageParamVO{}

	err := ctx.ShouldBindQuery(vo)
	//err := ctx.ShouldBindJSON(vo)
	if err != nil {
		println(err.Error())
		Error(ctx, err.Error())
		return
	}

	if vo.QueryValue != "" {
		Success(ctx, r.resourceService.GetResourceListPage(vo.ModelUid, vo.QueryValue, vo.Current, vo.PageSize))
	} else {
		if vo.QueryMap == nil {
			vo.QueryMap = &map[string]string{}
		}

		Success(ctx, r.resourceService.GetResourceListPageByMap(vo.ModelUid, vo.Current, vo.PageSize, vo.QueryMap))
	}

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

	err := r.resourceService.DeleteResource(strings.Split(ctx.Param("uuids"), ","))
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

func (r *ResourceApi) getModelInfoForIns(ctx *gin.Context) {
	//fmt.Println(ctx.Get("user"))
	result, err := r.resourceService.GetModelInfoForIns(ctx.Param("uid"))
	ResultHandle(ctx, result, err)
}
