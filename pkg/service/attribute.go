package service

import "github.com/yametech/devops-cmdb-service/pkg/store"

type AttributeService struct {
	AttributeGroup store.AttributeGroup
	Attribute      store.Attribute
}

func (as *AttributeService) CheckExists(obj interface{}) error {
	return nil
}
