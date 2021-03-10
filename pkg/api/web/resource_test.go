package web

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/store"
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

func TestMyTest(t *testing.T) {
	//model := &[]store.Resource{}
	model := make([]store.Resource, 0)
	(&store.Neo4jDomain{}).Get(&model, "modelUid", "host")
	printOut(model)
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
}
