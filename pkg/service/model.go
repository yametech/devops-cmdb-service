package service

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/store"
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
