package service

import (
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type ModeGroupService struct {
	Service
}

func (mgs ModeGroupService) Create(mg store.ModelGroup) error {
	return mg.Save()
}

func (mgs ModeGroupService) ListByUid(uid string) interface{} {
	m := &store.ModelGroup{}
	mgs.ManualQuery("match (a:ModelGroup) where a.uid = $uid return *", map[string]interface{}{"uid": uid}, m)
	return m
}

func (mgs ModeGroupService) SomeComplexityOperate() {

}

//func (s ModeGroupService) ManualQuery(query string, properties map[string]interface{}) interface{} {
//	var respObj = new(interface{})
//	getSession().Query(query, properties,respObj)
//	return respObj
//}
