package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"net/http"
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
	unstructured := make(map[string]string)
	if err := ctx.BindJSON(&unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	modelGroup, err := s.ModelService.GetModelGroupInstance(unstructured["uuid"])
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, modelGroup)
}

func (s *Server) createGroup(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroup := store.ModelGroup{}
	if err := json.Unmarshal(rawData, &modelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = modelGroup.Save(s.ModelService.Session)
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
	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	uuid := fmt.Sprintf("%v", unstructured["uuid"])

	if exists := s.ModelService.CheckExists("modelGroup", uuid); exists != true {
		api.RequestErr(ctx, fmt.Errorf("group not exists"))
		return
	}
	modelGroup := store.ModelGroup{}
	if err := json.Unmarshal(rawData, &modelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroup.UUID = uuid
	err = modelGroup.Update(s.ModelService.Session)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, modelGroup)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	unstructured := make(map[string]string)
	if err := ctx.BindJSON(&unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroup, err := s.ModelService.GetModelGroupInstance(unstructured["uuid"])
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = modelGroup.Delete(s.ModelService.Session)
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
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	data := make(map[string]string, 0)
	if err := json.Unmarshal(rawData, &data); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	model, err := s.ModelService.GetModelInstance(data["uuid"])
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, model)
}

func (s *Server) createModel(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	unstructured := make(map[string]string)
	if err := ctx.BindJSON(&unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	model := store.Model{}
	modelGroupUuid := fmt.Sprintf("%v", unstructured["modelgroup"])
	if !s.ModelService.CheckExists("modelGroup", modelGroupUuid) {
		api.RequestDataErr(ctx, "groupUUID not exists", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(rawData, &model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err = s.ModelService.ChangeModelGroup(&model, modelGroupUuid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, model)
}

func (s *Server) putModel(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	unstructured := make(map[string]string)
	if err := ctx.BindJSON(&unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	originModel := &store.Model{}
	if !s.ModelService.CheckExists("model", unstructured["uuid"]) {
		api.RequestErr(ctx, fmt.Errorf("get origin model error"))
		return
	}

	model := store.Model{}
	if err := json.Unmarshal(rawData, &model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	model.UUID = originModel.UUID

	modelGroupUuid := fmt.Sprintf("%v", unstructured["modelgroup"])
	if exists := s.ModelService.CheckExists("modelGroup", modelGroupUuid); exists != true {
		api.RequestErr(ctx, fmt.Errorf("modelgroup not exists"))
		return
	}
	if err := s.ModelService.ChangeModelGroup(&model, modelGroupUuid); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, s.Model)
}

func (s *Server) deleteModel(ctx *gin.Context) {
	unstructured := make(map[string]string)
	if err := ctx.BindJSON(&unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	model, err := s.ModelService.GetModelInstance(unstructured["uuid"])
	if err != nil {
		api.RequestErr(ctx, fmt.Errorf("get model fail"))
		return
	}
	err = model.Delete(s.ModelService.Session)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
