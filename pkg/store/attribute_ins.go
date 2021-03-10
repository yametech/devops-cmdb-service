package store

import "github.com/mindstand/gogm"

type AttributeIns struct {
	gogm.BaseNode
	AttributeCommon
	//ModelUid          string             `json:"modelUid" gogm:"name=modelUid"`
	AttributeGroupIns *AttributeGroupIns `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	// 值
	AttributeInsValue string `json:"attributeInsValue" gogm:"name=attributeInsValue"`
	CommonObj
}

func (obj *AttributeIns) Save() error {
	return GetSession(false).Save(obj)
}

//func (mg ModelGroup) List(uuid string)  {
//
//	//m := &[]ModelGroup{}
//	//err := getSession().LoadAll(m)
//
//
//}
