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

type UpdateModelRelationVO struct {
	Uid             string `json:"uid"`
	RelationshipUid string `json:"relationshipUid"`
	Constraint      string `json:"constraint"`
	SourceUid       string `json:"sourceUid"`
	TargetUid       string `json:"targetUid"`
	Comment         string `json:"comment"`
}

type ResourceListPageParamVO struct {
	PageSize   int                `form:"pageSize" json:"pageSize" binding:"required,gte=0"`
	Current    int                `form:"current" json:"current" binding:"required,gte=0"`
	ModelUid   string             `form:"modelUid" json:"modelUid" binding:"required"`
	QueryValue string             `json:"queryValue" binding:""`
	QueryMap   *map[string]string `json:"queryMap" binding:""`
}

type AddModelVO struct {
	ModelGroupUUID string `json:"modelGroupUUID" binding:"required"`
	Uid            string `json:"uid" binding:"required"`
	Name           string `json:"name" binding:"required"`
}

type RelationshipModelUpdateVO struct {
	Uid           string `json:"uid" form:"uid" binding:"required"`
	Name          string `json:"name" form:"name" binding:"required"`
	Source2Target string `json:"source2Target" form:"source2Target" binding:"required"`
	Target2Source string `json:"target2Source" form:"target2Source" binding:"required"`
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

type IdVO struct {
	Uid  string `json:"uid"`
	UUID string `json:"uuid"`
}
