package service

import (
	"errors"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"time"
)

type RelationService struct {
	Service
}

func (rs *RelationService) GetAllModelRelations() *[]common.ModelRelationVO {
	query := "match (a:Model)-[r:Relation]-(b:Model) return distinct  r"
	return rs.queryModelRelations(query, nil)
}

func (rs *RelationService) queryModelRelations(query string, properties map[string]interface{}) *[]common.ModelRelationVO {
	session := rs.GetSession(true)
	defer session.Close()
	result, _ := session.QueryRaw(query, properties)

	relations := make([]common.ModelRelationVO, 0)
	relationshipModelMap := make(map[string]store.RelationshipModel)
	for _, wrap := range result {
		relationshipWrap := wrap[0].(*gogm.RelationshipWrap)
		relation := &common.ModelRelationVO{}
		utils.SimpleConvert(relation, relationshipWrap.Props)
		relation.Id = relationshipWrap.Id

		if len(wrap) == 3 {
			// 模型关系列表需要查看关系模型名称
			relationshipModel, ok := relationshipModelMap[relation.RelationshipUid]
			if ok {
				relation.RelationshipName = relationshipModel.Name
			} else {
				model := &store.RelationshipModel{}
				err := rs.Neo4jDomain.Get(model, "uid", relation.RelationshipUid)
				if err != nil {
					fmt.Println(err)
				}
				relationshipModelMap[relation.RelationshipUid] = *model
				relation.RelationshipName = model.Name
			}

			// 模型信息
			nodeWrap := wrap[1].(*gogm.NodeWrap)
			if relation.SourceUid == nodeWrap.Props["uid"] {
				relation.SourceName = nodeWrap.Props["name"].(string)
				nodeWrap := wrap[2].(*gogm.NodeWrap)
				relation.TargetName = nodeWrap.Props["name"].(string)
			} else {
				relation.TargetName = nodeWrap.Props["name"].(string)
				nodeWrap := wrap[2].(*gogm.NodeWrap)
				relation.SourceName = nodeWrap.Props["name"].(string)
			}
		}
		relations = append(relations, *relation)
	}

	return &relations
}

func (rs *RelationService) GetModelRelationList(uid string) interface{} {
	query := "match (a:Model)-[r:Relation]-(b:Model) where a.uid = $modelUid or b.uid = $modelUid return distinct r, a, b"
	result := rs.queryModelRelations(query, map[string]interface{}{"modelUid": uid})
	// 处理重复的
	IdMap := map[int64]bool{}
	relations := make([]common.ModelRelationVO, 0)
	for _, vo := range *result {
		if !IdMap[vo.Id] {
			relations = append(relations, vo)
			IdMap[vo.Id] = true
		}
	}

	return relations
}

func (rs RelationService) DeleteModelRelation(uid string) ([][]interface{}, error) {
	result, err := rs.GetResourceRelationsByModelRelationUid(uid)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return nil, errors.New("该模型已被使用，禁止删除")
	}

	query := "match (a:Model)-[r:Relation]-(b:Model) where r.uid = $uid delete  r"
	return rs.ManualExecute(query, map[string]interface{}{"uid": uid})
}

func (rs *RelationService) UpdateModelRelation(vo *common.UpdateModelRelationVO, operator string) (interface{}, error) {
	src, _ := parseToModelRelation(vo, operator)

	result, err := rs.GetModelRelation("uid", src.Uid)
	if result == nil || len(result) == 0 || len(result[0]) == 0 && result[0][0] == nil {
		return nil, errors.New("不存在该模型关系")
	}

	updateCypher := "MATCH (a:Model)-[r:Relation]-(b:Model) WHERE r.uid = $uid " +
		"SET r.comment = $comment , r.updateTime = $updateTime , r.editor = $editor "

	properties := map[string]interface{}{}
	properties["uid"] = src.Uid
	properties["comment"] = src.Comment
	properties["editor"] = src.Editor
	properties["updateTime"] = time.Now().UnixNano() / 1000000
	//
	result, _ = rs.GetResourceRelationsByModelRelationUid(src.Uid)
	////如果已有数据关联此模型，则只能更新描述备注
	if result == nil {
		//	// 全部更新
		updateCypher += " , r.uid = $newUid , r.relationshipUid = $relationshipUid , r.targetUid = $targetUid , r.sourceUid = $sourceUid , r.constraint = $constraint"
		properties["relationshipUid"] = src.RelationshipUid
		properties["sourceUid"] = src.SourceUid
		properties["targetUid"] = src.TargetUid
		properties["constraint"] = src.Constraint
		properties["newUid"] = src.SourceUid + "_" + src.RelationshipUid + "_" + src.TargetUid
	}

	fmt.Println(updateCypher, properties)
	result, err = rs.ManualExecute(updateCypher, properties)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseToModelRelation(vo interface{}, operator string) (*store.ModelRelation, error) {
	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	relation := &store.ModelRelation{}
	utils.SimpleConvert(relation, vo)
	relation.CommonObj = *commonObj

	return relation, nil
}

func (rs *RelationService) GetResourceRelationsByModelRelationUid(modelRelationUid string) ([][]interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where r.uid = $uid return distinct  r"
	return rs.ManualQueryRaw(query, map[string]interface{}{"uid": modelRelationUid})
}

func (rs *RelationService) GetModelRelation(key, value string) ([][]interface{}, error) {
	query := "match (a:Model)-[r:Relation]-(b:Model) where r." + key + " = $value return distinct  r"
	return rs.ManualQueryRaw(query, map[string]interface{}{"value": value})
}

func (rs *RelationService) AddModelRelation(vo *common.AddModelRelationVO, operator string) (interface{}, error) {
	relationshipModel := &store.RelationshipModel{}
	if err := rs.Get(relationshipModel, "uid", vo.RelationshipUid); err != nil {
		if relationshipModel.UUID == "" {
			return nil, fmt.Errorf("关系模型%q不存在或已被删除", vo.RelationshipUid)
		}
		return nil, err
	}

	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	relation, err := parseToModelRelation(vo, operator)
	if err != nil {
		return nil, err
	}

	relation.Uid = relation.SourceUid + "_" + relation.RelationshipUid + "_" + relation.TargetUid
	result, err := rs.GetModelRelation("uid", relation.Uid)
	if result != nil {
		return nil, errors.New("已存在该模型关系")
	}

	query := "MATCH (a:Model), (b:Model) WHERE a.uid = $sourceUid AND b.uid = $targetUid " +
		"CREATE (a)-[:Relation {uid: $uid, relationshipUid: $relationshipUid, constraint: $constraint, sourceUid: $sourceUid, targetUid: $targetUid, comment: $comment }]->(b)"
	result, err = rs.ManualExecute(query, map[string]interface{}{"sourceUid": relation.SourceUid, "targetUid": relation.TargetUid,
		"uid": relation.Uid, "relationshipUid": relation.RelationshipUid, "constraint": relation.Constraint, "comment": relation.Comment})
	return result, err
}

func (rs *RelationService) GetResourceRelationList(uuid string) (interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where a.uuid = $uuid return a,r,b order by ID(r) ASC"
	result, err := rs.ManualQueryRaw(query, map[string]interface{}{"uuid": uuid})

	voList := make([]common.ResourceRelationListPageVO, 0)
	for _, row := range result {
		addResourceRelationList(&voList, row, uuid)
	}

	// 获取最新的模型名称
	modelMap := make(map[string]store.Model)
	for i := 0; i < len(voList); i++ {
		model, ok := modelMap[voList[i].SourceUid]
		if ok {
			voList[i].SourceName = model.Name
		} else {
			model = store.Model{}
			err := rs.Get(&model, "uid", voList[i].SourceUid)
			if err != nil {
				fmt.Println(err)
			}
			if model.UUID != "" {
				voList[i].SourceName = model.Name
				modelMap[voList[i].SourceUid] = model
			}
		}

		model, ok = modelMap[voList[i].TargetUid]
		if ok {
			voList[i].TargetName = model.Name
		} else {
			model = store.Model{}
			err := rs.Get(&model, "uid", voList[i].TargetUid)
			if err != nil {
				fmt.Println(err)
			}
			if model.UUID != "" {
				voList[i].TargetName = model.Name
				modelMap[voList[i].TargetUid] = model
			}
		}
	}

	return voList, err
}

func addResourceRelationList(result *[]common.ResourceRelationListPageVO, row []interface{}, uuid string) {
	pageVO := convert2ResourceRelationListPageVO(row, uuid)

	newRelation := true
	for _, vo := range *result {
		if vo.RelationshipUid == pageVO.RelationshipUid {
			newRelation = false
			// 资源信息添加进去
			*vo.Resources = append(*vo.Resources, (*pageVO.Resources)[0])
		}
	}
	// 新的数据
	if newRelation {
		// 资源字段
		//resourceService := &ResourceService{}
		modelAttributes := &[]common.ModelAttributeVisibleVO{}
		//utils.SimpleConvert(modelAttributes, resourceService.GetModelAttributeList((*pageVO.Resources)[0]["modelUid"]))
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
	//resourceService := ResourceService{}
	//res, _ := resourceService.GetResourceDetail(resource["uuid"])
	//for _, g := range res.AttributeGroupIns {
	//	for _, attribute := range g.AttributeIns {
	//		resource[attribute.Uid] = attribute.AttributeInsValue
	//	}
	//}

	return vo
}

func (rs *RelationService) DeleteResourceRelation(sourceUUID, targetUUID, uid string) ([][]interface{}, error) {
	query := "match (a:Resource)-[r:Relation]-(b:Resource) where r.uid = $uid and a.uuid = $sourceUUID and b.uuid = $targetUUID delete  r"
	return rs.ManualExecute(query, map[string]interface{}{"uid": uid, "sourceUUID": sourceUUID, "targetUUID": targetUUID})
}

func (rs *RelationService) AddResourceRelation(sourceUUID, targetUUID, uid string) ([][]interface{}, error) {
	modelRelation, err := rs.GetModelRelation("uid", uid)
	if err != nil {
		return nil, err
	}
	if modelRelation == nil || len(modelRelation) == 0 || len(modelRelation[0]) == 0 && modelRelation[0][0] == nil {
		return nil, errors.New("不存在该模型关系")
	}

	result, err := rs.GetResourceRelation(sourceUUID, targetUUID, uid)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return nil, errors.New("已存在该资源关系")
	}

	query := "MATCH (a:Resource), (b:Resource) WHERE a.uuid = $sourceUUID AND b.uuid= $targetUUID " +
		"CREATE (a)-[:Relation {uid: $uid}]->(b)"
	result, err = rs.ManualExecute(query, map[string]interface{}{"uid": uid, "sourceUUID": sourceUUID, "targetUUID": targetUUID})
	return result, err
}

func (rs *RelationService) GetResourceRelation(sourceUUID, targetUUID, uid string) ([][]interface{}, error) {
	query := "match (a:Resource)-[r:Relation]->(b:Resource) where r.uid = $uid and a.uuid = $sourceUUID and b.uuid = $targetUUID return distinct  r"
	return rs.ManualQueryRaw(query, map[string]interface{}{"uid": uid, "sourceUUID": sourceUUID, "targetUUID": targetUUID})
}
