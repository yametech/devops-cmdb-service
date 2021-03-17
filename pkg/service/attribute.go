package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"strconv"
	"sync"
	"time"
)

type AttributeService struct {
	Service
	store.Neo4jDomain
	Mutex sync.Mutex
}

func (as *AttributeService) ChangeModelGroup(attribute *store.Attribute, uuid string) error {
	attributeGroup, err := as.GetAttributeGroupInstance(uuid)
	if err != nil {
		return err
	}
	as.Mutex.Lock()
	defer as.Mutex.Unlock()

	query := fmt.Sprintf("match (a:Attribute)-[r:GroupBy]->(b:AttributeGroup)where a.uuid=$uuid delete r")
	properties := map[string]interface{}{
		"uuid": attribute.UUID,
	}

	as.ManualExecute(query, properties)

	attribute.AttributeGroup = attributeGroup
	if err := as.Save(attribute); err != nil {
		return err
	}
	return nil
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
	query := fmt.Sprintf("match (a:AttributeGroup) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}

	if err := as.ManualQuery(query, properties, &allAG); err != nil {
		return nil, err
	}
	/*for i, v := range allAG {
		attributes, err := as.Attribute.LoadAll(as.Session, v.UUID)
		if err != nil {
			continue
		}
		allAG[i].Attributes = *attributes
	}*/
	return &allAG, nil
}

func (as *AttributeService) GetAttributeGroupInstance(uuid string) (*store.AttributeGroup, error) {
	attributeGroup := &store.AttributeGroup{}
	if err := store.GetSession(true).Load(attributeGroup, uuid); err != nil {
		return nil, err
	}
	return attributeGroup, nil
}

func (as *AttributeService) CreateAttributeGroup(attributeGroupVO *common.AddAttributeGroupVO, operator string) (*store.AttributeGroup, error) {
	model := &store.Model{}
	if err := as.Neo4jDomain.Get(model, "uuid", attributeGroupVO.ModelUUID); err != nil {
		return nil, err
	}

	// TODO check exist

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

	//attributeGroup := &store.AttributeGroup{}
	//if err := as.Neo4jDomain.Get(attributeGroup, "uuid", uuid); err != nil {
	//	return err
	//}
	//if err := as.Neo4jDomain.Delete(attributeGroup); err != nil {
	//	return err
	//}
	return nil
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
	query := fmt.Sprintf("match (a:Attribute) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := store.GetSession(true).Query(query, properties, &attributeList); err != nil {
		return nil, err
	}
	return &attributeList, nil
}

func (as *AttributeService) GetAttribute(uuid string) (*store.Attribute, error) {

	attribute := &store.Attribute{}
	if err := as.Neo4jDomain.Get(attribute, "uuid", uuid); err != nil {
		return nil, err
	}
	return attribute, nil
}

func (as *AttributeService) UpdateAttribute(rawData []byte, uuid string) (*store.Attribute, error) {
	_, err := as.GetAttribute(uuid)
	if err != nil {
		return nil, errors.New("attribute not exists")
	}
	attribute := &store.Attribute{}
	if err := json.Unmarshal(rawData, attribute); err != nil {
		return nil, err
	}
	attribute.UUID = uuid

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
