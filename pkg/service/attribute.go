package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
	"sync"
)

type AttributeService struct {
	Service
	store.Neo4jDomain
	Mutex sync.Mutex
}

//func (as *AttributeService) CheckExists(modelType, uuid string) bool {
//	switch modelType {
//	case "attribute":
//		err := as.Attribute.Get(as.Session, uuid)
//		if err != nil {
//			return false
//		}
//		return true
//	case "attributeGroup":
//		err := as.AttributeGroup.Get(as.Session, uuid)
//		if err != nil {
//			return false
//		}
//		return true
//	}
//	return false
//}

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

	as.ManualExecute(query, properties)

	attribute.AttributeGroup = attributeGroup
	if err := as.Save(attribute); err != nil {
		return err
	}
	return nil
}

//func (as *AttributeService) CleanModelGroup(uuid string) error {
//	attributeGroup := &store.AttributeGroup{}
//	if err := as.Neo4jDomain.Get(attributeGroup, "uuid", uuid); err != nil {
//		return err
//	}
//	as.Attribute.AttributeGroup = &as.AttributeGroup
//	if err := as.AttributeGroup.Save(as.Session); err != nil {
//		return err
//	}
//	return nil
//}

func (as *AttributeService) GetAttributeGroupList(limit string, pageNumber string) (*[]store.AttributeGroup, error) {
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

	if err := as.ManualQuery(query, properties, &allAG); err != nil {
		return nil, err
	}
	/*for i, v := range allAG {
		attributes, err := as.Attribute.LoadAll(as.Session, v.UUID)
		if err != nil {
			continue
		}
		allAG[i].Attributes = *attributes
	}*/
	return &allAG, nil
}

func (as *AttributeService) GetAttributeGroupInstance(uuid string) (*store.AttributeGroup, error) {
	attributeGroup := &store.AttributeGroup{}
	if err := store.GetSession(true).Load(attributeGroup, uuid); err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) CreateAttributeGroup(rawData []byte, model *store.Model) (*store.AttributeGroup, error) {
	attributeGroup := &store.AttributeGroup{}
	if err := json.Unmarshal(rawData, attributeGroup); err != nil {
		return nil, err
	}

	attributeGroup.Model = model
	attributeGroup.ModelUid = model.Uid
	err := store.GetSession(false).SaveDepth(attributeGroup, 2)
	if err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) UpdateAttributeGroupInstance(rawData []byte, uuid string) (*store.AttributeGroup, error) {

	unstructured := make(map[string]interface{})
	if err := json.Unmarshal(rawData, &unstructured); err != nil {
		return nil, err
	}

	attributeGroup := &store.AttributeGroup{}
	if err := json.Unmarshal(rawData, attributeGroup); err != nil {
		return nil, err
	}
	attributeGroup.UUID = uuid

	err := store.GetSession(false).Save(attributeGroup)
	if err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) DeleteAttributeGroup(uuid string) error {

	attributeGroup := &store.AttributeGroup{}
	if err := as.Neo4jDomain.Get(attributeGroup, "uuid", uuid); err != nil {
		return err
	}
	if err := as.Neo4jDomain.Delete(attributeGroup); err != nil {
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

func (as *AttributeService) GetAttribute(uuid string) (*store.Attribute, error) {

	attribute := &store.Attribute{}
	if err := as.Neo4jDomain.Get(attribute, "uuid", uuid); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) UpdateAttribute(rawData []byte, uuid string) (*store.Attribute, error) {
	_, err := as.GetAttribute(uuid)
	if err != nil {
		return nil, errors.New("attribute not exists")
	}
	attribute := &store.Attribute{}
	if err := json.Unmarshal(rawData, attribute); err != nil {
		return nil, err
	}
	attribute.UUID = uuid

	if err := as.Neo4jDomain.Update(attribute); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) DeleteAttributeInstance(uuid string) error {
	attribute, err := as.GetAttribute(uuid)
	if err != nil {
		return err
	}
	if err := as.Neo4jDomain.Delete(attribute); err != nil {
		return err
	}
	return nil
}
