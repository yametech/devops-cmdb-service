package web

import (
	"encoding/json"
	"fmt"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/service"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/devops-cmdb-service/pkg/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGetModelAttributeList(t *testing.T) {
	baseUrl := "http://127.0.0.1:8080"
	// 表单数据
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=枯藤&age=18"
	// json
	contentType := "application/json"
	data := `{"modelUid":"host"}`
	resp, err := http.Post(baseUrl+"/cmdb/web/resource/model-attribute-list", contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed,err:%v\n", err)
		return
	}
	fmt.Println(string(b))
}
func printOut(obj interface{}) {
	b, _ := json.Marshal(obj)
	//fmt.Printf("%#v\n", string(b))
	fmt.Println(string(b))
}

func TestRelation(t *testing.T) {
	// RelationTest
	// MATCH (a:Model), (b:Model)
	//WHERE a.uid = 'host' AND b.uid = 'cabinet'
	//CREATE (a)-[:Relation {uid:"host_belong_cabinet", relationshipUid:"belong", constraint:"1 - 1", sourceUid:"host", targetUid:"cabinet"}]->(b);
	// query := "match (a:Model)-[r:Relation]-(b:Model) where r.uid = 'host_belong_cabinet' return distinct  r "
	query := "match (a:Model)-[r:Relation]-(b:Model) where a.uid = $modelUid or b.uid = $modelUid return distinct  r"

	session := store.GetSession(true)
	//relation := store.Relation{}
	result, _ := session.QueryRaw(query, map[string]interface{}{"modelUid": "hostss"})
	//fmt.Printf("%T\n", result[0][0])

	if result == nil {
		printOut("bingo")
	}
	for _, wrap := range result {
		relationship := wrap[0].(*gogm.RelationshipWrap)
		relation := &store.ModelRelation{}
		utils.SimpleConvert(relation, relationship.Props)
		printOut(relation)
	}
	printOut(result)
	relationship := result[0][0].(*gogm.RelationshipWrap)
	relation := &store.ModelRelation{}
	utils.SimpleConvert(relation, relationship.Props)
	printOut(relation)
	printOut(relationship.Props)

}

func TestMyTest(t *testing.T) {
	relationshipService := service.RelationshipService{}
	//res, _ := relationshipService.GetResourceRelationList("69906653-b834-4208-91e5-38fd262003d1")
	res, _ := relationshipService.GetResourceRelationList("6ae509e1-f6a7-438d-b336-c96d33b64e34")
	printOut(res)

	//resourceService := service.ResourceService{}
	//res := resourceService.GetModelAttributeList("host")
	//printOut(res)
}

func TestGetModel(t *testing.T) {
	session := store.GetSession(true)

	model := &[]store.Model{}
	err := session.LoadAllDepth(model, 2)
	if err != nil {
		panic(err)
	}

	marshal, _ := json.Marshal((*model)[0])
	//marshal, _ := json.Marshal(model)
	fmt.Println(string(marshal))
	//fmt.Printf("%#v\n", model)
}

func TestInit(t *testing.T) {
	// init
	session := store.GetSession(false)
	err := session.PurgeDatabase()
	if err != nil {
		panic(err)
	}
	//neo4j := store.Neo4jDomain{}

	// modelGroup
	modelGroup := &store.ModelGroup{Uid: "hardware", Name: "硬件资源"}
	//session.Save(modelGroup)

	model := &store.Model{Uid: "host", Name: "主机"}
	model.ModelGroup = modelGroup
	//session.Save(model)

	attributeGroup := &store.AttributeGroup{}
	json.Unmarshal([]byte("{\"modelUid\":\"host\",\"uid\":\"baseInfo\",\"name\":\"基本属性\"}"), attributeGroup)
	attributeGroup.Model = model
	//session.Save(attributeGroup)

	attributeCommon := &store.AttributeCommon{}
	jsonStr := "{\"uid\":\"ip\",\"name\":\"网址\",\"valueType\":\"短字符串\",\"editable\":true,\"required\":false,\"regular\":\"(([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\\\\.){3}([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\",\"comment\":\"网址信息\",\"modelUId\":\"host\",\"attributeGroupUid\":\"baseInfo\"}"
	json.Unmarshal([]byte(jsonStr), attributeCommon)

	attribute := &store.Attribute{AttributeCommon: *attributeCommon}
	attribute.AttributeGroup = attributeGroup
	session.SaveDepth(attribute, 10)

	attributeGroup2 := &store.AttributeGroup{}
	json.Unmarshal([]byte("{\"modelUid\":\"host\",\"uid\":\"otherInfo\",\"name\":\"其他属性\"}"), attributeGroup2)
	attributeGroup2.Model = model

	attributeCommon2 := &store.AttributeCommon{}
	jsonStr2 := "{\"uid\":\"test\",\"name\":\"cesi\",\"valueType\":\"短字符串\",\"editable\":true,\"required\":false,\"regular\":\"(([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\\\\.){3}([01]{0,1}\\\\d{0,1}\\\\d|2[0-4]\\\\d|25[0-5])\",\"comment\":\"网址信息\",\"modelUId\":\"host\",\"attributeGroupUid\":\"otherInfo\"}"
	json.Unmarshal([]byte(jsonStr2), attributeCommon2)

	attribute2 := &store.Attribute{AttributeCommon: *attributeCommon2}
	attribute2.AttributeGroup = attributeGroup2
	session.SaveDepth(attribute2, 10)

	// 使用服务
	resourceService := service.ResourceService{}
	body := "{\"modelUid\":\"host\",\"modelName\":\"主机\",\"attributeGroupIns\":[{\"uid\":\"baseInfo\",\"attributeIns\":[{\"uid\":\"ip\",\"attributeInsValue\":\"1\"},{\"uid\":\"cpu\",\"attributeInsValue\":\"12\"}]},{\"uid\":\"otherInfo\",\"attributeIns\":[{\"uid\":\"test\",\"attributeInsValue\":\"testttstset\"}]}]}"
	resourceResult, err := resourceService.AddResource(body, "")
	resourceResult3, err := resourceService.AddResource(body, "")

	model2 := &store.Model{Uid: "cabinet", Name: "机柜"}
	model2.ModelGroup = modelGroup
	session.SaveDepth(model2, 2)

	body2 := "{\"modelUid\":\"cabinet\",\"modelName\":\"机柜\",\"attributeGroupIns\":[{\"uid\":\"baseInfo\",\"attributeIns\":[{\"uid\":\"ip\",\"attributeInsValue\":\"1\"},{\"uid\":\"cpu\",\"attributeInsValue\":\"12\"}]},{\"uid\":\"otherInfo\",\"attributeIns\":[{\"uid\":\"test\",\"attributeInsValue\":\"testttstset\"}]}]}"
	resourceResult2, err := resourceService.AddResource(body2, "")

	relationshipService := service.RelationshipService{}
	resource1 := resourceResult.(*store.Resource)
	resource2 := resourceResult2.(*store.Resource)
	resource3 := resourceResult3.(*store.Resource)
	relationshipService.AddModelRelation("{\"relationshipUid\":\"belong\",\"constraint\":\"N - N\",\"sourceUid\":\"cabinet\",\"targetUid\":\"host\",\"comment\":\"描述信息\"}", "")

	relationshipService.AddResourceRelation(resource1.UUID, resource2.UUID, "host_belong_cabinet")
	relationshipService.AddResourceRelation(resource3.UUID, resource2.UUID, "host_belong_cabinet")
}
