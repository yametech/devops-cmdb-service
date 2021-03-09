package service

import (
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type ModeService struct {
	Service
}

//func (mgs ModeService) Create(mg store.Model) error {
//	return mg.Save()
//}

func (mgs *ModeService) List() interface{} {
	m := &[]store.Model{}
	mgs.ManualQuery("match (a:Model) return a", map[string]interface{}{}, m)
	return m
}
