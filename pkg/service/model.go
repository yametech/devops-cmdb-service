package service

import (
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"strconv"
	"sync"
	"time"
)

type ModelService struct {
	Service
	store.Neo4jDomain
	mutex sync.Mutex
}

func (ms *ModelService) GetAllModelGroup() (*[]store.ModelGroup, error) {
	modelGroups := make([]store.ModelGroup, 0)
	err := store.GetSession(true).LoadAllDepth(&modelGroups, 2)
	return &modelGroups, err
}

func (ms *ModelService) GetModelGroup(uuid string) (*store.ModelGroup, error) {
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", uuid); err != nil {
		return nil, err
	}
	return modelGroup, nil
}

func (ms *ModelService) CreateModelGroup(vo *common.AddModelGroupVO, operator string) (*store.ModelGroup, error) {
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uid", vo.Uid); err == nil {
		return nil, fmt.Errorf("已存在uid是%s的模型分组", vo.Uid)
	}

	utils.SimpleConvert(modelGroup, vo)

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	modelGroup.CommonObj = *commonObj
	err := ms.Neo4jDomain.Save(modelGroup)
	return modelGroup, err
}

func (ms *ModelService) UpdateModelGroup(vo *common.AddModelGroupVO, operator string) (*store.ModelGroup, error) {
	modelGroup := &store.ModelGroup{}

	if err := ms.Neo4jDomain.Get(modelGroup, "uid", vo.Uid); err != nil {
		return nil, err
	}

	modelGroup.Name = vo.Name
	modelGroup.UpdateTime = time.Now().Unix()
	modelGroup.CommonObj.Editor = operator

	err := ms.Neo4jDomain.Update(modelGroup)
	return modelGroup, err
}

func (ms *ModelService) CreateModel(vo *common.AddModelVO, operator string) (*store.Model, error) {
	model := &store.Model{}
	if err := ms.Neo4jDomain.Get(model, "uid", vo.Uid); err == nil {
		return nil, fmt.Errorf("已存在uid是%s的模型", vo.Uid)
	}

	utils.SimpleConvert(model, vo)

	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", vo.ModelGroupUUID); err != nil {
		return nil, err
	}

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	model.CommonObj = *commonObj
	model.ModelGroup = modelGroup
	err := ms.Neo4jDomain.Save(model)
	return model, err
}

func (ms *ModelService) DeleteModel(uuid, operator string) error {
	model := &store.Model{}
	err := ms.Neo4jDomain.Get(model, "uuid", uuid)
	if err != nil {
		return err
	}

	resource := &store.Resource{}
	err = ms.Neo4jDomain.Get(resource, "modelUid", model.Uid)
	if resource.UUID != "" {
		return errors.New("该模型已被使用，禁止删除")
	}

	// TODO test
	return ms.Neo4jDomain.Delete(model)
}

func (ms *ModelService) GetModelList(limit string, pageNumber string) (*[]store.Model, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		return nil, err
	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return nil, err
	}
	allModel := make([]store.Model, 0)
	query := fmt.Sprintf("match (a:Model) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumberInt - 1) * limitInt,
		"limit": limitInt,
	}
	if err := store.GetSession(true).Query(query, properties, &allModel); err != nil {
		return nil, err
	}
	return &allModel, nil
}

func (ms *ModelService) GetModel(uuid string) (*store.Model, error) {
	model := &store.Model{}
	err := ms.Neo4jDomain.Get(model, "uuid", uuid)
	if err != nil {
		return nil, err
	}

	attributeGroups := &[]*store.AttributeGroup{}
	query := "MATCH (a:Model {uuid: $uuid})-[]-(b:AttributeGroup) RETURN b"
	_ = ms.ManualQuery(query, map[string]interface{}{"uuid": uuid}, attributeGroups)

	model.AttributeGroups = *attributeGroups

	for _, ag := range *attributeGroups {
		attributes := &[]*store.Attribute{}
		query = "MATCH (a:AttributeGroup {uuid: $uuid})-[]-(b:Attribute) RETURN b"
		_ = ms.ManualQuery(query, map[string]interface{}{"uuid": ag.UUID}, attributes)
		ag.Attributes = *attributes
	}
	return model, err
}

func (ms *ModelService) GetRelationshipList(pageSize, pageNumber int) (*[]store.RelationshipModel, error) {
	allModel := make([]store.RelationshipModel, 0)
	query := fmt.Sprintf("match (a:RelationshipModel) return a ORDER BY a.createTime DESC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumber - 1) * pageSize,
		"limit": pageSize,
	}
	if err := store.GetSession(true).Query(query, properties, &allModel); err != nil {
		return nil, err
	}
	// get all ModelRelation, count the Relationship usage
	relationService := RelationService{}
	relations := relationService.GetAllModelRelations()
	for _, model := range allModel {
		for _, relation := range *relations {
			if relation.RelationshipUid == model.Uid {
				model.Usage += 1
			}
		}
	}

	return &allModel, nil
}

func (ms *ModelService) GetRelationship(uuid string) (*store.RelationshipModel, error) {
	relationship := &store.RelationshipModel{}
	if err := ms.Neo4jDomain.Get(relationship, "uuid", uuid); err != nil {
		return nil, err
	}
	return relationship, nil
}

func (ms *ModelService) SaveRelationship(relation *store.RelationshipModel) error {

	err := ms.Neo4jDomain.Save(relation)
	if err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) UpdateRelationship(vo *common.RelationshipModelUpdateVO, operator string) error {
	relation := &store.RelationshipModel{}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	err := ms.Neo4jDomain.Get(relation, "uid", vo.Uid)
	if err != nil {
		return err
	}

	relation.Name = vo.Name
	relation.Source2Target = vo.Source2Target
	relation.Target2Source = vo.Target2Source
	relation.Direction = vo.Direction
	relation.UpdateTime = time.Now().Unix()
	relation.Editor = operator
	return ms.Neo4jDomain.Update(relation)
}

func (ms *ModelService) DeleteRelationship(uuid string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	relationshipModel := &store.RelationshipModel{}
	err := ms.Neo4jDomain.Get(relationshipModel, "uuid", uuid)
	if err != nil {
		return err
	}
	// if the Relationship has been used, deny operation
	relationService := RelationService{}
	relations := relationService.GetAllModelRelations()
	for _, relation := range *relations {
		if relation.RelationshipUid == relationshipModel.Uid {
			return errors.New("该关系模型已被使用，禁止删除")
		}
	}
	return ms.Neo4jDomain.Delete(relationshipModel)
}
