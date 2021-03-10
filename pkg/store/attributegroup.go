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
	Attributes []*Attribute `json:"-" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (m *AttributeGroup) Get(uid string) error {
	query := fmt.Sprintf("match (a:AttributeGroup) where a.uid = $uid return a")
	properties := map[string]interface{}{
		"uid": uid,
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
