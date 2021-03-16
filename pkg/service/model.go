package service

import (
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"strconv"
	"sync"
	"time"
)

type ModelService struct {
	Service
	store.Neo4jDomain
	mutex sync.Mutex
}

func (ms *ModelService) ChangeModelGroup(model *store.Model, uuid string) error {
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", uuid); err != nil {
		return err
	}

	query := fmt.Sprintf("match (a:Model)-[r:GroupBy]->(b:ModelGroup)where a.uuid=$uuid delete r")
	properties := map[string]interface{}{
		"uuid": model.UUID,
	}
	_ = store.GetSession(false).Query(query, properties, nil)
	model.ModelGroup = modelGroup
	if err := ms.Neo4jDomain.Update(model); err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) GetAllGroup() (*[]store.ModelGroup, error) {
	modelGroups := make([]store.ModelGroup, 0)
	err := ms.Neo4jDomain.List(&modelGroups)
	return &modelGroups, err
}

func (ms *ModelService) GetModelGroup(uuid string) (*store.ModelGroup, error) {
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", uuid); err != nil {
		return nil, err
	}
	return modelGroup, nil
}

func (ms *ModelService) CreateModelGroup(model *store.ModelGroup) error {
	return ms.Neo4jDomain.Save(model)
}

func (ms *ModelService) CreateModel(model *store.Model) error {
	return ms.Neo4jDomain.Save(model)
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

func (ms *ModelService) GetModelInstance(uuid string) (*store.Model, error) {
	model := &store.Model{}
	err := store.GetSession(true).Load(model, uuid)
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
