package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"strconv"
	"strings"
)

func (s *Server) getAllGroup(ctx *gin.Context) {
	allMG, _ := s.ModelService.GetAllModelGroup()
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

	result, err := s.ModelService.CreateModelGroup(modelGroupVO, ctx.GetHeader("x-wrapper-username"))
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

	result, err := s.ModelService.UpdateModelGroup(modelGroupVO, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, result)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	idVO := &common.IdVO{}
	if err := ctx.BindJSON(idVO); err != nil || idVO.UUID == "" {
		api.RequestErr(ctx, errors.New("参数异常"))
		return
	}
	err := s.DeleteModelGroup(idVO.UUID)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllModel(ctx *gin.Context) {
	allModel, err := s.ModelService.GetSimpleModelList()
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

	result, err := s.ModelService.CreateModel(vo, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, result)
}

func (s *Server) putModel(ctx *gin.Context) {
	vo := &common.UpdateModelVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	err := s.ModelService.UpdateModel(vo, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) deleteModel(ctx *gin.Context) {
	idVO := &common.IdVO{}
	if err := ctx.ShouldBindJSON(idVO); err != nil || idVO.UUID == "" {
		api.RequestErr(ctx, errors.New("参数异常"))
		return
	}

	err := s.ModelService.DeleteModel(idVO.UUID, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	api.RequestOK(ctx, "")
}

func (s *Server) getAllRelationship(ctx *gin.Context) {
	limit := ctx.DefaultQuery("pageSize", "10000")
	pageNumber := ctx.DefaultQuery("current", "1")
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
	vo := &common.CreateRelationshipModelVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}
	result, err := s.ModelService.SaveRelationship(vo, ctx.GetHeader("x-wrapper-username"))
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}

	api.RequestOK(ctx, result)
}

func (s *Server) updateRelationship(ctx *gin.Context) {
	vo := &common.UpdateRelationshipModelVO{}
	if err := ctx.ShouldBindJSON(vo); err != nil {
		api.RequestErr(ctx, err)
		return
	}

	if err := s.ModelService.UpdateRelationship(vo, ctx.GetHeader("x-wrapper-username")); err != nil {
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
