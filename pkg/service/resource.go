package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ResourceService struct {
	Service
}

// 创建实例时候获取模型基础信息
func (rs ResourceService) GetModelInfoForIns(uid string) (*store.Resource, error) {
	m := &store.Model{}
	err := rs.Neo4jDomain.Get(m, "uid", uid)
	if err != nil {
		if m.UUID == "" {
			return nil, errors.New("模型已被删除")
		}
		return nil, err
	}

	resource := &store.Resource{}
	utils.SimpleConvert(resource, m)
	resource.ModelUid = m.Uid
	resource.ModelName = m.Name

	query := "MATCH (a:Model)-[]-(b:AttributeGroup)-[]-(c:Attribute) WHERE a.uid =$uid RETURN a,b,c ORDER BY c.createTime ASC"
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
	list := make([]*store.Attribute, 0)

	resource, err := rs.GetModelInfoForIns(modelUid)
	if err != nil {
		return nil
	}
	for _, group := range resource.AttributeGroupIns {
		for _, att := range group.AttributeIns {
			attribute := &store.Attribute{}
			utils.SimpleConvert(attribute, att)
			list = append(list, attribute)
		}
	}

	return list
}

// 设置预览属性
func (rs *ResourceService) SetModelAttribute(modelUid string, result *[]common.ModelAttributeVisibleVO) error {
	for _, vo := range *result {
		_, _ = rs.ManualExecute("MATCH (a:Attribute {modelUid:$modelUid, uid:$uid}) SET a.visible = $visible ,a.updateTime = $updateTime",
			map[string]interface{}{"modelUid": modelUid, "uid": vo.Uid, "visible": vo.Visible, "updateTime": time.Now().Unix()})
	}

	return nil
}

// 获取模型实例列表，
// 支持2种查询方式：1.根据指定字段查询，2.不指定字段查询
func (rs *ResourceService) GetResourceListPageByMap(uuid, modelUid, modelRelationUid string, pageNumber int, pageSize int, queryMap *map[string]string) interface{} {
	queryCommon := "MATCH p=(a:Resource)-[]-()-[]-(b:AttributeIns) "
	where := "WHERE a.modelUid ='" + modelUid + "' AND "
	for k, v := range *queryMap {
		if k == "ID" {
			where += "ID(a) =" + strings.TrimSpace(v) + " AND "
		} else {
			where += "b.uid='" + k + "' AND b.attributeInsValue=~'.*" + strings.TrimSpace(v) + ".*' AND "
		}
	}

	// uuid 不为空，则需排除跟此资源有关联的实例
	if uuid != "" {
		query := "MATCH (a:Resource)-[r:Relation]-(b:Resource) WHERE a.uuid = $uuid and r.uid = $modelRelationUid RETURN b.uuid"
		uuidArray, err := rs.ManualQueryRaw(query, map[string]interface{}{"uuid": uuid, "modelRelationUid": modelRelationUid})
		if err == nil && uuidArray != nil && len(uuidArray) > 0 {
			where += " NONE(x IN nodes(p) WHERE x.uuid in ["

			for _, v := range uuidArray {
				where += "'" + v[0].(string) + "',"
			}
			where = strings.TrimSuffix(where, ",") + "])"
		}
	}
	query := queryCommon + strings.TrimSuffix(strings.TrimSpace(where), "AND") + " "
	fmt.Println(query)
	srcList := &[]store.Resource{}
	totalRaw, err := rs.ManualQueryRaw(query+"RETURN COUNT(distinct a)", nil)
	if err != nil {
		panic(err)
	}
	total := totalRaw[0][0].(int64)
	if total <= 0 {
		return &common.PageResultVO{List: []interface{}{}}
	}

	rs.ManualQuery(query+"RETURN DISTINCT a ORDER BY a.createTime ASC SKIP $skip LIMIT $limit",
		map[string]interface{}{"modelUid": modelUid, "skip": (pageNumber - 1) * pageSize, "limit": pageSize}, srcList)

	pageResultVO := &common.PageResultVO{TotalCount: total}
	list := make([]interface{}, 0)
	for _, srcResource := range *srcList {
		resource := &store.Resource{}
		err = rs.GetSession(true).LoadDepth(resource, srcResource.UUID, 2)
		if err != nil {
			panic(err)
		}
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
func (rs *ResourceService) GetResourceListPageByQueryValue(modelUid, queryValue string, pageNumber int, pageSize int) interface{} {
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

	model := &store.Model{Uid: bodyObj.ModelUid}
	err = rs.Neo4jDomain.Get(model, "uid", bodyObj.ModelUid)
	if err != nil {
		if model.UUID == "" {
			return nil, errors.New("该模型已被删除")
		}
		return nil, err
	}

	// 获取模型详细
	fullModel := &store.Model{}
	err = rs.GetSession(true).LoadDepth(fullModel, model.UUID, 2)
	if err != nil {
		if fullModel.UUID == "" {
			return nil, errors.New("该模型已被删除")
		}
		return nil, err
	}

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)

	resource := &store.Resource{ModelUid: bodyObj.ModelUid, ModelName: bodyObj.ModelName, CommonObj: *commonObj}
	resource.Models = model

	for _, groupObj := range bodyObj.AttributeGroupIns {
		attributeGroup := fullModel.GetAttributeGroupByUid(groupObj.Uid)
		if attributeGroup != nil {
			attributeGroupIns := &store.AttributeGroupIns{Uid: attributeGroup.Uid, Name: attributeGroup.Name}
			for _, attributeObj := range groupObj.AttributeIns {
				attribute := attributeGroup.GetAttributeByUid(attributeObj.Uid)
				if attribute != nil {
					err := rs.attributeInsValueValidate(attribute, attributeObj.AttributeInsValue)
					if err != nil {
						return nil, err
					}
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

	err = rs.GetSession(false).SaveDepth(resource, 2)
	return resource, err
}

func (rs *ResourceService) attributeInsValueValidate(attribute *store.Attribute, attributeInsValue string) error {
	insValue := strings.TrimSpace(attributeInsValue)
	// 必填
	if attribute.Required && insValue == "" {
		return fmt.Errorf("字段%q必填", attribute.Name)
	}

	if insValue == "" {
		return nil
	}

	// 正则规范
	if attribute.Regular != "" {
		pattern, err := strconv.Unquote(`"` + attribute.Regular + `"`)
		if err != nil {
			fmt.Printf("Unquote%q,err%q\n ", attribute.Regular, err)
			//return err
		}
		match, err := regexp.MatchString(pattern, insValue)
		fmt.Println(attribute.Regular, pattern, insValue, match, err)
		if err != nil {
			return err
		}
		if !match {
			return fmt.Errorf("字段%q内容%q不符合正则规范%q", attribute.Name, insValue, attribute.Regular)
		}
	}

	// 类型:短字符,长字符,数字,浮点数,枚举,日期,时间,用户,布尔,列表
	// 短字符：长度 256 个英文字符或 85 个中文字符
	// 长字符：长度 2000 英文或 666 个中文字符
	// 数字：正负整数
	// 浮点数：可以包含小数点的数字
	// 枚举：包含 K-V 结构的列表
	// 日期：日期格式
	// 时间：时间格式
	// 用户：可以搜索【授权中心 - 用户管理】中已经录入的用户
	// 布尔：布尔类型，常用于开关
	// 列表：可以理解为数组类型，只包含值的列表
	switch {
	case attribute.ValueType == "短字符" && len(insValue) > 256:
		return fmt.Errorf("字段%q内容%q不符合规范，短字符：长度 256 个英文字符或 85 个中文字符", attribute.Name, insValue)
	case attribute.ValueType == "长字符" && len(insValue) > 2000:
		return fmt.Errorf("字段%q内容不符合规范，长字符：长度 2000 英文或 666 个中文字符", attribute.Name)
	case attribute.ValueType == "数字":
		match, err := regexp.MatchString("^-?[0-9]*$", insValue)
		if err != nil {
			return err
		}
		if !match {
			return fmt.Errorf("字段%q内容%q不符合%q规范", attribute.Name, insValue, "数字")
		}
	case attribute.ValueType == "浮点数":
		match, err := regexp.MatchString("^(-?\\d+)(\\.\\d+)?$", insValue)
		if err != nil {
			return err
		}
		if !match {
			return fmt.Errorf("字段%q内容%q不符合%q规范", attribute.Name, insValue, "浮点数")
		}
	case attribute.ValueType == "枚举":
		if attribute.Enums != nil {
			type enum struct {
				Id    interface{} `json:"id"`
				Value interface{} `json:"value"`
			}
			enums := make([]enum, 0)
			err := json.Unmarshal([]byte(attribute.Enums.(string)), &enums)
			if err != nil {
				return err
			}

			exit := false
			for _, item := range enums {
				if item.Id == insValue {
					exit = true
					break
				}
			}
			if !exit {
				return fmt.Errorf("字段%q内容%q不在枚举值里面", attribute.Name, insValue)
			}
		}
	case attribute.ValueType == "列表":
		if attribute.ListValues != nil {
			type list struct {
				Value interface{} `json:"value"`
			}
			lists := make([]list, 0)
			err := json.Unmarshal([]byte(attribute.ListValues.(string)), &lists)
			if err != nil {
				return err
			}

			exit := false
			for _, item := range lists {
				if item.Value == insValue {
					exit = true
					break
				}
			}
			if !exit {
				return fmt.Errorf("字段%q内容%q不在列表值里面", attribute.Name, insValue)
			}
		}
	case attribute.ValueType == "布尔" && (insValue != "true" && insValue != "false"):
		return fmt.Errorf("字段%q内容%q不符合%q规范", attribute.Name, insValue, "布尔")
	case attribute.ValueType == "用户":
	case attribute.ValueType == "日期":
		_, err := time.Parse("2006-01-02", insValue)
		if err != nil {
			return fmt.Errorf("字段%q内容%q不符合%q规范", attribute.Name, insValue, "日期")
		}
	case attribute.ValueType == "时间":
		_, err := time.Parse("2006-01-02 15:04:05", insValue)
		if err != nil {
			return fmt.Errorf("字段%q内容%q不符合%q规范", attribute.Name, insValue, "时间")
		}
	}
	return nil
}

func (rs *ResourceService) UpdateResource(body string, operator string) (*store.Resource, error) {
	fmt.Println(body)
	bodyObj := &store.Resource{}
	err := json.Unmarshal([]byte(body), bodyObj)
	if err != nil {
		return nil, err
	}

	source, err := rs.GetResourceDetail(bodyObj.UUID)
	if err != nil {
		return nil, err
	}

	// map: groupInsUid+attInsUid = attributeInsValue
	attributeInsValueMap := map[string]string{}
	for _, groupIns := range bodyObj.AttributeGroupIns {
		for _, attIns := range groupIns.AttributeIns {
			attributeInsValueMap[groupIns.Uid+attIns.Uid] = attIns.AttributeInsValue
		}
	}

	// validate and update
	// 获取创建模板
	modelService := ModelService{}
	resourceTemplate, err := modelService.GetModelDetail(source.ModelUid)
	if err != nil {
		return nil, err
	}
	// validate
	for _, group := range resourceTemplate.AttributeGroups {
		for _, att := range group.Attributes {
			v, ok := attributeInsValueMap[group.Uid+att.Uid]
			if ok {
				err := rs.attributeInsValueValidate(att, v)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// update
	for _, sourceGroupIns := range source.AttributeGroupIns {
		for _, sourceAttIns := range sourceGroupIns.AttributeIns {
			v, ok := attributeInsValueMap[sourceGroupIns.Uid+sourceAttIns.Uid]
			if ok {
				sourceAttIns.AttributeInsValue = v
				sourceAttIns.UpdateTime = time.Now().Unix()
				sourceAttIns.Editor = operator
			}
		}
	}
	// 如果有新增的属性
	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	for _, attributeGroupObj := range bodyObj.AttributeGroupIns {
		exits := false
		for _, sourceGroupIns := range source.AttributeGroupIns {
			if sourceGroupIns.Uid == attributeGroupObj.Uid {
				// 继续检测属性
				for _, attributeObj := range attributeGroupObj.AttributeIns {
					attributeObjExits := false
					for _, sourceAttributeIns := range sourceGroupIns.AttributeIns {
						if sourceAttributeIns.Uid == attributeObj.Uid {
							attributeObjExits = true
							break
						}
					}
					if !attributeObjExits {
						newAttributeIns := getNewAttributeFromTemplate(resourceTemplate, attributeObj.Uid)
						newAttributeIns.CommonObj = *commonObj
						newAttributeIns.AttributeInsValue = attributeObj.AttributeInsValue
						sourceGroupIns.AddAttributeIns(newAttributeIns)
					}
				}
				exits = true
				break
			}
		}
		// 新的分组
		if !exits {
			newGroupIns := &store.AttributeGroupIns{Uid: attributeGroupObj.Uid, Name: attributeGroupObj.Name}
			for _, attributeObj := range attributeGroupObj.AttributeIns {
				newAttributeIns := getNewAttributeFromTemplate(resourceTemplate, attributeObj.Uid)
				newAttributeIns.CommonObj = *commonObj
				newAttributeIns.AttributeInsValue = attributeObj.AttributeInsValue
				newGroupIns.AddAttributeIns(newAttributeIns)
			}
			source.AddAttributeGroupIns(newGroupIns)
		}
	}

	source.Editor = operator
	source.UpdateTime = time.Now().Unix()
	err = rs.GetSession(false).SaveDepth(source, 2)
	return source, err
}

func getNewAttributeFromTemplate(temp *store.Model, uid string) *store.AttributeIns {
	for _, groupIns := range temp.AttributeGroups {
		for _, attribute := range groupIns.Attributes {
			if attribute.Uid == uid {
				newIns := &store.AttributeIns{}
				utils.SimpleConvert(newIns, attribute)
				newIns.AttributeGroupIns = nil
				newIns.BaseNode = gogm.BaseNode{}
				return newIns
			}
		}
	}
	return nil
}

// 获取资源详情
func (rs *ResourceService) GetResourceDetail(uuid string) (*store.Resource, error) {
	r := &store.Resource{}
	err := rs.Neo4jDomain.Get(r, "uuid", uuid)
	if err != nil {
		if r.UUID == "" {
			return nil, fmt.Errorf("资源已被删除")
		}
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
			return errors.New("内容不符合正则规范:" + a.Regular)
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
