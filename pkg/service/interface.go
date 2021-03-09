package service

import (
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/core"
	"github.com/yametech/devops-cmdb-service/pkg/store"
)

type Service struct {
	store.IStore
}

type fakeService struct {
	Service
}

func (f *fakeService) GetMember(uuid string) core.IObject {
	obj, err := f.Get(uuid)
	if err != nil {
		//
	}
	return obj
}

func init() {
	fmt.Println("Neo4jInit....start")
	store.Neo4jInit("localhost", "neo4j", "123456")
	fmt.Println("Neo4jInit....end")
}

func (s Service) ManualQuery(query string, properties map[string]interface{}, respObj interface{}) {
	store.GetSession(true).Query(query, properties, respObj)
}

func (s Service) ManualExecute(query string, properties map[string]interface{}) ([][]interface{}, error) {
	return store.GetSession(false).QueryRaw(query, properties)
}
