package service

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
)

type ModelService struct {
	Model      store.Model
	ModelGroup store.ModelGroup
}

func (as *ModelService) CheckExists(modelType, uuid string) bool {
	switch modelType {
	case "model":
		err := as.Model.Get(uuid)
		if err != nil {
			return false
		}
		return true
	case "modelGroup":
		err := as.ModelGroup.Get(uuid)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (as *ModelService) ChangeModelGroup(uuid string) error {
	if err := as.ModelGroup.Get(uuid); err != nil {
		return err
	}
	query := fmt.Sprintf("match (a:Model)-[r:GroupBy]->(b:ModelGroup)where a.uuid=$uuid delete r")
	properties := map[string]interface{}{
		"uuid": as.Model.UUID,
	}
	_ = store.GetSession(false).Query(query, properties, nil)
	as.Model.ModelGroup = &as.ModelGroup
	if err := as.Model.Save(); err != nil {
		return err
	}
	return nil
}

func (as *ModelService) CleanModelGroup(uuid string) error {
	if err := as.ModelGroup.Get(uuid); err != nil {
		return err
	}
	as.Model.ModelGroup = &as.ModelGroup
	if err := as.Model.Save(); err != nil {
		return err
	}
	return nil
}

func (as *ModelService) GetGroupList(limit, pageNumber string) (*[]store.ModelGroup, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}
	allMG := make([]store.ModelGroup, 0)
	query := fmt.Sprintf("match (a:ModelGroup) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	err = store.GetSession(true).Query(query, properties, &allMG)
	if err != nil {
		return nil, err
	}
	for i, v := range allMG {
		models, err := as.Model.LoadAll(v.UUID)
		if err != nil {
			return nil, err
		}
		allMG[i].Models = models
	}
	return &allMG, nil
}

func (as *ModelService) GetModelList(limit string, pageNumber string) (*[]store.Model, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}
	allModel := make([]store.Model, 0)
	query := fmt.Sprintf("match (a:Model) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := store.GetSession(true).Query(query, properties, allModel); err != nil {
		return nil, err
	}
	return &allModel, nil
}

func (as *ModelService) GetModelGroupInstance(uuid string) (*store.ModelGroup, error) {
	modelGroup := &store.ModelGroup{}
	if err := modelGroup.Get(uuid); err != nil {
		return nil, err
	}
	models, err := as.Model.LoadAll(modelGroup.UUID)
	if err != nil {
		return nil, err
	}
	modelGroup.Models = models
	return modelGroup, nil
}

func (as *ModelService) GetModelInstance(uuid string) (*store.Model, error) {
	model := &store.Model{}
	if err := model.Get(uuid); err != nil {
		return nil, err
	}
	return model, nil
}
