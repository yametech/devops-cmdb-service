package service

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type AttributeService struct {
	AttributeGroup store.AttributeGroup
	Attribute      store.Attribute
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

