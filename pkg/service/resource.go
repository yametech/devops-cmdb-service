package service

import (
	"encoding/json"
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"time"
)

type ResourceService struct {
	Service
}

// 模型属性字段列表
func (rs *ResourceService) GetModelAttributeList(modelUid string) interface{} {
	a := &[]store.Attribute{}
	rs.ManualQuery("MATCH (a:Attribute {modelUid:$modelUid}) RETURN a", map[string]interface{}{"modelUid": modelUid}, a)
	fmt.Printf("%#v", a)
	return a
}

// 设置预览属性
func (rs *ResourceService) SetModelAttribute(modelUid string, result *[]common.ModelAttributeVisibleVO) error {
	for _, vo := range *result {
		_, _ = rs.ManualExecute("MATCH (a:Attribute {modelUid:$modelUid, uid:$uid}) SET a.visible = $visible ,a.updateTime = $updateTime",
			map[string]interface{}{"modelUid": modelUid, "uid": vo.Uid, "visible": vo.Visible, "updateTime": time.Now().Unix()})
	}

	return nil
}

// 获取模型实例列表
func (rs *ResourceService) GetResourceList(modelUid string, currentPage int, pageSize int) interface{} {
	r := &[]store.Resource{}
	rs.ManualQuery("MATCH (a:Resource {modelUid:$modelUid}) ORDER BY a.createTime DESC SKIP $skip LIMIT $limit",
		map[string]interface{}{"modelUid": modelUid, "skip": currentPage * pageSize, "limit": pageSize}, r)
	return r
}

// 获取资源详情
func (rs *ResourceService) GetResourceDetail(uuid string) (interface{}, error) {
	neo4j := store.Neo4jDomain{}
	r := &store.Resource{}
	err := neo4j.Get(r, "uuid", uuid)
	if err != nil {
		return nil, err
	}

	query := "MATCH (a:Resource)<-[]-(b:AttributeGroupIns)<-[]-(c:AttributeIns) WHERE a.uuid=$uuid RETURN *"
	result, err := rs.ManualQueryRaw(query, map[string]interface{}{"uuid": uuid})
	if err != nil {
		return nil, err
	}

	for _, row := range result {
		// 属性
		o := row[2].(*gogm.NodeWrap)
		attributeIns := &store.AttributeIns{}
		utils.SimpleConvert(attributeIns, &o.Props)

		// 属性分组
		o = row[1].(*gogm.NodeWrap)
		attributeGroupIns := &store.AttributeGroupIns{}
		utils.SimpleConvert(attributeGroupIns, &o.Props)

		attributeGroupIns.AddAttributeIns(attributeIns)
		r.AddAttributeGroupIns(attributeGroupIns)
	}
	return r, nil
}

func printOut(obj interface{}) {
	b, _ := json.Marshal(obj)
	fmt.Printf("%#v\n", string(b))
}
