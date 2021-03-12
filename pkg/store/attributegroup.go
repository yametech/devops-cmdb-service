package store

import (
	"fmt"
	"github.com/mindstand/gogm"
	"time"
)

type AttributeGroup struct {
	gogm.BaseNode
	Uid        string       `json:"uid" gogm:"name=uid"`
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

func (m *AttributeGroup) Get(uuid string) error {
	query := fmt.Sprintf("match (a:AttributeGroup) where a.uuid = $uuid return a")
	properties := map[string]interface{}{
		"uuid": uuid,
	}
	return GetSession(false).Query(query, properties, m)
}

func (m *AttributeGroup) Save() error {
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = time.Now().Unix()
	return GetSession(false).Save(m)
}

func (m *AttributeGroup) Update() error {
	m.UpdateTime = time.Now().Unix()
	return GetSession(false).Save(m)
}

func (m *AttributeGroup) Delete() error {
	return GetSession(false).Delete(m)
}
