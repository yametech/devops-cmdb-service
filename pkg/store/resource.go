package store

import (
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
)

type Resource struct {
	gogm.BaseNode
	ModelUid          string               `json:"modelUid" gogm:"index;name=modelUid"`
	ModelName         string               `json:"modelName" gogm:"name=modelName"`
	Models            *Model               `json:"-" gogm:"direction=outgoing;relationship=Instance"`
	Resource          *Resource            `json:"-" gogm:"direction=both;relationship=Relation"`
	AttributeGroupIns []*AttributeGroupIns `json:"attributeGroupIns" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (obj *Resource) AddAttributeGroupIns(target *AttributeGroupIns) {

	if target == nil {
		return
	}

	if obj.AttributeGroupIns == nil {
		obj.AttributeGroupIns = make([]*AttributeGroupIns, 0)
	}

	for _, attributeGroupIns := range obj.AttributeGroupIns {
		if attributeGroupIns.Uid == target.Uid {
			// 需要把属性保留下来
			for _, attributeIns := range target.AttributeIns {
				attributeGroupIns.AddAttributeIns(attributeIns)
			}
			return
		}
	}

	obj.AttributeGroupIns = append(obj.AttributeGroupIns, target)
}
