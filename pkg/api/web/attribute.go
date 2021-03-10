package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

func (s *Server) getAllAttributeGroup(ctx *gin.Context) {
	allAG := make([]store.AttributeGroup, 0)
	query := fmt.Sprintf("match (a:AttributeGroup) return a")
	err := store.GetSession(true).Query(query, map[string]interface{}{}, &allAG)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	for i, v := range allAG {
		attributes := make([]*store.Attribute, 0)
		if err := s.Attribute.LoadAll(&attributes, v.Uid); err != nil {
			api.RequestErr(ctx, err)
			return
		}
		allAG[i].Attributes = attributes
	}
	api.RequestOK(ctx, allAG)
}

func (s *Server) getAttributeGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.AttributeGroup.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.AttributeGroup)
}

func (s *Server) createAttributeGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	if err := json.Unmarshal(rawData, &s.AttributeGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = s.AttributeGroup.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.AttributeGroup)
}

func (s *Server) putAttributeGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uid := ctx.Param("uid")
	originAG := &store.AttributeGroup{}
	err = originAG.Get(uid)
	if err != nil {
		api.RequestErr(ctx, fmt.Errorf("get origin attributegroup error"))
		return
	}
	if err := json.Unmarshal(rawData, &s.AttributeGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.AttributeGroup.UUID = originAG.UUID
	err = s.AttributeGroup.Update()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.AttributeGroup)
}

func (s *Server) deleteAttributeGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	if err := s.AttributeGroup.Get(uid); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err := s.AttributeGroup.Delete()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllAttribute(ctx *gin.Context) {
	allAttribute := &[]store.Attribute{}
	query := fmt.Sprintf("match (a:Attribute) return a")
	properties := map[string]interface{}{}
	err := store.GetSession(true).Query(query, properties, allAttribute)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, allAttribute)
}

func (s *Server) getAttribute(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.Attribute.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.Attribute)
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
	attributeGroupUid := fmt.Sprintf("%v", unstructured["attributegroupuid"])

	if err := s.AttributeGroup.Get(attributeGroupUid); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	if err := json.Unmarshal(rawData, &s.Attribute); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.Attribute.AttributeGroup = &s.AttributeGroup
	err = s.Attribute.Save()
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
	uid := ctx.Param("uid")
	originAttribute := &store.Attribute{}
	err = originAttribute.Get(uid)
	if err != nil {
		api.RequestErr(ctx, fmt.Errorf("get origin attribute error"))
		return
	}
	if err := json.Unmarshal(rawData, &s.Attribute); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.Attribute.UUID = originAttribute.UUID
	err = s.Attribute.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.Attribute)
}

func (s *Server) deleteAttribute(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.Attribute.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = s.Attribute.Delete()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
