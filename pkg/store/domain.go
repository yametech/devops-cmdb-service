package store

import (
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
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
	Uid           string `json:"uid" gogm:"unique;name=uid"`
	Name          string `json:"name" gogm:"name=name"`
	Source2Target string `json:"source2Target" gogm:"name=source2Target"` // 源->目标描述
	Target2Source string `json:"target2Source" gogm:"name=target2Source"` // 目标->源描述
	Direction     string `json:"direction" gogm:"name=direction"`         // direction 方向：源指向目标，无方向，双方向
	CurrentUsage  int    `json:"currentUsage" gogm:"name=currentUsage"`
	CommonObj
}

// 模型关系
type ModelRelation struct {
	gogm.BaseNode
	Uid             string      `json:"uid"`
	RelationshipUid string      `json:"relationshipUid"`
	Constraint      string      `json:"constraint"`
	SourceUid       string      `json:"sourceUid"`
	TargetUid       string      `json:"targetUid"`
	Comment         interface{} `json:"comment"`
	CommonObj
}

type Model struct {
	gogm.BaseNode
	Uid             string            `json:"uid" gogm:"unique;name=uid"`
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
