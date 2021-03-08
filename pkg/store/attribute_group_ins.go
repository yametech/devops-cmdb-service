package store

import (
	"github.com/mindstand/gogm"
)

type AttributeGroupIns struct {
	gogm.BaseNode
	Uid      string    `json:"uid" gogm:"unique;name=uid"`
	Name     string    `json:"name" gogm:"name=name"`
	ModelUid string    `json:"modelUid" gogm:"name=modelUid"`
	Resource *Resource `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	CommonObj
}

func (obj AttributeGroupIns) Save() error {
	return GetSession().Save(obj)
}

//func (mg ModelGroup) List(uuid string)  {
//
//	//m := &[]ModelGroup{}
//	//err := getSession().LoadAll(m)
//
//
//}
