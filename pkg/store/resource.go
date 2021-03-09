package store

import (
	"github.com/mindstand/gogm"
)

type Resource struct {
	gogm.BaseNode
	ModelUid          string               `json:"modelUid" gogm:"name=modelUid"`
	ModelName         string               `json:"name" gogm:"name=modelName"`
	Models            *Model               `json:"-" gogm:"direction=outgoing;relationship=Instance"`
	AttributeGroupIns []*AttributeGroupIns `json:"attributeGroupIns" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (obj *Resource) Save() error {
	return GetSession(false).Save(obj)
}

func (obj *Resource) AddAttributeGroupIns(target *AttributeGroupIns) {

	if target == nil {
		return
	}

	if obj.AttributeGroupIns == nil {
		obj.AttributeGroupIns = make([]*AttributeGroupIns, 0)
	}

	for _, attributeGroupIns := range obj.AttributeGroupIns {
		if attributeGroupIns.UUID == target.UUID {
			return
		}
	}

	obj.AttributeGroupIns = append(obj.AttributeGroupIns, target)
}
