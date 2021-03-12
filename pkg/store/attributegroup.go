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

func (m *AttributeGroup) Get(session *gogm.Session, uuid string) error {

	query := fmt.Sprintf("match (a:AttributeGroup) where a.uuid = $uuid return a")
	properties := map[string]interface{}{
		"uuid": uuid,
	}
	return session.Query(query, properties, m)
}

func (ag *AttributeGroup) LoadAll(session *gogm.Session, uuid string) error {
	query := fmt.Sprintf("match (a:AttributeGroup)-[r:GroupBy]->(b:Model) where b.uuid = $uuid return a")
	properties := map[string]interface{}{
		"uuid": uuid,
	}
	return session.Query(query, properties, ag)
}

func (ag *AttributeGroup) Save(session *gogm.Session) error {
	ag.CreateTime = time.Now().Unix()
	ag.UpdateTime = time.Now().Unix()
	return session.Save(ag)
}

func (ag *AttributeGroup) Update(session *gogm.Session) error {
	ag.UpdateTime = time.Now().Unix()
	return session.Save(ag)
}

func (ag *AttributeGroup) Delete(session *gogm.Session) error {
	return session.Delete(ag)
}
