package store

import (
	"github.com/mindstand/gogm"
)

type AttributeIns struct {
	gogm.BaseNode
	Uid               string             `json:"uid" gogm:"unique;name=uid"`
	Name              string             `json:"name" gogm:"name=name"`
	ModelUid          string             `json:"modelUid" gogm:"name=modelUid"`
	AttributeGroupIns *AttributeGroupIns `json:"-" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (obj AttributeIns) Save() error {
	return GetSession().Save(obj)
}

//func (mg ModelGroup) List(uuid string)  {
//
//	//m := &[]ModelGroup{}
//	//err := getSession().LoadAll(m)
//
//
//}
