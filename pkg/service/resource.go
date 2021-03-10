package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"regexp"
	"time"
)

type ResourceService struct {
	Service
	store.Neo4jDomain
}

// 模型属性字段列表
func (rs *ResourceService) GetModelAttributeList(modelUid string) interface{} {
	a := &[]store.Attribute{}
	rs.Neo4jDomain.Get(a, "modelUid", modelUid)
	//rs.ManualQuery("MATCH (a:Attribute {modelUid:$modelUid}) RETURN a", map[string]interface{}{"modelUid": modelUid}, a)
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

func (rs *ResourceService) DeleteResource(uuid string) error {
	r := &store.Resource{}
	err := rs.Neo4jDomain.Get(r, "uuid", uuid)
	if err != nil {
		return err
	}

	return rs.Neo4jDomain.Delete(r)
}

func (rs *ResourceService) AddResource(body string, operator string) (interface{}, error) {
	bodyObj := &store.Resource{}
	err := json.Unmarshal([]byte(body), bodyObj)
	if err != nil {
		return nil, err
	}

	model := store.Model{Uid: bodyObj.ModelUid}
	err = rs.Neo4jDomain.Get(model, "uid", bodyObj.ModelUid)
	if err != nil {
		return nil, err
	}

	// 获取模型详细
	fullModel := fakeGetFullModel()

	commonObj := initCommonObj(operator)
	resource := &store.Resource{ModelUid: bodyObj.ModelUid, ModelName: bodyObj.ModelName, CommonObj: *commonObj}
	resource.Models = &model

	for _, groupObj := range bodyObj.AttributeGroupIns {
		attributeGroup := fullModel.GetAttributeGroupByUid(groupObj.Uid)
		if attributeGroup != nil {
			attributeGroupIns := &store.AttributeGroupIns{Uid: attributeGroup.Uid, Name: attributeGroup.Name}
			resource.AddAttributeGroupIns(attributeGroupIns)
			for _, attributeObj := range groupObj.AttributeIns {
				attribute := attributeGroup.GetAttributeByUid(attributeObj.Uid)
				attributeIns := &store.AttributeIns{
					AttributeCommon:   attribute.AttributeCommon,
					AttributeInsValue: attributeObj.AttributeInsValue,
					CommonObj:         *commonObj,
				}
				attributeGroupIns.AddAttributeIns(attributeIns)
			}
		}
	}

	err = resource.Save()
	return resource, err
}

// 获取资源详情
func (rs *ResourceService) GetResourceDetail(uuid string) (interface{}, error) {
	r := &store.Resource{}
	err := rs.Neo4jDomain.Get(r, "uuid", uuid)
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

func (rs *ResourceService) UpdateResourceAttribute(uuid string, attributeInsValue string, editor string) error {
	a := &store.AttributeIns{}
	err := rs.Neo4jDomain.Get(a, "uuid", uuid)
	if err != nil {
		return err
	}
	if len(a.Regular) > 0 {
		match, _ := regexp.MatchString(a.Regular, attributeInsValue)
		if !match {
			return errors.New("内容不符合正则规范")
		}
	}

	a.AttributeInsValue = attributeInsValue
	a.UpdateTime = time.Now().Unix()
	a.Editor = editor
	return rs.Save(a)
}

func initCommonObj(creator string) *store.CommonObj {
	return &store.CommonObj{Creator: creator, Editor: creator, CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
}

func fakeGetFullModel() *store.Model {
	return &store.Model{}
}

func printOut(obj interface{}) {
	b, _ := json.Marshal(obj)
	fmt.Printf("%#v\n", string(b))
}
