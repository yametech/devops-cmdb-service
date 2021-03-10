package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

func (s *Server) getAllGroup(ctx *gin.Context) {
	allMG := &[]store.ModelGroup{}
	query := fmt.Sprintf("match (a:ModelGroup) return a")
	properties := map[string]interface{}{}
	err := store.GetSession(true).Query(query, properties, allMG)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, allMG)
}

func (s *Server) getGroup(ctx *gin.Context) {
	originMG := &store.ModelGroup{}
	uid := ctx.Param("uid")
	err := originMG.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, originMG)
}

func (s *Server) createGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroup := &store.ModelGroup{}
	if err := json.Unmarshal(rawData, modelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	err = modelGroup.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, modelGroup)
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
	newMG := &store.ModelGroup{}
	if err := json.Unmarshal(rawData, newMG); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	newMG.UUID = originMG.UUID
	err = newMG.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, newMG)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	originMG := &store.ModelGroup{}
	err := originMG.Delete(uid)
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
	err := s.model.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.model)
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
	if err := s.modelGroup.Get(modelGroupUid); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	model := &store.Model{}
	if err := json.Unmarshal(rawData, model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	model.ModelGroup = s.modelGroup
	err = model.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	//s.modelGroup.Models = append(s.modelGroup.Models, model)
	//err = s.modelGroup.Save()
	//if err != nil {
	//	api.RequestErr(ctx, err)
	//	return
	//}
	api.RequestOK(ctx, model)
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
		api.RequestErr(ctx, err)
		return
	}
	newModel := &store.Model{}
	if err := json.Unmarshal(rawData, newModel); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	newModel.UUID = originModel.UUID
	err = newModel.Save()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, newModel)
}

func (s *Server) deleteModel(ctx *gin.Context) {
	uid := ctx.Param("uid")
	model := &store.Model{}
	err := model.Delete(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
