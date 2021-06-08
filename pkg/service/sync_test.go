package service

import (
	"encoding/json"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"log"
	"testing"
)

func TestSyncDomainIns(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	var host, username, password = "10.200.10.51", "neo4j", "test123qwe"
	store.Neo4jInit(host, username, password)
	service := SyncService{}
	result, err := service.SyncAliDomain("sync")
	if err != nil {
		panic(err)
	}
	log.Println(result)

}

func TestSyncResource(t *testing.T) {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	var host, username, password = "localhost", "neo4j", "123456"
	store.Neo4jInit(host, username, password)
	//service := ResourceService{}
	//attributes := []map[string]interface{}{
	//	{"uid":"test", "attributeInsValue":234},
	//	{"uid":"test2", "attributeInsValue":234},
	//	{"uid":"test3", "attributeInsValue":234},
	//}
	//service.GenADDResourceVO("room","base",attributes)
	service := SyncService{}
	log.Println(service.SyncResource("built_domain", "sss"))

}

func TestFixUserName(t *testing.T) {
	query := "match (a:Resource {modelUid:'business'})-[]-(b:AttributeGroupIns)-[]-(c:AttributeIns {uid:'business_master'}) return id(c),c.attributeInsValue"
	common.LdapAuthPassword = "h1W8QZLdoa"
	log.SetFlags(log.Llongfile | log.LstdFlags)
	//var host, username, password = "10.200.10.51", "neo4j", "test123qwe"
	var host, username, password = "10.200.64.10", "neo4j", "Nesf2Ld"
	store.Neo4jInit(host, username, password)
	service := SyncService{}
	result, err := service.ResourceService.ManualQueryRaw(query, nil)
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(result)
	fmt.Println(string(marshal))
	ldapService := &LdapService{}
	ldapUserMap := ldapService.GetLdapUserMap()
	//fmt.Println(ldapUserMap)
	for _, row := range result {
		if user, ok := ldapUserMap[row[1].(string)]; ok && user.Name != "" {
			fmt.Println(fmt.Sprintf("match (c:AttributeIns) where id(c) = %v set c.attributeInsValue='%v';", row[0], user.Name))
			//r, err := service.ResourceService.ManualExecute("match (c:AttributeIns) where id(c) = $id set c.attributeInsValue=$value", map[string]interface{}{"id": row[0], "value": user.Name})
			//log.Println(r, err)
		}
	}
}

func TestK8s(t *testing.T) {
	vo := common.K8sResource{}
	vo.CompassName = "compass开发环境"
	vo.ResourceName = "Pod"
	vo.ResourceAttribute = map[string]string{
		"name":            "worker",
		"cpu_cores":       "112",
		"cpu_request":     "87.3",
		"cpu_limits":      "387.99",
		"memory_capacity": "581.1",
		"memory_request":  "87.3",
		"pod_capacity":    "5000",
	}
	vo.ResourceRelation = map[string][]string{
		"Stone": {"sdf"},
		"PVC":   {"1", "2"},
	}

	marshal, _ := json.Marshal(vo)
	fmt.Println(string(marshal))
}
