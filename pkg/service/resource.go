package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ResourceService struct {
	Service
}

// 创建实例时候获取模型基础信息
func (rs *ResourceService) GetModelInfoForIns(uid string) (*store.Resource, error) {
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
			map[string]interface{}{"modelUid": modelUid, "uid": vo.Uid, "visible": vo.Visible, "updateTime": time.Now().UnixNano() / 1000000})
	}

	return nil
}

func normalize(target string) string {
	var str = "~`!@#$%^&*()_-+={}[]:;\"'|\\,.<>/?"
	newRune := make([]rune, 0)
	chars := []rune(strings.TrimSpace(target))
	slic := []rune("\\")
	for _, c := range chars {
		if strings.Index(str, string(c)) > 0 {
			newRune = append(newRune, slic[0])
		}
		newRune = append(newRune, c)
	}
	return string(newRune)
}

// 获取模型实例列表，
// 支持2种查询方式：1.根据指定字段查询，2.不指定字段查询
func (rs *ResourceService) GetResourceListPage(queryVO *common.ResourceListPageParamVO) interface{} {
	// 1、分页信息
	queryCommon := "MATCH p=(a:Resource)-[r1]-(c:AttributeGroupIns)-[r2]-(b:AttributeIns) "
	where := "WHERE a.modelUid ='" + queryVO.ModelUid + "' AND "
	if len(queryVO.QueryTags) > 0 {
		for k, v := range queryVO.QueryTags {
			where += "(b.uid='" + k + "' AND b.attributeInsValue in ["
			for _, value := range v {
				where += "'" + value + "',"
			}
			where = strings.TrimSuffix(where, ",") + "]) OR "
		}
		where = strings.TrimSuffix(strings.TrimSpace(where), "OR")
	} else {
		for k, v := range *queryVO.QueryMap {
			if k == "ID" {
				where += "ID(a) =" + strings.TrimSpace(v) + " AND "
			} else {
				where += "b.uid='" + k + "' AND b.attributeInsValue=~'(?i).*" + normalize(v) + ".*' AND "
			}
		}
	}

	// uuid 不为空，结果需要进行过滤
	if queryVO.UUID != "" {
		query := "MATCH (a:Resource)-[r:Relation]-(b:Resource) WHERE a.uuid = $uuid and r.uid = $modelRelationUid RETURN b.uuid"
		uuidArray, err := rs.ManualQueryRaw(query, map[string]interface{}{"uuid": queryVO.UUID, "modelRelationUid": queryVO.ModelRelationUid})
		if err == nil && uuidArray != nil && len(uuidArray) > 0 {
			if queryVO.HasRelation == 0 {
				where += " NONE(x IN nodes(p) WHERE x.uuid in ["
				for _, v := range uuidArray {
					where += "'" + v[0].(string) + "',"
				}
				where = strings.TrimSuffix(where, ",") + "])"
			} else {
				where += " a.uuid in ["
				for _, v := range uuidArray {
					where += "'" + v[0].(string) + "',"
				}
				where = strings.TrimSuffix(where, ",") + "]"
			}
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

	rs.ManualQuery(query+"RETURN DISTINCT a ORDER BY id(a) DESC SKIP $skip LIMIT $limit ", //ORDER BY a.createTime ASC SKIP $skip LIMIT $limit",
		map[string]interface{}{"modelUid": queryVO.ModelUid, "skip": (queryVO.Current - 1) * queryVO.PageSize, "limit": queryVO.PageSize}, srcList)

	list := make([]interface{}, 0)
	// 2、批量查询资源完整信息
	// MATCH (a:Resource)-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns) WHERE a.uuid in ['',''] return a,b,c
	//resourceMap := make(map[string]*store.Resource)
	query = "MATCH (a:Resource)-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns) WHERE a.uuid in ["
	for _, srcResource := range *srcList {
		query += "'" + srcResource.UUID + "',"
	}
	query = strings.TrimSuffix(query, ",") + "]  RETURN a,b,c " //ORDER BY a.createTime ASC"
	//fmt.Println(query)
	result, err := rs.ManualQueryRaw(query, nil)
	if err != nil {
		panic(err)
	}

	ldapUserMap := map[string]common.LdapUserVO{}
	resources := utils.GetResourceFromNeo4jResult(result)
	for _, resource := range resources {
		attributes := make(map[string]interface{})
		for _, srcAttributeGroupIns := range resource.AttributeGroupIns {
			for _, srcAttributeIns := range srcAttributeGroupIns.AttributeIns {
				if srcAttributeIns.ValueType == "用户" {
					ldapService := &LdapService{}
					if len(ldapUserMap) == 0 {
						ldapUserMap = ldapService.GetLdapUserMap()
					}
					if ldapUserMap[srcAttributeIns.AttributeInsValue].Name != "" {
						srcAttributeIns.AttributeInsValue = ldapUserMap[srcAttributeIns.AttributeInsValue].Name
					}
				}
				attributes[srcAttributeIns.Uid] = srcAttributeIns.AttributeInsValue
			}
		}
		attributes["id"] = resource.Id
		attributes["uuid"] = resource.UUID
		list = append(list, attributes)
	}
	pageResultVO := &common.PageResultVO{TotalCount: total}
	// 根据id倒排序
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].(map[string]interface{})["id"].(int64) > list[j].(map[string]interface{})["id"].(int64)
	})
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

func (rs *ResourceService) AddResource(body string, operator string) (*store.Resource, error) {
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
	modelService := ModelService{}
	fullModel, err := modelService.GetModel(model.UUID)
	if err != nil && fullModel != nil {
		if fullModel.UUID == "" {
			return nil, errors.New("该模型已被删除")
		}
		return nil, err
	}

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)

	resource := &store.Resource{ModelUid: model.Uid, ModelName: model.Name, CommonObj: *commonObj}
	resource.Models = model

	attributeInsValueMap := map[string]string{}
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
					attributeInsValueMap[attributeGroupIns.Uid+attributeIns.Uid] = attributeIns.AttributeInsValue
				}
			}
		}
	}

	// 验证数据
	if len(resource.AttributeGroupIns) == 0 {
		return nil, errors.New("禁止插入不完整资源实例数据")
	}
	// 唯一性校验
	result, err := rs.saveResourceWithUniqueCheck(resource, fullModel, attributeInsValueMap)
	if err != nil {
		return result, err
	}
	// 阿里账号同步
	if resource.ModelUid == "aliyun_account" {
		syncService := &SyncService{}
		go syncService.SyncAliDomainByResource(resource, operator)
	}

	return resource, nil
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
		//match, err := regexp.MatchString("^-?[0-9]*$", insValue)
		value, err := strconv.Atoi(insValue)
		if err != nil {
			return fmt.Errorf("字段%q内容%q不符合%q规范", attribute.Name, insValue, "数字")
		}
		if attribute.Maximum != "" {
			maximum, err := strconv.Atoi(attribute.Maximum)
			if err == nil && value > maximum {
				return fmt.Errorf("字段%q内容%q大于最大值%q", attribute.Name, insValue, maximum)
			}
		}
		if attribute.Minimum != "" {
			minimum, err := strconv.Atoi(attribute.Minimum)
			if err == nil && value < minimum {
				return fmt.Errorf("字段%q内容%q小于最小值%q", attribute.Name, insValue, minimum)
			}
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

func (rs *ResourceService) UpdateResource(body string, operator string) (interface{}, error) {
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
	resourceTemplate, err := modelService.GetModelDetail(bodyObj.ModelUid)
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
				sourceAttIns.UpdateTime = time.Now().UnixNano() / 1000000
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
	source.UpdateTime = time.Now().UnixNano() / 1000000

	return rs.saveResourceWithUniqueCheck(source, resourceTemplate, attributeInsValueMap)
}

func (rs *ResourceService) saveResourceWithUniqueCheck(resource *store.Resource, resourceTemplate *store.Model, attributeInsValueMap map[string]string) (*store.Resource, error) {
	session := rs.GetSession(false)
	defer session.Close()
	// 属性唯一性验证
	var attributeMaps = getUniqueAttributeMaps(resourceTemplate, attributeInsValueMap)
	if len(attributeMaps) > 0 {
		ec := utils.EtcdClient{}
		//TODO 粗糙实现
		mutex := ec.NewMutex(resourceTemplate.Uid)
		if err := mutex.Lock(context.TODO()); err != nil {
			log.Println("get etcd mutex failed " + err.Error())
			return nil, err
		}
		defer mutex.Unlock(context.TODO())
		// 查询重复记录
		for _, attributeMap := range attributeMaps {
			query := "MATCH (a:Resource)<-[]-(b:AttributeGroupIns)<-[]-(c:AttributeIns) " +
				"WHERE a.modelUid=$modelUid and b.uid=$groupUid and c.uid=$uid and c.attributeInsValue=$insValue RETURN a.uuid"
			result, err := rs.ManualQueryRaw(query, attributeMap)
			if err != nil {
				return nil, err
			}
			if result != nil && len(result[0]) > 0 {
				if len(result[0]) > 1 || result[0][0].(string) != resource.UUID {
					//b, _ := json.Marshal(attributeMap)
					return nil, fmt.Errorf("数据唯一性校验失败，%q重复", attributeMap["name"])
				}
			}
		}
	}

	err := session.SaveDepth(resource, 2)
	return resource, err
}

func getUniqueAttributeMaps(model *store.Model, attributeInsValueMap map[string]string) []map[string]interface{} {
	var attributeMaps = make([]map[string]interface{}, 0)
	for _, group := range model.AttributeGroups {
		for _, att := range group.Attributes {
			if att.Unique && attributeInsValueMap[group.Uid+att.Uid] != "" {
				attributeMap := map[string]interface{}{}
				attributeMap["uid"] = att.Uid
				attributeMap["name"] = att.Name
				attributeMap["insValue"] = strings.TrimSpace(attributeInsValueMap[group.Uid+att.Uid])
				attributeMap["groupUid"] = group.Uid
				attributeMap["modelUid"] = model.Uid
				attributeMaps = append(attributeMaps, attributeMap)
			}
		}
	}
	return attributeMaps
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

	query := "MATCH (a:Resource)<-[]-(b:AttributeGroupIns)<-[]-(c:AttributeIns) WHERE a.uuid=$uuid RETURN *  ORDER BY b.createTime, c.createTime ASC"
	result, err := rs.ManualQueryRaw(query, map[string]interface{}{"uuid": uuid})
	if err != nil {
		return nil, err
	}

	modelService := ModelService{}
	resourceTemplate, err := modelService.GetModelDetail(r.ModelUid)
	if err != nil || resourceTemplate == nil {
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
		for _, group := range resourceTemplate.AttributeGroups {
			if group.Uid == attributeGroupIns.Uid {
				attributeGroupIns.Name = group.Name
				for _, attribute := range group.Attributes {
					if attribute.Uid == attributeIns.Uid {
						attributeIns.AttributeCommon = attribute.AttributeCommon
						break
					}
				}
				break
			}
		}
		r.AddAttributeGroupIns(attributeGroupIns)
	}
	r.ModelName = resourceTemplate.Name
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
	a.UpdateTime = time.Now().UnixNano() / 1000000
	a.Editor = editor
	return rs.Save(a)
}
