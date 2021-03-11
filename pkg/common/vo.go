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

type ResourcePageListVO struct {
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
