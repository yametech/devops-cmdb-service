package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"regexp"
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
	groupRoute.POST("/resource-list", r.getResourceListPage)
	groupRoute.GET("/resource/:uuid", r.getResourceDetail)
	groupRoute.POST("/resource", r.addResource)
	groupRoute.PUT("/resource", r.updateResource)
	groupRoute.DELETE("/resource", r.deleteResource)
	groupRoute.PUT("/resource-attribute/:uuid", r.updateResourceAttribute)
}

// 获取资源实例字段列表
func (r *ResourceApi) getModelAttribute(ctx *gin.Context) {
	//result := &[]common.ModelAttributeVisibleVO{}
	//utils.SimpleConvert(result, r.resourceService.GetModelAttributeList(ctx.Param("uid")))
	Success(ctx, r.resourceService.GetModelAttributeList(ctx.Param("uid")))
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
	modelService := service.ModelService{}
	list, _ := modelService.GetSimpleModelList()
	utils.SimpleConvert(result, list)
	Success(ctx, result)
}

func (r *ResourceApi) addResource(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.resourceService.AddResource(string(rawData), ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		Error(ctx, err.Error())
		return
	}
	Success(ctx, result)
}

func (r *ResourceApi) updateResource(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	result, err := r.resourceService.UpdateResource(string(rawData), ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		Error(ctx, err.Error())
		return
	}
	Success(ctx, result)
}

// 获取模型实例列表
func (r *ResourceApi) getResourceListPage(ctx *gin.Context) {
	vo := &common.ResourceListPageParamVO{}
	err := ctx.ShouldBindJSON(vo)
	if err != nil {
		println(err.Error())
		Error(ctx, err.Error())
		return
	}

	if vo.QueryValue != "" {
		Success(ctx, r.resourceService.GetResourceListPageByQueryValue(vo.ModelUid, vo.QueryValue, vo.Current, vo.PageSize))
	} else {
		if vo.QueryMap == nil {
			vo.QueryMap = &map[string]string{}
		}
		id, ok := (*vo.QueryMap)["ID"]
		if ok {
			match, err := regexp.MatchString("^\\d*$", strings.TrimSpace(id))
			if err != nil {
				Error(ctx, err.Error())
				return
			}
			if !match {
				Error(ctx, "ID内容必须是整数")
				return
			}
		}
		Success(ctx, r.resourceService.GetResourceListPageByMap(vo.UUID, vo.ModelUid, vo.ModelRelationUid, vo.Current, vo.PageSize, vo.QueryMap))
	}

}

func (r *ResourceApi) getResourceDetail(ctx *gin.Context) {
	result, err := r.resourceService.GetResourceDetail(ctx.Param("uuid"))
	if err != nil {
		Error(ctx, err.Error())
		return
	}
	Success(ctx, result)
}

func (r *ResourceApi) deleteResource(ctx *gin.Context) {
	rawData, _ := ctx.GetRawData()
	dataMap := map[string]string{}
	_ = json.Unmarshal(rawData, &dataMap)
	err := r.resourceService.DeleteResource(strings.Split(dataMap["uuids"], ","))
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
	err := r.resourceService.UpdateResourceAttribute(ctx.Param("uuid"), dataMap["attributeInsValue"], ctx.GetHeader("x-wrapper-username"))
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
