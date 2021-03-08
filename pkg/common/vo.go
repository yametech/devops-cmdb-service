package common

type ApiResponseVO struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

type ModelAttributeVisibleVO struct {
	Uid     string `json:"uid"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
}
