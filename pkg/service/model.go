package service

import (
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type ModelService struct {
	Model      store.Model
	ModelGroup store.ModelGroup
}

func (as *ModelService) CheckExists(modelType, uid string) bool {
	switch modelType {
	case "Model":
		err := as.Model.Get(uid)
		if err != nil {
			return false
		}
		return true
	case "ModelGroup":
		err := as.ModelGroup.Get(uid)
		if err != nil {
			return false
		}
		return true

	}
	return false
}
