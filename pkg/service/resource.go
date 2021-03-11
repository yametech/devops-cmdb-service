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

func (rs *ResourceService) GetAllModeList() interface{} {
	modelList := &[]store.Model{}
	rs.Neo4jDomain.List(modelList)
	return modelList
}

// 获取模型实例列表
func (rs *ResourceService) GetResourcePageList(modelUid string, pageNumber int, pageSize int) interface{} {
	srcList := &[]store.Resource{}
	totalRaw, err := rs.ManualQueryRaw("MATCH (a:Resource {modelUid:$modelUid}) RETURN COUNT(a)",
		map[string]interface{}{"modelUid": modelUid})
	printOut(totalRaw[0][0])
	total := totalRaw[0][0].(int64)
	if err != nil {
		panic(err)
	}
	if total <= 0 {
		return common.PageResultVO{}
	}

	rs.ManualQuery("MATCH (a:Resource {modelUid:$modelUid}) RETURN a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit",
		map[string]interface{}{"modelUid": modelUid, "skip": (pageNumber - 1) * pageSize, "limit": pageSize}, srcList)

	printOut(srcList)

	pageResultVO := &common.PageResultVO{TotalCount: total}
	//list := make([]common.ResourcePageListVO, 0)
	list := make([]interface{}, 0)
	for _, srcResource := range *srcList {
		resource := &store.Resource{}
		store.GetSession(true).LoadDepth(resource, srcResource.UUID, 10)
		vo := &common.ResourcePageListVO{}
		utils.SimpleConvert(vo, resource)
		attributes := make(map[string]string)
		for _, srcAttributeGroupIns := range resource.AttributeGroupIns {
			for _, srcAttributeIns := range srcAttributeGroupIns.AttributeIns {
				attributes[srcAttributeIns.Uid] = srcAttributeIns.AttributeInsValue
			}
		}
		vo.Attributes = attributes
		list = append(list, vo)
	}
	pageResultVO.List = list
	return pageResultVO
}

func (rs *ResourceService) DeleteResource(uuid string) error {
	r := &store.Resource{}
	err := rs.Neo4jDomain.Get(r, "uuid", uuid)
	if err != nil {
		return err
	}

	query := "match (a:Resource)-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns) where a.uuid = $uuid detach delete a,b,c"
	_, err = rs.ManualExecute(query, map[string]interface{}{"uuid": uuid})
	return err
}

func (rs *ResourceService) AddResource(body string, operator string) (interface{}, error) {
	bodyObj := &store.Resource{}
	err := json.Unmarshal([]byte(body), bodyObj)
	if err != nil {
		return nil, err
	}
	//printOut(bodyObj)

	model := &store.Model{Uid: bodyObj.ModelUid}
	err = rs.Neo4jDomain.Get(model, "uid", bodyObj.ModelUid)
	if err != nil {
		return nil, err
	}

	// 获取模型详细
	//fullModel := fakeGetFullModel(rs)
	fullModel := &store.Model{}
	_ = store.GetSession(true).LoadDepth(fullModel, model.UUID, 2)

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)

	resource := &store.Resource{ModelUid: bodyObj.ModelUid, ModelName: bodyObj.ModelName, CommonObj: *commonObj}
	resource.Models = fullModel

	for _, groupObj := range bodyObj.AttributeGroupIns {
		attributeGroup := fullModel.GetAttributeGroupByUid(groupObj.Uid)
		if attributeGroup != nil {
			attributeGroupIns := &store.AttributeGroupIns{Uid: attributeGroup.Uid, Name: attributeGroup.Name}
			for _, attributeObj := range groupObj.AttributeIns {
				attribute := attributeGroup.GetAttributeByUid(attributeObj.Uid)
				if attribute != nil {
					attribute.AttributeCommon.Visible = true
					attributeIns := &store.AttributeIns{
						AttributeCommon:   attribute.AttributeCommon,
						AttributeInsValue: attributeObj.AttributeInsValue,
						CommonObj:         *commonObj,
					}
					attributeGroupIns.AddAttributeIns(attributeIns)
					resource.AddAttributeGroupIns(attributeGroupIns)
				}
			}
		}
	}

	err = store.GetSession(false).SaveDepth(resource, 10)
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

func fakeGetFullModel(rs *ResourceService) *store.Model {
	jsonStr := "{\"id\":144,\"uuid\":\"fdec6cdf-72e8-4966-a23b-5d4990574094\",\"uid\":\"host\",\"name\":\"主机\",\"iconUrl\":\"\",\"attributeGroups\":[{\"id\":143,\"uuid\":\"83e3088e-dd38-4469-a5c8-e703a3863e32\",\"uid\":\"baseInfo\",\"name\":\"基本属性\",\"modelUid\":\"host\",\"attributes\":[{\"id\":142,\"uuid\":\"e066f7f1-d932-4994-8ab8-03332d61279c\",\"uid\":\"ip\",\"name\":\"网址\",\"valueType\":\"短字符串\",\"editable\":true,\"required\":false,\"defaultValue\":\"\",\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":\"\",\"listValues\":\"\",\"tips\":\"\",\"regular\":\"(([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\\\\.){3}([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\",\"comment\":\"网址信息\",\"visible\":false,\"modelUid\":\"host\",\"creator\":\"\",\"editor\":\"\",\"createTime\":0,\"updateTime\":0}],\"creator\":\"\",\"editor\":\"\",\"createTime\":0,\"updateTime\":0},{\"id\":133,\"uuid\":\"bdb28935-6409-4344-88c5-b5b7a7bd117a\",\"uid\":\"otherInfo\",\"name\":\"其他属性\",\"modelUid\":\"host\",\"attributes\":[{\"id\":145,\"uuid\":\"67db51cc-180a-4c62-9926-b126a3961f00\",\"uid\":\"test\",\"name\":\"cesi\",\"valueType\":\"短字符串\",\"editable\":true,\"required\":false,\"defaultValue\":\"\",\"unit\":\"\",\"maximum\":\"\",\"minimum\":\"\",\"enums\":\"\",\"listValues\":\"\",\"tips\":\"\",\"regular\":\"(([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\\\\.){3}([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\",\"comment\":\"网址信息\",\"visible\":false,\"modelUid\":\"host\",\"creator\":\"\",\"editor\":\"\",\"createTime\":0,\"updateTime\":0}],\"creator\":\"\",\"editor\":\"\",\"createTime\":0,\"updateTime\":0}],\"resources\":null,\"creator\":\"\",\"editor\":\"\",\"createTime\":0,\"updateTime\":0}"
	a := &store.Model{}

	json.Unmarshal([]byte(jsonStr), a)
	//printOut(a)
	return a
}

func printOut(obj interface{}) {
	b, _ := json.Marshal(obj)
	fmt.Printf("%#v\n", string(b))
}
