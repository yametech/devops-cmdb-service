package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

func (s *Server) getAllGroup(ctx *gin.Context) {
	allMG := make([]store.ModelGroup, 0)
	query := fmt.Sprintf("match (a:ModelGroup) return a")
	err := store.GetSession(true).Query(query, map[string]interface{}{}, &allMG)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	for i, v := range allMG {
		models := make([]*store.Model, 0)
		if err := s.Model.LoadAll(&models, v.Uid); err != nil {
			api.RequestErr(ctx, err)
			return
		}
		allMG[i].Models = models
	}
	api.RequestOK(ctx, allMG)
}

func (s *Server) getGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.ModelGroup.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.ModelGroup)
}

func (s *Server) createGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	if err := json.Unmarshal(rawData, &s.ModelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = s.ModelGroup.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.ModelGroup)
}

func (s *Server) putGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uid := ctx.Param("uid")
	originMG := &store.ModelGroup{}
	err = originMG.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	if err := json.Unmarshal(rawData, &s.ModelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.ModelGroup.UUID = originMG.UUID
	err = s.ModelGroup.Update()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.ModelGroup)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.ModelGroup.Delete(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllModel(ctx *gin.Context) {
	allModel := &[]store.Model{}
	query := fmt.Sprintf("match (a:Model) return a")
	properties := map[string]interface{}{}
	err := store.GetSession(true).Query(query, properties, allModel)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, allModel)
}

func (s *Server) getModel(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.Model.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.Model)
}

func (s *Server) createModel(ctx *gin.Context) {
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
	modelGroupUid := fmt.Sprintf("%v", unstructured["modelgroupuid"])

	if err := s.ModelGroup.Get(modelGroupUid); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	if err := json.Unmarshal(rawData, &s.Model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.Model.ModelGroup = &s.ModelGroup
	err = s.Model.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, s.Model)
}

func (s *Server) putModel(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uid := ctx.Param("uid")
	originModel := &store.Model{}
	err = originModel.Get(uid)
	if err != nil {
		api.RequestErr(ctx, fmt.Errorf("get origin model error"))
		return
	}
	if err := json.Unmarshal(rawData, &s.Model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.Model.UUID = originModel.UUID
	err = s.Model.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.Model)
}

func (s *Server) deleteModel(ctx *gin.Context) {
	uid := ctx.Param("uid")
	err := s.Model.Delete(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
