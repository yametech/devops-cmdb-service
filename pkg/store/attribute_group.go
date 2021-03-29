package store

import (
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
)

type AttributeGroup struct {
	gogm.BaseNode
	Uid        string       `json:"uid" gogm:"index;name=uid"`
	Name       string       `json:"name" gogm:"name=name"`
	ModelUid   string       `json:"modelUid" gogm:"index;name=modelUid"`
	Model      *Model       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	Attributes []*Attribute `json:"attributes" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (m *AttributeGroup) AddAttribute(target *Attribute) {
	if target == nil {
		return
	}

	if m.Attributes == nil {
		m.Attributes = make([]*Attribute, 0)
	}

	for _, attribute := range m.Attributes {
		if attribute.Uid == target.Uid {
			return
		}
	}

	m.Attributes = append(m.Attributes, target)
}
