package service

import (
	"encoding/json"
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
	"sync"
)

type AttributeService struct {
	AttributeGroup store.AttributeGroup
	Attribute      store.Attribute
	Session        *gogm.Session
	Mutex          sync.Mutex
}

func (as *AttributeService) CheckExists(modelType, uuid string) bool {
	switch modelType {
	case "attribute":
		err := as.Attribute.Get(as.Session, uuid)
		if err != nil {
			return false
		}
		return true
	case "attributeGroup":
		err := as.AttributeGroup.Get(as.Session, uuid)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (as *AttributeService) ChangeModelGroup(attribute *store.Attribute, uuid string) error {
	attributeGroup, err := as.GetAttributeGroupInstance(uuid)
	if err != nil {
		return err
	}
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	query := fmt.Sprintf("match (a:Attribute)-[r:GroupBy]->(b:AttributeGroup)where a.uuid=$uuid delete r")
	properties := map[string]interface{}{
		"uuid": attribute.UUID,
	}
	_ = as.Session.Query(query, properties, nil)
	attribute.AttributeGroup = attributeGroup
	if err := attribute.Save(as.Session); err != nil {
		return err
	}
	return nil
}

func (as *AttributeService) CleanModelGroup(uuid string) error {
	if err := as.AttributeGroup.Get(as.Session, uuid); err != nil {
		return err
	}
	as.Attribute.AttributeGroup = &as.AttributeGroup
	if err := as.AttributeGroup.Save(as.Session); err != nil {
		return err
	}
	return nil
}

func (as *AttributeService) GetAttributeGroupList(limit string, pageNumber string) (*[]store.AttributeGroup, error) {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}

	allAG := make([]store.AttributeGroup, 0)
	query := fmt.Sprintf("match (a:AttributeGroup) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := as.Session.Query(query, properties, &allAG); err != nil {
		return nil, err
	}
	for i, v := range allAG {
		attributes, err := as.Attribute.LoadAll(as.Session, v.UUID)
		if err != nil {
			continue
		}
		allAG[i].Attributes = *attributes
	}
	return &allAG, nil
}

func (as *AttributeService) GetAttributeGroupInstance(uuid string) (*store.AttributeGroup, error) {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	attributeGroup := &store.AttributeGroup{}
	if err := attributeGroup.Get(as.Session, uuid); err != nil {
		return nil, err
	}

	attribute, err := as.Attribute.LoadAll(as.Session, attributeGroup.UUID)
	if err != nil {
		return attributeGroup, nil
	}
	attributeGroup.Attributes = *attribute
	return attributeGroup, nil
}

func (as *AttributeService) CreateAttributeGroup(rawData []byte, model *store.Model) (*store.AttributeGroup, error) {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	attributeGroup := &store.AttributeGroup{}
	if err := json.Unmarshal(rawData, attributeGroup); err != nil {
		return nil, err
	}

	attributeGroup.Model = model
	attributeGroup.ModelUid = model.Uid
	err := attributeGroup.Save(as.Session)
	if err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) UpdateAttributeGroupInstance(rawData []byte, uuid string) (*store.AttributeGroup, error) {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		return nil, err
	}

	attributeGroup := &store.AttributeGroup{}
	if err := json.Unmarshal(rawData, attributeGroup); err != nil {
		return nil, err
	}
	attributeGroup.UUID = uuid
	err := attributeGroup.Update(as.Session)
	if err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) DeleteAttributeGroupInstance(uuid string) error {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	if err := as.AttributeGroup.Get(as.Session, uuid); err != nil {
		return err
	}
	if err := as.AttributeGroup.Delete(as.Session); err != nil {
		return err
	}
	return nil
}

func (as *AttributeService) GetAttributeList(limit string, pageNumber string) (*[]store.Attribute, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}

	attributeList := make([]store.Attribute, 0)
	query := fmt.Sprintf("match (a:Attribute) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := store.GetSession(true).Query(query, properties, &attributeList); err != nil {
		return nil, err
	}
	return &attributeList, nil
}

func (as *AttributeService) GetAttributeInstance(uuid string) (*store.Attribute, error) {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()
	attribute := &store.Attribute{}
	if err := attribute.Get(as.Session, uuid); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) UpdateAttributeInstance(rawData []byte, uuid string) error {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	if !as.CheckExists("attribute", uuid) {
		return fmt.Errorf("attribute not exists")
	}
	if err := json.Unmarshal(rawData, &as.Attribute); err != nil {
		return err
	}
	as.Attribute.UUID = uuid
	if err := as.Attribute.Save(as.Session); err != nil {
		return err
	}
	return nil
}

func (as *AttributeService) DeleteAttributeInstance(uuid string) error {
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	if err := as.Attribute.Get(as.Session, uuid); err != nil {
		return err
	}
	if err := as.Attribute.Delete(as.Session); err != nil {
		return err
	}
	return nil
}
