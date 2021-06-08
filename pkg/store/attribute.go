package store

import (
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
)

type AttributeCommon struct {
	Uid          string      `json:"uid" gogm:"index;name=uid"`             //  唯一标识
	Name         string      `json:"name" gogm:"name=name"`                 // 名称
	ValueType    string      `json:"valueType" gogm:"name=valueType"`       // 类型:短字符,长字符,数字,浮点,枚举,日期,时间,用户,布尔,列表
	Editable     bool        `json:"editable" gogm:"name=editable"`         // 是否可编辑
	Required     bool        `json:"required" gogm:"name=required"`         // 是否必填
	Unique       bool        `json:"unique" gogm:"name=unique"`             // 是否唯一
	DefaultValue interface{} `json:"defaultValue" gogm:"name=defaultValue"` // 默认值
	Unit         string      `json:"unit" gogm:"name=unit"`                 // 单位
	Maximum      string      `json:"maximum" gogm:"name=maximum"`           // 最大值
	Minimum      string      `json:"minimum" gogm:"name=minimum"`           // 最小值
	Enums        interface{} `json:"enums" gogm:"name=enums"`               // 枚举值：{id1:value1,id2:value2...}
	ListValues   interface{} `json:"listValues" gogm:"name=listValues"`     // 列表：value1,value2
	Tips         string      `json:"tips" gogm:"name=tips"`                 // 用户提示内容
	Regular      string      `json:"regular" gogm:"name=regular"`           // 正则内容
	Comment      string      `json:"comment" gogm:"name=comment"`           // 备注描述
	Visible      bool        `json:"visible" gogm:"name=visible"`           // 字段是否可见
	ModelUid     string      `json:"modelUid" gogm:"index;name=modelUid"`   // 模型唯一标识
}

type Attribute struct {
	gogm.BaseNode
	AttributeCommon
	AttributeGroup *AttributeGroup `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	// 公共字段：创建人，更新人，创建时间，更新时间
	CommonObj
}
