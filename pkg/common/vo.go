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
	ModelUid   string             `form:"model_uid" json:"modelUid" binding:"required"`
	QueryValue string             `json:"queryValue" binding:""`
	QueryMap   *map[string]string `json:"queryMap" binding:""`
}

type AddModelVO struct {
	ModelGroupUUID string `json:"modelGroupUUID" binding:"required"`
	Uid            string `json:"uid" binding:"required"`
	Name           string `json:"name" binding:"required"`
}
