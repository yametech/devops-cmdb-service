package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

func (s *Server) getAllGroup(ctx *gin.Context) {
	limit := ctx.DefaultQuery("page_size", "10")
	pageNumber := ctx.DefaultQuery("page_number", "1")

	allMG, err := s.ModelService.GetGroupList(limit, pageNumber)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, allMG)
}

func (s *Server) getGroup(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	err := s.ModelGroup.Get(uuid)
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
	uuid := ctx.Param("uuid")
	if exists := s.ModelService.CheckExists("modelGroup", uuid); exists != true {
		api.RequestErr(ctx, fmt.Errorf("group not exists"))
		return
	}
	if err := json.Unmarshal(rawData, &s.ModelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.ModelGroup.UUID = uuid
	err = s.ModelGroup.Update()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, s.ModelGroup)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	err := s.ModelGroup.Delete(uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllModel(ctx *gin.Context) {
	limit := ctx.DefaultQuery("page_size", "10")
	pageNumber := ctx.DefaultQuery("page_number", "1")

	allModel, err := s.ModelService.GetModelList(limit, pageNumber)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, allModel)
}

func (s *Server) getModel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	err := s.Model.Get(uuid)
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

	modelGroupUuid := fmt.Sprintf("%v", unstructured["modelgroup"])
	if exists := s.ModelService.CheckExists("modelGroup", modelGroupUuid); exists != true {
		api.RequestErr(ctx, fmt.Errorf("groupUUID not exists"))
		return
	}
	if err := json.Unmarshal(rawData, &s.Model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = s.ModelService.ChangeModelGroup(modelGroupUuid)
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
	uuid := ctx.Param("uuid")
	originModel := &store.Model{}
	err = originModel.Get(uuid)
	if err != nil {
		api.RequestErr(ctx, fmt.Errorf("get origin model error"))
		return
	}
	if err := json.Unmarshal(rawData, &s.Model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	s.Model.UUID = originModel.UUID
	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroupUuid := fmt.Sprintf("%v", unstructured["modelgroup"])
	if exists := s.ModelService.CheckExists("modelGroup", modelGroupUuid); exists != true {
		api.RequestErr(ctx, fmt.Errorf("modelgroup not exists"))
		return
	}
	if err := s.ModelService.ChangeModelGroup(modelGroupUuid); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, s.Model)
}

func (s *Server) deleteModel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	err := s.Model.Delete(uuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
