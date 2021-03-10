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
	allMG := &[]store.ModelGroup{}
	query := fmt.Sprintf("match (a:ModelGroup) return a")
	properties := map[string]interface{}{}
	err := store.GetSession(true).Query(query, properties, allMG)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, allMG)
}

func (s *Server) getGroup(ctx *gin.Context) {
	originMG := &store.ModelGroup{}
	uid := ctx.Param("uid")
	err := originMG.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, originMG)
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
	ctx.JSON(http.StatusOK, modelGroup)
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
	ctx.JSON(http.StatusOK, newMG)
}

func (s *Server) deleteGroup(ctx *gin.Context) {
	uid := ctx.Param("uid")
	originMG := &store.ModelGroup{}
	err := originMG.Delete(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, "")
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
	ctx.JSON(http.StatusOK, allModel)
}

func (s *Server) getModel(ctx *gin.Context) {
	model := &store.Model{}
	uid := ctx.Param("uid")
	err := model.Get(uid)
	if err != nil {
		api.RequestErr(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, model)
}

func (s *Server) createModel(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, modelGroup)
}

