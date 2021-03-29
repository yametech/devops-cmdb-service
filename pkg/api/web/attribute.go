package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
)

func (s *Server) getAllAttributeGroup(ctx *gin.Context) {
	limit := ctx.DefaultQuery("page_size", "10")
	pageNumber := ctx.DefaultQuery("page_number", "1")

	allAG, err := s.AttributeService.GetAttributeGroupList(limit, pageNumber)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, allAG)
}

func (s *Server) getAttributeGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uuid := fmt.Sprintf("%v", unstructured["uuid"])
	attributeGroup, err := s.AttributeService.GetAttributeGroup(uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, attributeGroup)
}

func (s *Server) createAttributeGroup(ctx *gin.Context) {
	vo := &common.AddAttributeGroupVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	result, err := s.AttributeService.CreateAttributeGroup(vo, ctx.GetHeader("x-wrapper-username"))

	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, result)
}

func (s *Server) putAttributeGroup(ctx *gin.Context) {
	vo := &common.UpdateAttributeGroupVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	attributeGroup, err := s.AttributeService.UpdateAttributeGroup(vo, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, attributeGroup)
}

func (s *Server) deleteAttributeGroup(ctx *gin.Context) {
	idVO := &common.IdVO{}
	if err := ctx.BindJSON(idVO); err != nil || idVO.UUID == "" {
		api.RequestErr(ctx, errors.New("参数异常"))
		return
	}
	err := s.DeleteAttributeGroup(idVO.UUID)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllAttribute(ctx *gin.Context) {
	limit := ctx.DefaultQuery("page_size", "10")
	pageNumber := ctx.DefaultQuery("page_number", "1")

	attributeList, err := s.AttributeService.GetAttributeList(limit, pageNumber)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, attributeList)
}

func (s *Server) getAttribute(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uuid := fmt.Sprintf("%v", unstructured["uuid"])
	attribute, err := s.AttributeService.GetAttribute(uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, attribute)
}

func (s *Server) createAttribute(ctx *gin.Context) {
	vo := &common.CreateAttributeVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	result, err := s.AttributeService.CreateAttribute(vo, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, result)
}

func (s *Server) putAttribute(ctx *gin.Context) {
	vo := &common.UpdateAttributeVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	result, err := s.AttributeService.UpdateAttribute(vo, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, result)
}

func (s *Server) deleteAttribute(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uuid := fmt.Sprintf("%v", unstructured["uuid"])
	if err := s.AttributeService.DeleteAttributeInstance(uuid); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
