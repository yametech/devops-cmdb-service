package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
	"strings"
)

func (s *Server) getAllGroup(ctx *gin.Context) {
	allMG, _ := s.ModelService.GetAllModelGroup()
	//if err != nil {
	//	api.RequestErr(ctx, err)
	//	return
	//}
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
	modelGroupVO := &common.AddModelGroupVO{}
	if err := ctx.ShouldBindJSON(modelGroupVO); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroupVO.Uid = strings.TrimSpace(modelGroupVO.Uid)
	modelGroupVO.Name = strings.TrimSpace(modelGroupVO.Name)

	result, err := s.ModelService.CreateModelGroup(modelGroupVO, "")
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, result)
}

func (s *Server) putGroup(ctx *gin.Context) {
	modelGroupVO := &common.AddModelGroupVO{}
	if err := ctx.ShouldBindJSON(modelGroupVO); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	modelGroupVO.Name = strings.TrimSpace(modelGroupVO.Name)

	result, err := s.ModelService.UpdateModelGroup(modelGroupVO, "")
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, result)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	// TODO
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

	model, err := s.ModelService.GetModel(data["uuid"])
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, model)
}

func (s *Server) createModel(ctx *gin.Context) {
	vo := &common.AddModelVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	vo.Uid = strings.TrimSpace(vo.Uid)
	vo.Name = strings.TrimSpace(vo.Name)

	result, err := s.ModelService.CreateModel(vo, "")
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, result)
}

func (s *Server) putModel(ctx *gin.Context) {
	// TODO
	api.RequestOK(ctx, "")
}

func (s *Server) deleteModel(ctx *gin.Context) {
	idVO := &common.IdVO{}
	if err := ctx.ShouldBindJSON(idVO); err != nil || idVO.UUID == "" {
		api.RequestErr(ctx, errors.New("参数异常"))
		return
	}

	err := s.ModelService.DeleteModel(idVO.UUID, "")
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
		api.RequestOK(ctx, []map[string]string{})
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
