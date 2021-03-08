package store

import (
	"github.com/mindstand/gogm"
)

type CommonObj struct {
	//Id int64 `json:"id"`
	Creator    string `json:"creator" gogm:"name=creator"`
	Editor     string `json:"editor" gogm:"name=editor"`
	CreateTime string `json:"createTime" gogm:"name=createTime"`
	UpdateTime string `json:"updateTime" gogm:"name=updateTime"`
}

type Model struct {
	gogm.BaseNode
	Uid             string            `json:"uid" gogm:"unique;name=uid"`
	Name            string            `json:"name" gogm:"name=name"`
	IconUrl         string            `json:"iconUrl" gogm:"name=iconUrl"`
	ModelGroup      *ModelGroup       `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	AttributeGroups []*AttributeGroup `json:"-" gogm:"direction=incoming;relationship=GroupBy"`
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
