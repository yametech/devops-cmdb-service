package store

import (
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
)

type AttributeGroupIns struct {
	gogm.BaseNode
	Uid          string          `json:"uid" gogm:"name=uid"`
	Name         string          `json:"name" gogm:"name=name"`
	ModelUid     string          `json:"modelUid" gogm:"index;name=modelUid"`
	Resource     *Resource       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	AttributeIns []*AttributeIns `json:"attributeIns" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (obj *AttributeGroupIns) AddAttributeIns(target *AttributeIns) {

	if target == nil {
		return
	}

	if obj.AttributeIns == nil {
		obj.AttributeIns = make([]*AttributeIns, 0)
	}

	for _, attributeIns := range obj.AttributeIns {
		if attributeIns.Uid == target.Uid {
			return
		}
	}

	obj.AttributeIns = append(obj.AttributeIns, target)
}
