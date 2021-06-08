package common

type K8sResource struct {
	CompassName  string `json:"compassName"`
	ResourceName string `json:"resourceName"`
	// 每个资源模型的属性参考http://wiki.ym/pages/viewpage.action?pageId=75893663
	ResourceAttribute map[string]string `json:"resourceAttribute"`
	ResourceRelation  map[string]string `json:"resourceRelation"`
}

type Relation struct {
	Name string
	Type string
}
