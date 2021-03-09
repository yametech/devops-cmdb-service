package service

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"time"
)

type ResourceService struct {
	Service
}

// 模型属性字段列表
func (rs *ResourceService) GetModelAttributeList(modelUid string) interface{} {
	a := &[]store.Attribute{}
	rs.ManualQuery("match (a:Attribute {modelUid:$modelUid}) return a", map[string]interface{}{"modelUid": modelUid}, a)
	fmt.Printf("%#v", a)
	return a
}

// 设置预览属性
func (rs *ResourceService) SetModelAttribute(modelUid string, result *[]common.ModelAttributeVisibleVO) error {
	for _, vo := range *result {
		_, _ = rs.ManualExecute("match (a:Attribute {modelUid:$modelUid, uid:$uid}) SET a.visible = $visible ,a.updateTime = $updateTime",
			map[string]interface{}{"modelUid": modelUid, "uid": vo.Uid, "visible": vo.Visible, "updateTime": time.Now().Unix()})
	}

	return nil
}
