package store

import (
	"fmt"
	"github.com/mindstand/gogm"
	"time"
)

type CommonObj struct {
	//Id int64 `json:"id"`
	Creator    string `json:"creator" gogm:"name=creator"`
	Editor     string `json:"editor" gogm:"name=editor"`
	CreateTime int64  `json:"createTime" gogm:"name=createTime"`
	UpdateTime int64  `json:"updateTime" gogm:"name=updateTime"`
}

type Model struct {
	gogm.BaseNode
	Uid             string            `json:"uid" gogm:"unique;name=uid"`
	Name            string            `json:"name" gogm:"name=name"`
	IconUrl         string            `json:"iconUrl" gogm:"name=iconUrl"`
	ModelGroup      *ModelGroup       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	AttributeGroups []*AttributeGroup `json:"-" gogm:"direction=incoming;relationship=GroupBy"`
	Resources       []*Resource       `json:"-" gogm:"direction=incoming;relationship=Instance"`
	CommonObj
}

type AttributeGroup struct {
	gogm.BaseNode
	Uid        string       `json:"uid" gogm:"unique;name=uid"`
	Name       string       `json:"name" gogm:"name=name"`
	ModelUid   string       `json:"modelUid" gogm:"index;name=modelUid"`
	Model      *Model       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	Attributes []*Attribute `json:"-" gogm:"direction=incoming;relationship=GroupBy"`
	CommonObj
}

func (m *Model) Get(uid string) error {
	query := fmt.Sprintf("match (a:Model) where a.uid = $uid return a")
	properties := map[string]interface{}{
		"uid": uid,
	}
	return GetSession(false).Query(query, properties, m)
}


func (m *Model) LoadAll(mList *[]*Model, groupId string) error {
	query := fmt.Sprintf("match (a:Model)-[r:GroupBy]->(b:ModelGroup)where b.uid=$uid return a")
	properties := map[string]interface{}{
		"uid": groupId,
	}
	err := GetSession(true).Query(query, properties, mList)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) Save() error {
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = time.Now().Unix()
	return GetSession(false).Save(m)
}

func (m *Model) Update() error {
	m.UpdateTime = time.Now().Unix()
	return GetSession(false).Save(m)
}

func (m *Model) GetAttributeGroupByUid(uid string) *AttributeGroup {
	for _, group := range m.AttributeGroups {
		if group.Uid == uid {
			return group
		}
	}

	return nil
}

func (m *AttributeGroup) GetAttributeByUid(uid string) *Attribute {
	for _, attributes := range m.Attributes {
		if attributes.Uid == uid {
			return attributes
		}
	}

	return nil
}
