package service

import (
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"regexp"
	"strings"
	"time"
)

type ModelService struct {
	Service
}

func (ms *ModelService) GetAllModelGroup() ([]*store.ModelGroup, error) {
	modelGroups := make([]*store.ModelGroup, 0)
	query := "MATCH (a:ModelGroup) RETURN a ORDER BY a.createTime ASC"
	result, err := ms.ManualQueryRaw(query, nil)

	for _, row := range result {
		modelGroup := &store.ModelGroup{}
		nodeWrap := row[0].(*gogm.NodeWrap)
		utils.SimpleConvert(modelGroup, &nodeWrap.Props)

		models := make([]*store.Model, 0)
		query = "MATCH (a:ModelGroup {uuid: $uuid})<-[]-(b:Model) RETURN b ORDER BY b.createTime ASC"
		_ = ms.ManualQuery(query, map[string]interface{}{"uuid": modelGroup.UUID}, &models)
		modelGroup.Models = models
		addModelGroup(&modelGroups, modelGroup)
	}

	return modelGroups, err
}

func addModelGroup(self *[]*store.ModelGroup, target *store.ModelGroup) {
	for _, group := range *self {
		if group.Uid == target.Uid {
			for _, model := range target.Models {
				group.AddModel(model)
			}
			return
		}
	}
	*self = append(*self, target)
}

func (ms *ModelService) GetModelGroup(uuid string) (*store.ModelGroup, error) {
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", uuid); err != nil {
		return nil, err
	}
	return modelGroup, nil
}

func (ms *ModelService) DeleteModelGroup(uuid string) error {
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", uuid); err != nil {
		if modelGroup.UUID == "" {
			return errors.New("模型分组已被删除")
		}
		return err
	}

	session := ms.GetSession(true)
	defer session.Close()
	session.LoadDepth(modelGroup, uuid, 1)
	if len(modelGroup.Models) > 0 {
		return errors.New("此分组下存在模型，不可删除")
	}

	wSession := ms.GetSession(false)
	defer wSession.Close()
	if err := wSession.DeleteUUID(uuid); err != nil {
		return err
	}
	return nil
}

func UidNameValidate(uid, name string) error {
	uid = strings.TrimSpace(uid)
	regexpUid := "^[a-z]([_a-z0-9]{0,29})$"
	match, err := regexp.MatchString(regexpUid, uid)
	if err != nil {
		return err
	}
	if !match {
		return fmt.Errorf("唯一标识%q不符合规范: 小写英文开头，下划线，数字，小写英文的组合，且长度不超过30位", uid)
	}

	name = strings.TrimSpace(name)
	if len(name) < 1 || len([]rune(name)) > 30 {
		return fmt.Errorf("名称%q不符合规范：长度不超过30位", name)
	}

	return nil
}

func (ms *ModelService) CreateModelGroup(vo *common.AddModelGroupVO, operator string) (*store.ModelGroup, error) {
	vo.Uid = strings.TrimSpace(vo.Uid)
	vo.Name = strings.TrimSpace(vo.Name)
	// validate
	if err := UidNameValidate(vo.Uid, vo.Name); err != nil {
		return nil, err
	}

	ms.mutex.Lock()
	ms.mutex.Unlock()
	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uid", vo.Uid); err == nil {
		return nil, fmt.Errorf("已存在唯一标识是%q的模型分组", vo.Uid)
	}

	if err := ms.Neo4jDomain.Get(modelGroup, "name", vo.Name); err == nil {
		return nil, fmt.Errorf("已存在名称是%q的模型分组", vo.Name)
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

	modelGroup.Name = strings.TrimSpace(vo.Name)
	modelGroup.UpdateTime = time.Now().UnixNano() / 1000000
	modelGroup.CommonObj.Editor = operator

	err := ms.Neo4jDomain.Update(modelGroup)
	return modelGroup, err
}

func (ms *ModelService) CreateModel(vo *common.AddModelVO, operator string) (*store.Model, error) {
	vo.ModelGroupUUID = strings.TrimSpace(vo.ModelGroupUUID)
	vo.Uid = strings.TrimSpace(vo.Uid)
	vo.Name = strings.TrimSpace(vo.Name)
	// validate
	if err := UidNameValidate(vo.Uid, vo.Name); err != nil {
		return nil, err
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	model := &store.Model{}
	if err := ms.Neo4jDomain.Get(model, "uid", vo.Uid); err == nil {
		return nil, fmt.Errorf("已存在唯一标识是%q的模型", vo.Uid)
	}

	if err := ms.Neo4jDomain.Get(model, "name", vo.Name); err == nil {
		return nil, fmt.Errorf("已存在名称是%q的模型", vo.Name)
	}

	utils.SimpleConvert(model, vo)

	modelGroup := &store.ModelGroup{}
	if err := ms.Neo4jDomain.Get(modelGroup, "uuid", vo.ModelGroupUUID); err != nil {
		if modelGroup.UUID == "" {
			return nil, errors.New("对应模型分组已被删除")
		}
		return nil, err
	}

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	model.CommonObj = *commonObj
	model.ModelGroup = modelGroup
	err := ms.Neo4jDomain.Save(model)
	return model, err
}

func (ms *ModelService) UpdateModel(vo *common.UpdateModelVO, operator string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	vo.Name = strings.TrimSpace(vo.Name)
	vo.ModelGroupUUID = strings.TrimSpace(vo.ModelGroupUUID)
	model := &store.Model{}
	err := ms.Get(model, "uuid", vo.UUID)
	if err != nil {
		if model.UUID == "" {
			return errors.New("模型已被删除")
		}
		return err
	}

	session := ms.GetSession(false)
	defer session.Close()
	// 更新名称
	if model.Name != vo.Name {
		models := make([]store.Model, 0)
		err := ms.Get(&models, "name", vo.Name)
		if err != nil {
			fmt.Println(err)
		}

		// 校验
		if len(models) > 1 {
			return errors.New("已存在该名称的模型")
		}
		if len(models) == 1 {
			if models[0].UUID != vo.UUID {
				return errors.New("已存在该名称的模型")
			} else {
				return nil
			}
		}
	}

	// 更新模型分组
	if vo.ModelGroupUUID != "" {
		modelGroup := &store.ModelGroup{}
		err := ms.Get(modelGroup, "uuid", vo.ModelGroupUUID)
		if err != nil {
			if modelGroup.UUID == "" {
				return fmt.Errorf("模型分组已被删除")
			}
			return err
		}

		// 删除关系
		session.QueryRaw("MATCH (:Model{uid:$modelUid})-[r1:GroupBy]-(:ModelGroup) DELETE r1",
			map[string]interface{}{"modelUid": vo.Uid})

		// 新增关系
		session.QueryRaw("match (n:Model{uid:$modelUid}),(m:ModelGroup{uuid:$modelGroupUUID}) create (n)-[r2:GroupBy]->(m)",
			map[string]interface{}{"modelUid": vo.Uid, "modelGroupUUID": vo.ModelGroupUUID})
	}
	// 更新模型
	model.Name = vo.Name
	model.UpdateTime = time.Now().UnixNano() / 1000000
	model.Editor = operator
	return session.Save(model)
}

func (ms *ModelService) DeleteModel(uuid, operator string) error {
	model := &store.Model{}
	err := ms.Neo4jDomain.Get(model, "uuid", uuid)
	if err != nil {
		if model.UUID == "" {
			return errors.New("模型分组已被删除")
		}
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

func (ms *ModelService) GetSimpleModelList() ([]*store.Model, error) {
	allModel := make([]*store.Model, 0)
	query := fmt.Sprintf("match (a:Model) return a ORDER BY a.createTime ASC")
	properties := map[string]interface{}{}
	session := ms.GetSession(true)
	defer session.Close()
	if err := session.Query(query, properties, &allModel); err != nil {
		return allModel, err
	}
	return allModel, nil
}

func (ms *ModelService) GetModel(uuid string) (*store.Model, error) {
	model := &store.Model{}
	err := ms.Get(model, "uuid", uuid)
	if err != nil {
		return nil, errors.New("该模型已被删除")
	}

	attributeGroups := &[]*store.AttributeGroup{}
	query := "MATCH (a:Model {uuid: $uuid})-[]-(b:AttributeGroup) RETURN b ORDER BY b.createTime ASC"
	_ = ms.ManualQuery(query, map[string]interface{}{"uuid": uuid}, attributeGroups)

	model.AttributeGroups = *attributeGroups

	for _, ag := range *attributeGroups {
		attributes := &[]*store.Attribute{}
		query = "MATCH (a:AttributeGroup {uuid: $uuid})-[]-(b:Attribute) RETURN b ORDER BY b.createTime ASC"
		_ = ms.ManualQuery(query, map[string]interface{}{"uuid": ag.UUID}, attributes)
		ag.Attributes = *attributes
	}
	return model, err
}

func (ms *ModelService) GetModelDetail(uid string) (*store.Model, error) {
	model := &store.Model{}
	err := ms.Get(model, "uid", uid)
	if err != nil {
		return nil, errors.New("该模型已被删除")
	}

	attributeGroups := &[]*store.AttributeGroup{}
	query := "MATCH (a:Model {uid: $uid})-[]-(b:AttributeGroup) RETURN b ORDER BY b.createTime ASC"
	_ = ms.ManualQuery(query, map[string]interface{}{"uid": uid}, attributeGroups)

	model.AttributeGroups = *attributeGroups

	for _, ag := range *attributeGroups {
		attributes := &[]*store.Attribute{}
		query = "MATCH (a:AttributeGroup {uuid: $uuid})-[]-(b:Attribute) RETURN b ORDER BY b.createTime ASC"
		_ = ms.ManualQuery(query, map[string]interface{}{"uuid": ag.UUID}, attributes)
		ag.Attributes = *attributes
	}
	return model, err
}

func (ms *ModelService) GetRelationshipList(pageSize, pageNumber int) ([]*store.RelationshipModel, error) {
	allModel := make([]*store.RelationshipModel, 0)
	query := fmt.Sprintf("match (a:RelationshipModel) return a ORDER BY a.createTime ASC SKIP $skip LIMIT $limit")
	properties := map[string]interface{}{
		"skip":  (pageNumber - 1) * pageSize,
		"limit": pageSize,
	}
	session := ms.GetSession(true)
	defer session.Close()
	if err := session.Query(query, properties, &allModel); err != nil {
		return nil, err
	}
	// get all ModelRelation, count the Relationship usage
	relationService := RelationService{}
	relations := relationService.GetAllModelRelations()
	for _, model := range allModel {
		for _, relation := range *relations {
			if relation.RelationshipUid == model.Uid {
				model.CurrentUsage += 1
			}
		}
	}

	return allModel, nil
}

func (ms *ModelService) GetRelationship(uuid string) (*store.RelationshipModel, error) {
	relationship := &store.RelationshipModel{}
	if err := ms.Neo4jDomain.Get(relationship, "uuid", uuid); err != nil {
		return nil, err
	}
	return relationship, nil
}

func (ms *ModelService) SaveRelationship(vo *common.CreateRelationshipModelVO, operator string) (*store.RelationshipModel, error) {
	vo.Name = strings.TrimSpace(vo.Name)
	vo.Uid = strings.TrimSpace(vo.Uid)
	vo.Source2Target = strings.TrimSpace(vo.Source2Target)
	vo.Target2Source = strings.TrimSpace(vo.Target2Source)
	// validate
	if err := UidNameValidate(vo.Uid, vo.Name); err != nil {
		return nil, err
	}

	if len([]rune(vo.Source2Target)) > 5 || len([]rune(vo.Source2Target)) < 1 {
		return nil, fmt.Errorf("%q不满足位数限制", vo.Source2Target)
	}

	if len([]rune(vo.Target2Source)) > 5 || len([]rune(vo.Target2Source)) < 1 {
		return nil, fmt.Errorf("%q不满足位数限制", vo.Target2Source)
	}

	exist := &store.RelationshipModel{}
	err := ms.Neo4jDomain.Get(exist, "uid", vo.Uid)
	if err == nil {
		return nil, fmt.Errorf("已存在唯一标识是%q的关系模型", vo.Uid)
	}

	err = ms.Neo4jDomain.Get(exist, "name", vo.Name)
	if err == nil {
		return nil, fmt.Errorf("已存在名称是%q的关系模型", vo.Name)
	}

	relationship := &store.RelationshipModel{}
	utils.SimpleConvert(relationship, vo)
	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	relationship.CommonObj = *commonObj

	err = ms.Neo4jDomain.Save(relationship)
	if err != nil {
		return nil, err
	}
	return relationship, nil
}

func (ms *ModelService) UpdateRelationship(vo *common.UpdateRelationshipModelVO, operator string) error {
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
	relation.UpdateTime = time.Now().UnixNano() / 1000000
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
	query := "match (a:Model)-[r:Relation]-(b:Model) where r.relationshipUid = $relationshipUid return COUNT(distinct r)"
	count, err := ms.ManualQueryRaw(query, map[string]interface{}{"relationshipUid": relationshipModel.Uid})
	if err != nil {
		return err
	}
	if count[0][0].(int64) > 0 {
		return errors.New("该关系模型已被使用，禁止删除")
	}

	return ms.Neo4jDomain.Delete(relationshipModel)
}
