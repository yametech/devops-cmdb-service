package store

import "github.com/yametech/devops-cmdb-service/pkg/gogm"

type AttributeIns struct {
	gogm.BaseNode
	AttributeCommon
	AttributeGroupIns *AttributeGroupIns `json:"-" gogm:"direction=outgoing;relationship=GroupBy"`
	AttributeInsValue string             `json:"attributeInsValue" gogm:"name=attributeInsValue"`
	CommonObj
}
