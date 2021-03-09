package store

import (
	"github.com/mindstand/gogm"
)

type AttributeGroupIns struct {
	gogm.BaseNode
	Uid          string          `json:"uid" gogm:"unique;name=uid"`
	Name         string          `json:"name" gogm:"name=name"`
	ModelUid     string          `json:"modelUid" gogm:"name=modelUid"`
	Resource     *Resource       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	AttributeIns []*AttributeIns `json:"attributeIns" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (obj *AttributeGroupIns) Save() error {
	return GetSession(false).Save(obj)
}

func (obj *AttributeGroupIns) AddAttributeIns(target *AttributeIns) {

	if target == nil {
		return
	}

	if obj.AttributeIns == nil {
		obj.AttributeIns = make([]*AttributeIns, 0)
	}

	for _, attributeIns := range obj.AttributeIns {
		if attributeIns.UUID == target.UUID {
			return
		}
	}

	obj.AttributeIns = append(obj.AttributeIns, target)
}
