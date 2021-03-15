package service

import (
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
	"sync"
)

type ModelService struct {
	Model      store.Model
	ModelGroup store.ModelGroup
	Session    *gogm.Session
	mutex      sync.Mutex
}

func (ms *ModelService) CheckExists(modelType, uuid string) bool {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	switch modelType {
	case "model":
		model := store.Model{}
		err := model.Get(ms.Session, uuid)
		if err != nil {
			return false
		}
		return true
	case "modelGroup":
		modelGroup := store.ModelGroup{}
		err := modelGroup.Get(ms.Session, uuid)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (ms *ModelService) ChangeModelGroup(model *store.Model, uuid string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	modelGroup := store.ModelGroup{}
	if err := modelGroup.Get(ms.Session, uuid); err != nil {
		return err
	}
	query := fmt.Sprintf("match (a:Model)-[r:GroupBy]->(b:ModelGroup)where a.uuid=$uuid delete r")
	properties := map[string]interface{}{
		"uuid": model.UUID,
	}
	_ = store.GetSession(false).Query(query, properties, nil)
	model.ModelGroup = &modelGroup
	if err := model.Save(ms.Session); err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) CleanModelGroup(uuid string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if err := ms.ModelGroup.Get(ms.Session, uuid); err != nil {
		return err
	}
	ms.Model.ModelGroup = &ms.ModelGroup
	if err := ms.Model.Save(ms.Session); err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) GetGroupList(limit, pageNumber string) (*[]store.ModelGroup, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
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
	err = ms.Session.Query(query, properties, &allMG)
	//if err != nil {
	//	return nil, err
	//}
	for i, v := range allMG {
		models, err := ms.Model.LoadAll(ms.Session, v.UUID)
		if err != nil {
			return nil, err
		}
		allMG[i].Models = models
	}
	return &allMG, nil
}

func (ms *ModelService) GetModelGroupInstance(uuid string) (*store.ModelGroup, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	modelGroup := &store.ModelGroup{}
	if err := modelGroup.Get(ms.Session, uuid); err != nil {
		return nil, err
	}
	models, err := ms.Model.LoadAll(ms.Session, modelGroup.UUID)
	if err != nil {
		return nil, err
	}
	modelGroup.Models = models
	return modelGroup, nil
}

func (ms *ModelService) GetModelList(limit string, pageNumber string) (*[]store.Model, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
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
	if err := store.GetSession(true).Query(query, properties, &allModel); err != nil {
		return nil, err
	}
	return &allModel, nil
}

func (ms *ModelService) GetModelInstance(uuid string) (*store.Model, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	model := &store.Model{}
	if err := model.Get(ms.Session, uuid); err != nil {
		return nil, err
	}
	agInstance := &store.AttributeGroup{}
	attributeInstance := &store.Attribute{}
	attributeGroup, err := agInstance.LoadAll(ms.Session, model.UUID)
	if err != nil {
		return model, nil
	}
	for i, _ := range attributeGroup {
		attribute, err := attributeInstance.LoadAll(ms.Session, attributeGroup[i].UUID)
		if err != nil {
			continue
		}
		attributeGroup[i].Attributes = *attribute
	}
	model.AttributeGroups = attributeGroup
	return model, nil
}

func (ms *ModelService) GetRelationshipList(limit string, pageNumber string) (*[]store.RelationshipModel, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}
	allModel := make([]store.RelationshipModel, 0)
	query := fmt.Sprintf("match (a:RelationshipModel) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := store.GetSession(true).Query(query, properties, &allModel); err != nil {
		//return nil, err
	}
	return &allModel, err
}

func (ms *ModelService) GetRelationshipInstance(uuid string) (*store.RelationshipModel, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	relation := &store.RelationshipModel{}
	if err := relation.Get(ms.Session, uuid); err != nil {
		return nil, err
	}
	return relation, nil
}

func (ms *ModelService) SaveRelationship(relation *store.RelationshipModel) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	err := relation.Save(ms.Session)
	if err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) UpdateRelation(relation *store.RelationshipModel, uuid string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	relation.UUID = uuid
	err := relation.Update(ms.Session)
	if err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) DeleteRelationship(relation *store.RelationshipModel) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	err := relation.Delete(ms.Session)
	if err != nil {
		return err
	}
	return nil
}
