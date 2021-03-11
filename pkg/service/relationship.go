package service

import "github.com/yametech/devops-cmdb-service/pkg/store"

type RelationshipService struct {
	Service
	store.Neo4jDomain
}

func (ms *RelationshipService) GetModelRelationList(modelUid string) interface{} {
	relations := make([]store.ModelRelation, 0)
	ms.Neo4jDomain.Get(&relations, "sourceUid", modelUid)

	relations2 := make([]store.ModelRelation, 0)
	ms.Neo4jDomain.Get(&relations2, "targetUid", modelUid)

	for _, r := range relations2 {
		relations = append(relations, r)
	}

	return relations
}

func (ms *RelationshipService) AddModelRelation(relationshipUid, constraint, sourceUid, targetUid, comment, operator string) (interface{}, error) {
	commonObj := &store.CommonObj{}
	commonObj.InitCommonObj(operator)
	relation := &store.ModelRelation{}
	relation.RelationshipUid = relationshipUid
	relation.Constraint = constraint
	relation.SourceUid = sourceUid
	relation.TargetUid = targetUid
	relation.Comment = comment
	relation.CommonObj = *commonObj
	err := ms.Neo4jDomain.Save(relation)
	return relation, err
}
