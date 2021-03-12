package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/common"
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

func (ms RelationshipService) GetResourceRelationList(uuid string) (interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where a.uuid = $uuid return a,r,b"
	result, err := ms.ManualQueryRaw(query, map[string]interface{}{"uuid": uuid})
	printOut(result)

	voList := &[]common.ResourceRelationListPageVO{}
	for _, row := range result {
		addResourceRelationList(voList, row, uuid)
	}

	return voList, err
}

func addResourceRelationList(result *[]common.ResourceRelationListPageVO, row []interface{}, uuid string) {
	if row == nil {
		return
	}
	if result == nil {
		r := make([]common.ResourceRelationListPageVO, 0)
		result = &r
	}
	pageVO := convert2ResourceRelationListPageVO(row, uuid)

	newRelation := false
	for _, vo := range *result {
		if vo.SourceUid == pageVO.SourceUid {
			newRelation = true
			// 资源信息添加进去
			*vo.Resources = append(*vo.Resources, (*pageVO.Resources)[0])
		}
	}
	// 新的数据
	if !newRelation {
		// 资源字段
		resourceService := &ResourceService{}
		modelAttributes := &[]common.ModelAttributeVisibleVO{}
		utils.SimpleConvert(modelAttributes, resourceService.GetModelAttributeList((*pageVO.Resources)[0]["modelUid"]))
		pageVO.ModelAttributes = modelAttributes

		*result = append(*result, *pageVO)
	}
}

func convert2ResourceRelationListPageVO(row []interface{}, uuid string) *common.ResourceRelationListPageVO {
	a := row[0].(*gogm.NodeWrap)
	r := row[1].(*gogm.RelationshipWrap)
	b := row[2].(*gogm.NodeWrap)
	startSource := &store.Resource{}
	endSource := &store.Resource{}

	// 根据关系信息找到方向,比如：a-r->b
	if r.StartId == a.Id {
		utils.SimpleConvert(startSource, &a.Props)
		utils.SimpleConvert(endSource, &b.Props)
	} else {
		utils.SimpleConvert(endSource, &a.Props)
		utils.SimpleConvert(startSource, &b.Props)
	}
	vo := &common.ResourceRelationListPageVO{}
	vo.SourceUid = startSource.ModelUid
	vo.SourceName = startSource.ModelName
	vo.TargetUid = endSource.ModelUid
	vo.TargetName = endSource.ModelName
	vo.RelationshipUid = r.Props["uid"].(string)

	// 关联资源实例
	resource := map[string]string{}
	if startSource.UUID == uuid {
		resource["uuid"] = endSource.UUID
		resource["modelUid"] = endSource.ModelUid
	} else {
		resource["uuid"] = startSource.UUID
		resource["modelUid"] = startSource.ModelUid
	}
	resources := make([]map[string]string, 0)
	resources = append(resources, resource)
	vo.Resources = &resources
	//补充资源实例
	resourceService := ResourceService{}
	res, _ := resourceService.GetResourceDetail(resource["uuid"])
	for _, g := range res.(*store.Resource).AttributeGroupIns {
		for _, attribute := range g.AttributeIns {
			resource[attribute.Uid] = attribute.AttributeInsValue
		}
	}

	return vo
}

func (ms RelationshipService) DeleteResourceRelation(sourceUUID, targetUUID, uid string) ([][]interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where r.uid = $uid and a.uuid = $sourceUid and b.uuid = $targetUid delete  r"
	return ms.ManualExecute(query, map[string]interface{}{"uid": uid, "sourceUUID": sourceUUID, "targetUUID": targetUUID})
}

func (ms RelationshipService) AddResourceRelation(sourceUUID, targetUUID, uid string) ([][]interface{}, error) {
	result, err := ms.GetResourceRelation(sourceUUID, targetUUID, uid)
	if result != nil {
		return nil, errors.New("已存在该资源关系")
	}

	query := "MATCH (a:Resource), (b:Resource) WHERE a.uuid = $sourceUUID AND b.uuid= $targetUUID " +
		"CREATE (a)-[:Relation {uid: $uid}]->(b)"
	result, err = ms.ManualExecute(query, map[string]interface{}{"uid": uid, "sourceUUID": sourceUUID, "targetUUID": targetUUID})
	return result, err
}

func (ms RelationshipService) GetResourceRelation(sourceUUID, targetUUID, uid string) ([][]interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where r.uid = $uid and a.uuid = $sourceUUID and b.uuid = $targetUUID return distinct  r"
	return ms.ManualQueryRaw(query, map[string]interface{}{"uid": uid, "sourceUUID": sourceUUID, "targetUUID": targetUUID})
}
