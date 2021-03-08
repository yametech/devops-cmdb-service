package store

import (
	"github.com/mindstand/gogm"
)

type ModelGroup struct {
	gogm.BaseNode
	Uid    string   `json:"uid" gogm:"unique;name=uid"`
	Name   string   `json:"name" gogm:"name=name"`
	Models []*Model `json:"-" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (mg ModelGroup) Save() error {
	return GetSession().Save(mg)
}

//func (mg ModelGroup) List(uuid string)  {
//
//	//m := &[]ModelGroup{}
//	//err := getSession().LoadAll(m)
//
//
//}
