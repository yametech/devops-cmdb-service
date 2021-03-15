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
	"strconv"
	"strings"
	"time"
)

type ResourceService struct {
	Service
	store.Neo4jDomain
}

// 创建实例时候获取模型基础信息
func (rs ResourceService) GetModelInfoForIns(uid string) (interface{}, error) {
	m := &store.Model{}
	err := rs.Neo4jDomain.Get(m, "uid", uid)
	if err != nil {
		return nil, err
	}

	resource := &store.Resource{}
	utils.SimpleConvert(resource, m)
	resource.ModelUid = m.Uid
	resource.ModelName = m.Name

	query := "MATCH (a:Model)-[]-(b:AttributeGroup)-[]-(c:Attribute) WHERE a.uid =$uid RETURN a,b,c"
	result, err := rs.ManualQueryRaw(query, map[string]interface{}{"uid": uid})
	if err != nil {
		return nil, err
	}

	for _, row := range result {
		// 属性
		o := row[2].(*gogm.NodeWrap)
		attribute := &store.AttributeIns{}
		utils.SimpleConvert(attribute, &o.Props)

		// 属性分组
		o = row[1].(*gogm.NodeWrap)
		attributeGroup := &store.AttributeGroupIns{}
		utils.SimpleConvert(attributeGroup, &o.Props)

		attributeGroup.AddAttributeIns(attribute)
		resource.AddAttributeGroupIns(attributeGroup)
	}
	return resource, nil
}

// 模型属性字段列表
func (rs *ResourceService) GetModelAttributeList(modelUid string) interface{} {
	a := &[]store.Attribute{}
	rs.Neo4jDomain.Get(a, "modelUid", modelUid)
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

// 获取模型实例列表，
// 支持2种查询方式：1.根据指定字段查询，2.不指定字段查询
func (rs ResourceService) GetResourceListPageByMap(modelUid string, pageNumber int, pageSize int, queryMap *map[string]string) interface{} {
	queryCommon := "MATCH (a:Resource)-[]-()-[]-(b:AttributeIns) "
	where := "WHERE a.modelUid ='" + modelUid + "' AND "
	for k, v := range *queryMap {
		where += "b." + k + "='" + v + "' AND "
	}
	query := queryCommon + strings.TrimSuffix(strings.TrimSpace(where), "AND") + " "
	fmt.Println(query)
	srcList := &[]store.Resource{}
	totalRaw, err := rs.ManualQueryRaw(query+"RETURN COUNT(distinct a)", nil)
	if err != nil {
		panic(err)
	}
	printOut(totalRaw[0][0])
	total := totalRaw[0][0].(int64)
	if total <= 0 {
		return common.PageResultVO{}
	}

	rs.ManualQuery(query+"RETURN a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit",
		map[string]interface{}{"modelUid": modelUid, "skip": (pageNumber - 1) * pageSize, "limit": pageSize}, srcList)

	printOut(srcList)

	pageResultVO := &common.PageResultVO{TotalCount: total}
	list := make([]interface{}, 0)
	for _, srcResource := range *srcList {
		resource := &store.Resource{}
		err = store.GetSession(true).LoadDepth(resource, srcResource.UUID, 2)
		if err != nil {
			panic(err)
		}
		//vo := &common.ResourceListPageVO{}
		//utils.SimpleConvert(vo, resource)
		attributes := make(map[string]string)
		for _, srcAttributeGroupIns := range resource.AttributeGroupIns {
			for _, srcAttributeIns := range srcAttributeGroupIns.AttributeIns {
				attributes[srcAttributeIns.Uid] = srcAttributeIns.AttributeInsValue
			}
		}
		//vo.Attributes = attributes
		attributes["id"] = strconv.FormatInt(resource.Id, 10)
		attributes["uuid"] = resource.UUID
		list = append(list, attributes)
	}
	pageResultVO.List = list
	return pageResultVO
}

// 不指定字段查询
func (rs *ResourceService) GetResourceListPage(modelUid, queryValue string, pageNumber int, pageSize int) interface{} {
	// TODO
	return nil
}

func (rs *ResourceService) DeleteResource(uuidArray []string) error {
	for _, uuid := range uuidArray {
		r := &store.Resource{}
		err := rs.Neo4jDomain.Get(r, "uuid", uuid)
		if err != nil {
			return err
		}

		query := "match (a:Resource)-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns) where a.uuid = $uuid detach delete a,b,c"
		_, err = rs.ManualExecute(query, map[string]interface{}{"uuid": uuid})
		if err != nil {
			return err
		}
	}

	return nil
}

func (rs *ResourceService) AddResource(body string, operator string) (interface{}, error) {
	fmt.Println(body)
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

func printOut(obj interface{}) {
	b, _ := json.Marshal(obj)
	fmt.Printf("%#v\n", string(b))
}
