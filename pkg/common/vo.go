package common

type ApiResponseVO struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

type PageResultVO struct {
	TotalCount int64         `json:"totalCount"`
	List       []interface{} `json:"list"`
}

type ModelAttributeVisibleVO struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
}

type ModelMenuVO struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

type ResourceListPageVO struct {
	Id         int64             `json:"id"`
	Uuid       string            `json:"uuid"`
	ModelUid   string            `json:"modelUid"`
	ModelName  string            `json:"modelName"`
	Attributes map[string]string `json:"attributes"`
}

type ConfigModelAttributeVO struct {
	Uid     string                     `json:"uid"`
	Columns *[]ModelAttributeVisibleVO `json:"columns"`
}

type ResourceRelationListPageVO struct {
	RelationshipUid  string                     `json:"relationshipUid"`
	RelationshipName string                     `json:"relationshipName"`
	SourceUid        string                     `json:"sourceUid"`
	SourceName       string                     `json:"sourceName"`
	TargetUid        string                     `json:"targetUid"`
	TargetName       string                     `json:"targetName"`
	ModelAttributes  *[]ModelAttributeVisibleVO `json:"modelAttributes"`
	Resources        *[]map[string]string       `json:"resources"`
}

type ResourceListPageParamVO struct {
	PageSize         int                `form:"pageSize" json:"pageSize" binding:"required,gte=0"`
	Current          int                `form:"current" json:"current" binding:"required,gte=0"`
	UUID             string             `form:"uuid" json:"uuid" binding:""`
	ModelRelationUid string             `form:"modelRelationUid" json:"modelRelationUid" binding:""`
	ModelUid         string             `form:"modelUid" json:"modelUid" binding:"required"`
	QueryValue       string             `json:"queryValue" binding:""`
	QueryMap         *map[string]string `json:"queryMap" binding:""`
}

type AddModelGroupVO struct {
	Uid  string `json:"uid" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type AddModelVO struct {
	ModelGroupUUID string `json:"modelGroupUUID" binding:"required"`
	Uid            string `json:"uid" binding:"required"`
	Name           string `json:"name" binding:"required"`
}

type UpdateModelVO struct {
	UUID string `json:"uuid" binding:"required"`
	Uid  string `json:"uid" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type AddAttributeGroupVO struct {
	ModelUUID string `json:"modelUUID" binding:"required"`
	Uid       string `json:"uid" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type UpdateAttributeGroupVO struct {
	UUID      string `json:"uuid" binding:"required"`
	ModelUUID string `json:"modelUUID" binding:"required"`
	Uid       string `json:"uid" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type CreateRelationshipModelVO struct {
	Uid           string `json:"uid" form:"uid" binding:"required"`
	Name          string `json:"name" form:"name" binding:"required"`
	Source2Target string `json:"source2Target" form:"source2Target" binding:"required,gte=1,lte=15"`
	Target2Source string `json:"target2Source" form:"target2Source" binding:"required,gte=1,lte=15"`
	Direction     string `json:"direction" form:"direction" binding:"required"`
}

type UpdateRelationshipModelVO struct {
	Uid           string `json:"uid" form:"uid" binding:"required"`
	Name          string `json:"name" form:"name" binding:"required"`
	Source2Target string `json:"source2Target" form:"source2Target" binding:"required,gte=1,lte=15"`
	Target2Source string `json:"target2Source" form:"target2Source" binding:"required,gte=1,lte=15"`
	Direction     string `json:"direction" form:"direction" binding:"required"`
}

type ModelRelationVO struct {
	Id               int64       `json:"id"`
	UUID             string      `json:"uuid"`
	Uid              string      `json:"uid"`
	RelationshipUid  string      `json:"relationshipUid"`
	RelationshipName string      `json:"relationshipName"`
	Constraint       string      `json:"constraint"`
	SourceUid        string      `json:"sourceUid"`
	SourceName       string      `json:"sourceName"`
	TargetUid        string      `json:"targetUid"`
	TargetName       string      `json:"targetName"`
	Comment          interface{} `json:"comment"`
}

type AddModelRelationVO struct {
	SourceUid       string      `json:"sourceUid" form:"sourceUid" binding:"required"`
	TargetUid       string      `json:"targetUid" form:"targetUid" binding:"required"`
	RelationshipUid string      `json:"relationshipUid" form:"relationshipUid" binding:"required"`
	Constraint      string      `json:"constraint" form:"constraint" binding:"required"`
	Comment         interface{} `json:"comment"`
}

type UpdateModelRelationVO struct {
	Uid             string      `json:"uid" form:"uid" binding:"required"`
	SourceUid       string      `json:"sourceUid" form:"sourceUid" binding:"required"`
	TargetUid       string      `json:"targetUid" form:"targetUid" binding:"required"`
	RelationshipUid string      `json:"relationshipUid" form:"relationshipUid" binding:"required"`
	Constraint      string      `json:"constraint" form:"constraint" binding:"required"`
	Comment         interface{} `json:"comment"`
}

type IdVO struct {
	Uid  string `json:"uid"`
	UUID string `json:"uuid"`
}

type ResourceRelationVO struct {
	SourceUUID string `json:"sourceUUID" form:"sourceUUID" binding:"required"`
	TargetUUID string `json:"targetUUID" form:"targetUUID" binding:"required"`
	Uid        string `json:"uid" form:"uid" binding:"required"`
}

type CreateAttributeVO struct {
	ModelUId           string      `json:"modelUId" form:"modelUId" binding:"required"`
	AttributeGroupUUID string      `json:"attributeGroupUUID" form:"attributeGroupUUID" binding:"required"`
	Uid                string      `json:"uid" form:"uid" binding:"required"`
	Name               string      `json:"name" form:"name" binding:"required"`
	ValueType          string      `json:"valueType" form:"valueType" binding:"required"`
	Editable           bool        `json:"editable" form:"editable" binding:""`
	Required           bool        `json:"required" form:"required" binding:""`
	Regular            string      `json:"regular" form:"regular" binding:""`
	Comment            string      `json:"comment" form:"comment" binding:""`
	DefaultValue       interface{} `json:"defaultValue" form:"defaultValue"`
	Unit               string      `json:"unit" form:"unit"`
	Maximum            string      `json:"maximum" form:"maximum"`
	Minimum            string      `json:"minimum" form:"minimum"`
	Enums              interface{} `json:"enums" form:"enums"`
	ListValues         interface{} `json:"listValues" form:"listValues"`
	Tips               string      `json:"tips" form:"tips"`
}

type UpdateAttributeVO struct {
	UUID         string      `json:"uuid" form:"uuid" binding:"required"`
	ModelUId     string      `json:"modelUId" form:"modelUId" binding:"required"`
	Uid          string      `json:"uid" form:"uid" binding:"required"`
	Name         string      `json:"name" form:"name" binding:"required"`
	ValueType    string      `json:"valueType" form:"valueType" binding:"required"`
	Editable     bool        `json:"editable" form:"editable" binding:""`
	Required     bool        `json:"required" form:"required" binding:""`
	Regular      string      `json:"regular" form:"regular" binding:""`
	Comment      string      `json:"comment" form:"comment" binding:""`
	DefaultValue interface{} `json:"defaultValue" form:"defaultValue"`
	Unit         string      `json:"unit" form:"unit"`
	Maximum      string      `json:"maximum" form:"maximum"`
	Minimum      string      `json:"minimum" form:"minimum"`
	Enums        interface{} `json:"enums" form:"enums"`
	ListValues   interface{} `json:"listValues" form:"listValues"`
	Tips         string      `json:"tips" form:"tips"`
}

type LdapUserVO struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}
