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
	//var host, username, password string
	//flag.StringVar(&host, "host", "localhost", "-host xxxx")
	//flag.StringVar(&username, "username", "neo4j", "-username xxxx")
	//flag.StringVar(&password, "password", "123456", "-password xxxx")
	//flag.Parse()
	fmt.Println("Neo4jInit....start")
	store.Neo4jInit("localhost", "neo4j", "123456")
	fmt.Println("Neo4jInit....end")
}

func (s *Service) ManualQuery(query string, properties map[string]interface{}, respObj interface{}) error {
	return store.GetSession(true).Query(query, properties, respObj)
}

func (s *Service) ManualQueryRaw(query string, properties map[string]interface{}) ([][]interface{}, error) {
	return store.GetSession(true).QueryRaw(query, properties)
}

func (s *Service) ManualExecute(query string, properties map[string]interface{}) ([][]interface{}, error) {
	return store.GetSession(false).QueryRaw(query, properties)
}
