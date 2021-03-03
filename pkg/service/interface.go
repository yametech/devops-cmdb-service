package service

import (
	"github.com/yametech/devops-cmdb-service/pkg/core"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type Service struct {
	store.IStore
}

type fakeService struct {
	Service
}

func (f *fakeService) GetMember(uuid string) core.IObject {
	obj, err := f.Get(uuid)
	if err != nil {
		//
	}
	return obj
}
