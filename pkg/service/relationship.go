package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"time"
)

type RelationshipService struct {
	Service
	store.Neo4jDomain
}

func (ms *RelationshipService) GetModelRelationList(uid string) interface{} {
	session := store.GetSession(true)
	query := "match (a:Model)-[r:Relation]-(b:Model) where a.uid = $modelUid or b.uid = $modelUid return distinct  r"
	result, _ := session.QueryRaw(query, map[string]interface{}{"modelUid": uid})

	relations := make([]store.ModelRelation, 0)
	for _, wrap := range result {
		relationship := wrap[0].(*gogm.RelationshipWrap)
		relation := &store.ModelRelation{}
		utils.SimpleConvert(relation, relationship.Props)
		relations = append(relations, *relation)
	}

	return &relations
}

func (ms RelationshipService) DeleteModelRelation(uid string) ([][]interface{}, error) {
	result, _ := ms.GetResourceRelationByModelRelationUid(uid)
	if result != nil {
		return nil, errors.New("该模型已被使用，禁止删除")
	}

	query := "match (a:Model)-[r:Relation]-(b:Model) where r.uid = $uid delete  r"
	return ms.ManualExecute(query, map[string]interface{}{"uid": uid})
}

func (ms *RelationshipService) UpdateModelRelation(body string, operator string) (interface{}, error) {
	src, _ := parseToModelRelation(body, operator)
	if src == nil || src.Uid == "" {
		return nil, errors.New("更新模型关系缺少uid")
	}

	result, err := ms.GetModelRelationByUid(src.Uid)
	if result != nil && len(result) == 0 && len(result[0]) == 0 && result[0][0] == nil {
		return nil, errors.New("不存在该模型关系")
	}

	modelRelation := result[0][0].(*store.ModelRelation)

	result, _ = ms.GetResourceRelationByModelRelationUid(src.Uid)
	//如果已有数据关联此模型，则只能更新描述备注
	if result == nil {
		// 全部更新
		modelRelation.RelationshipUid = src.RelationshipUid
		modelRelation.TargetUid = src.TargetUid
		modelRelation.SourceUid = src.SourceUid
		modelRelation.Constraint = src.Constraint
	}

	modelRelation.Comment = src.Comment
	modelRelation.Editor = operator
	modelRelation.UpdateTime = time.Now().Unix()

	err = ms.Neo4jDomain.Update(modelRelation)
	if err != nil {
		return nil, err
	}

	return modelRelation, nil
}

func parseToModelRelation(body string, operator string) (*store.ModelRelation, error) {
	fmt.Println(body)
	bodyObj := &store.ModelRelation{}
	err := json.Unmarshal([]byte(body), bodyObj)
	if err != nil {
		return nil, err
	}

	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	relation := &store.ModelRelation{}
	relation.RelationshipUid = bodyObj.RelationshipUid
	relation.Constraint = bodyObj.Constraint
	relation.SourceUid = bodyObj.SourceUid
	relation.TargetUid = bodyObj.TargetUid
	relation.Comment = bodyObj.Comment
	relation.CommonObj = *commonObj

	return relation, nil
}

func (ms *RelationshipService) GetResourceRelationByModelRelationUid(modelRelationUid string) ([][]interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where r.modelRelationUid = $uid return distinct  r"
	return ms.ManualQueryRaw(query, map[string]interface{}{"uid": modelRelationUid})
}

func (ms *RelationshipService) GetModelRelationByUid(uid string) ([][]interface{}, error) {
	query := "match (a:Model)-[r:Relation]-(b:Model) where r.uid = $uid return distinct  r"
	return ms.ManualQueryRaw(query, map[string]interface{}{"uid": uid})
}

func (ms *RelationshipService) AddModelRelation(body string, operator string) (interface{}, error) {
	relation, err := parseToModelRelation(body, operator)
	if err != nil {
		return nil, err
	}
	relation.Uid = relation.SourceUid + "_" + relation.RelationshipUid + "_" + relation.TargetUid

	result, err := ms.GetModelRelationByUid(relation.Uid)
	if result != nil {
		return nil, errors.New("已存在该模型关系")
	}

	query := "MATCH (a:Model), (b:Model) WHERE a.uid = $sourceUid AND b.uid = $targetUid " +
		"CREATE (a)-[:Relation {uid: $uid, relationshipUid: $relationshipUid, constraint: $constraint, sourceUid: $sourceUid, targetUid: $targetUid, comment: $comment }]->(b)"
	result, err = ms.ManualExecute(query, map[string]interface{}{"sourceUid": relation.SourceUid, "targetUid": relation.TargetUid,
		"uid": relation.Uid, "relationshipUid": relation.RelationshipUid, "constraint": relation.Constraint, "comment": relation.Comment})
	return result, err
}
