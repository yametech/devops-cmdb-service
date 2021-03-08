package service

import (
	"fmt"
	"github.com/mindstand/gogm"
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

func getSession() *gogm.Session {
	//param is readonly, we're going to make stuff so we're going to do read write
	sess, err := gogm.NewSession(false)
	//sess, err := gogm.NewSessionWithConfig(gogm.SessionConfig{DatabaseName:"cmdb"})
	if err != nil {
		panic(err)
	}

	//close the session
	defer sess.Close()
	return sess
}

func (s Service) ManualQuery(query string, properties map[string]interface{}, respObj interface{}) {
	getSession().Query(query, properties, respObj)
}
