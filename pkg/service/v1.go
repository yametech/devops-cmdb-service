package service

import (
	"github.com/yametech/devops-cmdb-service/pkg/common"
)

type V1 struct {
	Service
}

func (s *V1) GetAppTree() (interface{}, error) {
	cypher := `
			MATCH (a:Resource {modelUid:'business'})<-[]-(b:AttributeGroupIns {uid:'business_info'})<-[]-(c:AttributeIns)
			WHERE c.uid in ['business_name','business_master']
			OPTIONAL MATCH (a)-[:Relation {uid:'business_including_business_domain'}]-(a2:Resource {modelUid:'business_domain'})<-[]-(b2:AttributeGroupIns)<-[]-(c2:AttributeIns {uid:'domain_name'}) 
			OPTIONAL MATCH (a2)-[:Relation {uid:'business_domain_including_business_service'}]-(a3:Resource {modelUid:'business_service'})<-[]-(b3:AttributeGroupIns)<-[]-(c3:AttributeIns) 
			RETURN id(a),c.uid,c.attributeInsValue,id(a2),c2.attributeInsValue,id(a3),c3.uid,c3.attributeInsValue
			`

	raw, err := s.ManualQueryRaw(cypher, nil)
	if err != nil {
		return nil, err
	}
	businessMap := map[int64]*common.Business{}
	for _, row := range raw {
		id := row[0].(int64)
		b := businessMap[id]
		if businessMap[id] == nil {
			b = &common.Business{Id: id}
			businessMap[id] = b
		}
		b.AddAttribute(row[1].(string), row[2].(string))
		if row[3] != nil {
			b.AddDomain(&common.Domain{Id: row[3].(int64), Name: row[4].(string)})
			if row[5] != nil {
				b.AddService(row[3].(int64), row[5].(int64), row[6].(string), row[7].(string))
			}
		}
	}
	businessList := make([]*common.Business, 0)
	for _, business := range businessMap {
		if business.Children == nil {
			business.Children = []*common.Domain{}
		}
		for _, child := range business.Children {
			if child.Children == nil {
				child.Children = []*common.Service{}
			}
		}
		businessList = append(businessList, business)
	}
	return businessList, nil
}
