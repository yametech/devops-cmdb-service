package common

type ApiResponseVO struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

//type Neo4jNode struct {
//	Id     int64                  `json:"id"`
//	Labels []string               `json:"labels"`
//	Props  map[string]interface{} `json:"props"`
//}

type ModelAttributeVisibleVO struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
}

type SimpleModelVO struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

type ConfigModelAttributeVO struct {
	ModelUid string                     `json:"modelUid"`
	Columns  *[]ModelAttributeVisibleVO `json:"columns"`
}
