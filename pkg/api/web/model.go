package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"strconv"
)

func (s *Server) getAllGroup(ctx *gin.Context) {
	allMG, err := s.ModelService.GetAllGroup()
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

	modelGroup, err := s.ModelService.GetModelGroup(unstructured["uuid"])
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
	modelGroup := &store.ModelGroup{}
	if err := json.Unmarshal(rawData, &modelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	err = s.ModelService.Neo4jDomain.Save(modelGroup)
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

	modelGroup := &store.ModelGroup{}
	if err := s.ModelService.Neo4jDomain.Get(modelGroup, "uuid", uuid); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	if err := json.Unmarshal(rawData, modelGroup); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroup.UUID = uuid
	err = s.ModelService.Neo4jDomain.Update(modelGroup)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, modelGroup)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
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
	vo := &common.AddModelVO{}
	err := ctx.BindJSON(vo)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	modelGroup := &store.ModelGroup{}
	if err := s.ModelService.Neo4jDomain.Get(modelGroup, "uuid", vo.ModelGroupUUID); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	model := &store.Model{}
	utils.SimpleConvert(model, vo)

	model.ModelGroup = modelGroup
	err = s.ModelService.Neo4jDomain.Save(model)
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

	model := &store.Model{}
	if err := s.ModelService.Neo4jDomain.Get(model, "uuid", unstructured["uuid"]); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	if err := json.Unmarshal(rawData, model); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	model.UUID = unstructured["uuid"]
	if err := s.ModelService.Neo4jDomain.Update(model); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, model)
}

func (s *Server) deleteModel(ctx *gin.Context) {
	unstructured := make(map[string]string)
	if err := ctx.BindJSON(&unstructured); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	model := &store.Model{}
	if err := s.ModelService.Neo4jDomain.Get(model, "uuid", unstructured["uuid"]); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	err := s.ModelService.Neo4jDomain.Delete(model)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllRelationship(ctx *gin.Context) {
	limit := ctx.DefaultQuery("page_size", "10")
	pageNumber := ctx.DefaultQuery("page_number", "1")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		api.RequestErr(ctx, errors.New("page_size参数不能小于0"))
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		api.RequestErr(ctx, errors.New("page_number"))
	}

	returnData, err := s.ModelService.GetRelationshipList(limitInt, pageNumberInt)
	if err != nil {
		fmt.Println(err)
		api.RequestOK(ctx, "")
		return
	}
	api.RequestOK(ctx, returnData)
}

func (s *Server) createRelationship(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelRelation := store.RelationshipModel{}

	if err := json.Unmarshal(rawData, &modelRelation); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	if err := s.ModelService.SaveRelationship(&modelRelation); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, modelRelation)
}

func (s *Server) updateRelationship(ctx *gin.Context) {
	vo := &common.RelationshipModelUpdateVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	if err := s.ModelService.UpdateRelationship(vo, ""); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, nil)
}

func (s *Server) deleteRelationship(ctx *gin.Context) {
	idVO := &common.IdVO{}
	if _ = ctx.BindJSON(idVO); idVO.UUID == "" {
		api.RequestErr(ctx, errors.New("uuid参数不能为空"))
		return
	}

	err := s.ModelService.DeleteRelationship(idVO.UUID)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}
