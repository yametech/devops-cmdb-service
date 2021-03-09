package store

import (
	"fmt"
	"github.com/mindstand/gogm"
)

type AttributeCommon struct {
	gogm.BaseNode
	//  唯一标识
	Uid string `json:"uid" gogm:"unique;name=uid"`
	// 名称
	Name string `json:"name" gogm:"name=name"`
	// 类型:短字符,长字符,数字,浮点,枚举,日期,时间,用户,布尔,列表
	ValueType string `json:"valueType" gogm:"name=valueType"`
	// 是否可编辑
	Editable bool `json:"editable" gogm:"name=editable"`
	// 是否必填
	Required bool `json:"required" gogm:"name=required"`
	// 默认值
	DefaultValue string `json:"defaultValue" gogm:"name=defaultValue"`
	// 单位
	Unit string `json:"unit" gogm:"name=unit"`
	// 最大值
	Maximum string `json:"maximum" gogm:"name=maximum"`
	// 最小值
	Minimum string `json:"minimum" gogm:"name=minimum"`
	// 枚举值：{id1:value1,id2:value2...}
	Enums string `json:"enums" gogm:"name=enums"`
	// 列表：{value1,value2}
	ListValues string `json:"listValues" gogm:"name=listValues"`
	// 用户提示内容
	Tips string `json:"tips" gogm:"name=tips"`
	// 正则内容
	Regular string `json:"regular" gogm:"name=regular"`
	// 备注描述
	Comment string `json:"comment" gogm:"name=comment"`
	// 字段是否可见
	Visible bool `json:"visible" gogm:"name=visible"`
	// 模型唯一标识
	ModelUid string `json:"modelUid" gogm:"name=modelUid"`
	// 公共字段：创建人，更新人，创建时间，更新时间
	CommonObj
}

type Attribute struct {
	AttributeCommon
	AttributeGroup *AttributeGroup `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
}

func (obj *Attribute) Save() error {
	//obj.Visible = true
	return GetSession(false).Save(obj)
}

func (obj *Attribute) List(field string, value interface{}) interface{} {
	result := &[]Attribute{}
	query := fmt.Sprintf("MATCH (a:Attribute {%s:$%s})", field, field)
	properties := map[string]interface{}{field: value}
	_ = GetSession(true).Query(query, properties, result)
	return result
}
