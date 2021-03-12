package service

import (
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
)

type AttributeService struct {
	AttributeGroup store.AttributeGroup
	Attribute      store.Attribute
	Session        *gogm.Session
}

func (as *AttributeService) CheckExists(modelType, uuid string) bool {
	switch modelType {
	case "attribute":
		err := as.Attribute.Get(uuid)
		if err != nil {
			return false
		}
		return true
	case "attributeGroup":
		err := as.AttributeGroup.Get(uuid)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (as *AttributeService) ChangeModelGroup(uuid string) error {
	if err := as.AttributeGroup.Get(uuid); err != nil {
		return err
	}
	query := fmt.Sprintf("match (a:Model)-[r:GroupBy]->(b:ModelGroup)where a.uuid=$uuid delete r")
	properties := map[string]interface{}{
		"uuid": as.Attribute.UUID,
	}
	_ = store.GetSession(false).Query(query, properties, nil)
	as.Attribute.AttributeGroup = &as.AttributeGroup
	if err := as.Attribute.Save(); err != nil {
		return err
	}
	return nil
}

func (as *AttributeService) CleanModelGroup(uuid string) error {
	if err := as.AttributeGroup.Get(uuid); err != nil {
		return err
	}
	as.Attribute.AttributeGroup = &as.AttributeGroup
	if err := as.AttributeGroup.Save(); err != nil {
		return err
	}
	return nil
}

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
	if err := store.GetSession(true).Query(query, properties, allAG); err != nil {
		return nil, err
	}
	for i, v := range allAG {
		attribute := make([]*store.Attribute, 0)
		if err := as.Attribute.LoadAll(&attribute, v.UUID); err != nil {
			return nil, err
		}
		allAG[i].Attributes = attribute
	}
	return &allAG, nil
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
	if err := store.GetSession(true).Query(query, properties, attributeList); err != nil {
		return nil, err
	}
	return &attributeList, nil
}
