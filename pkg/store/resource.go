package store

import (
	"github.com/mindstand/gogm"
)

type Resource struct {
	gogm.BaseNode
	Uid      string `json:"uid" gogm:"unique;name=uid"`
	Name     string `json:"name" gogm:"name=name"`
	ModelUid string `json:"modelUid" gogm:"name=modelUid"`
	Models   Model  `json:"-" gogm:"direction=outgoing;relationship=Instance"`
	CommonObj
}

func (obj *Resource) Save() error {
	return GetSession(false).Save(obj)
}

//func (mg ModelGroup) List(uuid string)  {
//
//	//m := &[]ModelGroup{}
//	//err := getSession().LoadAll(m)
//
//
//}
