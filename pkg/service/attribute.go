package service

import (
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type AttributeService struct {
	Service
}

func (as *AttributeService) GetAttributeGroupList(limit string, pageNumber string) (*[]store.AttributeGroup, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}

	allAG := make([]store.AttributeGroup, 0)
	query := fmt.Sprintf("match (a:AttributeGroup) return a ORDER BY a.createTime ASC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}

	if err := as.ManualQuery(query, properties, &allAG); err != nil {
		return nil, err
	}
	return &allAG, nil
}

func (as *AttributeService) GetAttributeGroup(uuid string) (*store.AttributeGroup, error) {
	attributeGroup := &store.AttributeGroup{}
	if err := as.GetSession(true).Load(attributeGroup, uuid); err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) CreateAttributeGroup(attributeGroupVO *common.AddAttributeGroupVO, operator string) (*store.AttributeGroup, error) {
	attributeGroupVO.ModelUUID = strings.TrimSpace(attributeGroupVO.ModelUUID)
	attributeGroupVO.Uid = strings.TrimSpace(attributeGroupVO.Uid)
	attributeGroupVO.Name = strings.TrimSpace(attributeGroupVO.Name)
	// validate
	if err := UidNameValidate(attributeGroupVO.Uid, attributeGroupVO.Name); err != nil {
		return nil, err
	}

	as.mutex.Lock()
	defer as.mutex.Unlock()
	model := &store.Model{}
	if err := as.Neo4jDomain.Get(model, "uuid", attributeGroupVO.ModelUUID); err != nil {
		if model.UUID == "" {
			return nil, fmt.Errorf("该模型已被删除")
		}
		return nil, err
	}

	// check exist
	query := "MATCH (a:Model {uuid: $modelUUID})<-[]-(b:AttributeGroup {uid: $uid}) RETURN COUNT(distinct b)"
	totalRaw, err := as.ManualQueryRaw(query, map[string]interface{}{"modelUUID": attributeGroupVO.ModelUUID, "uid": attributeGroupVO.Uid})
	if err != nil {
		return nil, err
	}
	var total int64
	if totalRaw != nil && len(totalRaw) > 0 && len(totalRaw[0]) > 0 {
		total = totalRaw[0][0].(int64)
	}
	if total > 0 {
		return nil, errors.New("该模型下已存在此分组唯一标识:" + attributeGroupVO.Uid)
	}

	query = "MATCH (a:Model {uuid: $modelUUID})<-[]-(b:AttributeGroup {name: $name}) RETURN COUNT(distinct b)"
	totalRaw, err = as.ManualQueryRaw(query, map[string]interface{}{"modelUUID": attributeGroupVO.ModelUUID, "name": attributeGroupVO.Name})
	if err != nil {
		return nil, err
	}

	if totalRaw != nil && len(totalRaw) > 0 && len(totalRaw[0]) > 0 {
		total = totalRaw[0][0].(int64)
	}
	if total > 0 {
		return nil, errors.New("该模型下已存在此分组名称:" + attributeGroupVO.Name)
	}

	attributeGroup := &store.AttributeGroup{}
	utils.SimpleConvert(attributeGroup, attributeGroupVO)
	attributeGroup.ModelUid = model.Uid
	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	attributeGroup.CommonObj = *commonObj
	attributeGroup.Model = model

	if err := as.Neo4jDomain.Save(attributeGroup); err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) isAttributeGroupExist(modelUid, uid, name string) bool {
	exists := &[]store.AttributeGroup{}
	query := "MATCH (a:AttributeGroup {modelUid: $modelUid, uid: $uid}) RETURN a"
	properties := map[string]interface{}{
		"modelUid": modelUid,
		"uid":      uid,
		"name":     name,
	}
	_ = as.ManualQuery(query, properties, exists)

	return len(*exists) > 0
}

func (as *AttributeService) UpdateAttributeGroup(attributeGroupVO *common.UpdateAttributeGroupVO, operator string) (*store.AttributeGroup, error) {
	attributeGroupVO.Name = strings.TrimSpace(attributeGroupVO.Name)

	attributeGroup := &store.AttributeGroup{}
	if err := as.Neo4jDomain.Get(attributeGroup, "uuid", attributeGroupVO.UUID); err != nil {
		return nil, err
	}

	attributeGroup.Name = attributeGroupVO.Name
	attributeGroup.UpdateTime = time.Now().Unix()
	attributeGroup.Editor = operator
	if err := as.Neo4jDomain.Update(attributeGroup); err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) DeleteAttributeGroup(uuid string) error {
	attributeGroup := &store.AttributeGroup{}
	if err := as.Neo4jDomain.Get(attributeGroup, "uuid", uuid); err != nil {
		if attributeGroup.UUID == "" {
			return errors.New("属性分组已被删除")
		}
		return err
	}

	as.GetSession(true).LoadDepth(attributeGroup, uuid, 1)
	if len(attributeGroup.Attributes) > 0 {
		return errors.New("此分组下存在属性，不可删除")
	}

	if err := as.GetSession(false).DeleteUUID(uuid); err != nil {
		return err
	}
	return nil
}

func (as *AttributeService) CreateAttribute(vo *common.CreateAttributeVO, operator string) (*store.Attribute, error) {
	vo.Uid = strings.TrimSpace(vo.Uid)
	vo.Name = strings.TrimSpace(vo.Name)
	// validate
	if err := UidNameValidate(vo.Uid, vo.Name); err != nil {
		return nil, err
	}
	if strings.Index(common.ATTRIBUTE_TYPE, vo.ValueType) < 0 {
		return nil, fmt.Errorf("%q不在属性类型范围内：%q", vo.ValueType, common.ATTRIBUTE_TYPE)
	}

	as.mutex.Lock()
	defer as.mutex.Unlock()
	group, _ := as.GetAttributeGroup(vo.AttributeGroupUUID)
	if group == nil {
		return nil, errors.New("数据异常，属性分组不存在，uuid:" + vo.AttributeGroupUUID)
	}

	if group.ModelUid != vo.ModelUId {
		return nil, errors.New("数据异常，modelUid不匹配，modelUid:" + group.ModelUid)
	}

	// check exist
	query := "MATCH (a:AttributeGroup {uuid: $attributeGroupUUID})<-[]-(b:Attribute {uid: $uid}) RETURN COUNT(distinct b)"
	totalRaw, err := as.ManualQueryRaw(query, map[string]interface{}{"attributeGroupUUID": group.UUID, "uid": vo.Uid})
	if err != nil {
		return nil, err
	}
	total := totalRaw[0][0].(int64)
	if total > 0 {
		return nil, errors.New("该分组下已存在此属性唯一标识:" + vo.Uid)
	}

	query = "MATCH (a:AttributeGroup {uuid: $attributeGroupUUID})<-[]-(b:Attribute {name: $name}) RETURN COUNT(distinct b)"
	totalRaw, err = as.ManualQueryRaw(query, map[string]interface{}{"attributeGroupUUID": group.UUID, "name": vo.Name})
	if err != nil {
		return nil, err
	}
	total = totalRaw[0][0].(int64)
	if total > 0 {
		return nil, errors.New("该分组下已存在此属性名称:" + vo.Name)
	}

	attribute := &store.Attribute{}
	utils.SimpleConvert(attribute, vo)
	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	attribute.CommonObj = *commonObj
	attribute.AttributeGroup = group

	if err := as.Neo4jDomain.Save(attribute); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) GetAttributeList(limit string, pageNumber string) (*[]store.Attribute, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}

	attributeList := make([]store.Attribute, 0)
	query := fmt.Sprintf("match (a:Attribute) return a ORDER BY a.createTime ASC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := as.GetSession(true).Query(query, properties, &attributeList); err != nil {
		return nil, err
	}
	return &attributeList, nil
}

func (as *AttributeService) GetAttribute(uuid string) (*store.Attribute, error) {
	attribute := &store.Attribute{}
	if err := as.Get(attribute, "uuid", uuid); err != nil {
		if attribute.UUID == "" {
			return nil, errors.New("属性已被删除")
		}
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) UpdateAttribute(vo *common.UpdateAttributeVO, operator string) (*store.Attribute, error) {
	source, err := as.GetAttribute(vo.UUID)
	if err != nil {
		return nil, errors.New("属性不存在")
	}
	attribute := &store.Attribute{}
	utils.SimpleConvert(attribute, vo)
	attribute.UpdateTime = time.Now().Unix()
	attribute.Editor = operator

	// 固定不变的值
	attribute.Uid = source.Uid
	attribute.ModelUid = source.ModelUid
	attribute.CreateTime = source.CreateTime
	attribute.Creator = source.Creator

	if err := as.Neo4jDomain.Update(attribute); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) DeleteAttributeInstance(uuid string) error {
	attribute, err := as.GetAttribute(uuid)
	if err != nil {
		return err
	}
	if err := as.Neo4jDomain.Delete(attribute); err != nil {
		return err
	}
	return nil
}
