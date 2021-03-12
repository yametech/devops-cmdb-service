package store

import (
	"fmt"
	"github.com/mindstand/gogm"
	"strings"
	"time"
)

type CommonObj struct {
	//Id int64 `json:"id"`
	Creator    string `json:"creator" gogm:"name=creator"`
	Editor     string `json:"editor" gogm:"name=editor"`
	CreateTime int64  `json:"createTime" gogm:"name=createTime"`
	UpdateTime int64  `json:"updateTime" gogm:"name=updateTime"`
}

// 关系模型
type RelationshipModel struct {
	gogm.BaseNode
	Uid  string `json:"uid" gogm:"unique;name=uid"`
	Name string `json:"name" gogm:"name=name"`
	// 源->目标描述
	Source2Target string `json:"source2Target" gogm:"name=source2Target"`
	// 目标->源描述
	Target2Source string `json:"target2Source" gogm:"name=target2Source"`
	// direction 方向：源指向目标，无方向，双方向
	Direction string `json:"direction" gogm:"name=direction"`
	CommonObj
}

//type ModelRelationTest struct {
//	gogm.BaseNode
//	Start      *Model
//	End        *Model
//	Constraint string `json:"constraint" gogm:"name=constraint"`
//	// 源uid
//	SourceUid string `json:"sourceUid" gogm:"name=sourceUid"`
//	// 目标uid
//	TargetUid string `json:"targetUid" gogm:"name=targetUid"`
//	// 关系类型uid
//	RelationshipUid string `json:"relationshipUid" gogm:"name=relationshipUid"`
//	CommonObj
//}

// 模型关系
type ModelRelation struct {
	//gogm.BaseNode
	//Uid             string `json:"uid" gogm:"name=uid"`
	//RelationshipUid string `json:"relationshipUid" gogm:"name=relationshipUid"`
	//Constraint      string `json:"constraint" gogm:"name=constraint"`
	//SourceUid       string `json:"sourceUid" gogm:"name=sourceUid"`
	//TargetUid       string `json:"targetUid" gogm:"name=targetUid"`
	//Comment         string `json:"comment" gogm:"name=comment"`
	Uid             string `json:"uid"`
	RelationshipUid string `json:"relationshipUid"`
	Constraint      string `json:"constraint"`
	SourceUid       string `json:"sourceUid"`
	TargetUid       string `json:"targetUid"`
	Comment         string `json:"comment"`
	CommonObj
}

type Model struct {
	gogm.BaseNode
	Uid string `json:"uid" gogm:"unique;name=uid"`
	//Uuid             string            `json:"Uuid" gogm:"unique;name=uuid"`
	Name            string            `json:"name" gogm:"name=name"`
	IconUrl         string            `json:"iconUrl" gogm:"name=iconUrl"`
	Model           *Model            `json:"-" gogm:"direction=both;relationship=Relation"`
	ModelGroup      *ModelGroup       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	AttributeGroups []*AttributeGroup `json:"attributeGroups" gogm:"direction=incoming;relationship=GroupBy"`
	Resources       []*Resource       `json:"resources" gogm:"direction=incoming;relationship=Instance"`
	CommonObj
}

func (cm *CommonObj) InitCommonObj(creator string) {
	cm.CreateTime = time.Now().Unix()
	cm.UpdateTime = time.Now().Unix()
	cm.Creator = creator
	cm.Editor = creator
}

func (m *Model) AddAttributeGroup(target *AttributeGroup) {
	if target == nil {
		return
	}

	if m.AttributeGroups == nil {
		m.AttributeGroups = make([]*AttributeGroup, 0)
	}

	for _, attributeGroup := range m.AttributeGroups {
		if attributeGroup.Uid == target.Uid {
			return
		}
	}

	m.AttributeGroups = append(m.AttributeGroups, target)
}

func (m *Model) Get(session *gogm.Session, uuid string) error {
	query := fmt.Sprintf("match (a:Model) where a.uuid = $uuid return a")
	properties := map[string]interface{}{
		"uuid": uuid,
	}

	return session.Query(query, properties, m)
}

func (m *Model) LoadAll(session *gogm.Session, groupId string) ([]*Model, error) {
	mList := make([]*Model, 0)
	query := fmt.Sprintf("match (a:Model)-[r:GroupBy]->(b:ModelGroup)where b.uuid=$uuid return a")
	properties := map[string]interface{}{
		"uuid": groupId,
	}
	err := session.Query(query, properties, &mList)

	if err != nil {
		if strings.Contains(err.Error(), "data not found") {
			return nil, nil
		}
		return nil, err
	}
	return mList, nil
}

func (m *Model) Save(session *gogm.Session) error {
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = time.Now().Unix()
	return session.Save(m)
}

func (m *Model) Update(session *gogm.Session) error {
	m.UpdateTime = time.Now().Unix()
	return session.Save(m)
}

func (m *Model) Delete(session *gogm.Session) error {
	return session.Delete(m)
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
