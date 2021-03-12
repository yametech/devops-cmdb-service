package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/store"
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
	attributeGroup, err := s.AttributeService.GetAttributeGroupInstance(uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, attributeGroup)
}

func (s *Server) createAttributeGroup(ctx *gin.Context) {
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
	modelUUID := fmt.Sprintf("%v", unstructured["modeluuid"])

	model, err := s.ModelService.GetModelInstance(modelUUID)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	attributeGroup, err := s.AttributeService.CreateAttributeGroup(rawData, model)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, attributeGroup)
}

func (s *Server) putAttributeGroup(ctx *gin.Context) {
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

	attributeGroup, err := s.AttributeService.UpdateAttributeGroupInstance(rawData, uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, attributeGroup)
}

func (s *Server) deleteAttributeGroup(ctx *gin.Context) {
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
	err = s.AttributeService.DeleteAttributeGroupInstance(uuid)
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
	uuid := ctx.Param("uuid")
	attribute, err := s.AttributeService.GetAttributeInstance(uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, attribute)
}

func (s *Server) createAttribute(ctx *gin.Context) {
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
	attributeGroupUuid := fmt.Sprintf("%v", unstructured["attributegroupuuid"])
	attribute := &store.Attribute{}
	if err := json.Unmarshal(rawData, attribute); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = s.AttributeService.ChangeModelGroup(attribute, attributeGroupUuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, s.Attribute)
}

func (s *Server) putAttribute(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uuid := ctx.Param("uuid")
	err = s.AttributeService.UpdateAttributeInstance(rawData, uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, s.Attribute)
}

func (s *Server) deleteAttribute(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if err := s.AttributeService.DeleteAttributeInstance(uuid); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
