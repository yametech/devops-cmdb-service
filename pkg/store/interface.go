package store

import (
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/core"
)

type IStore interface {
	List(string) ([]core.IObject, error)
	DeepSearch(obj core.IObject, edge string) ([]core.IObject, error)
	Get(string) (core.IObject, error)
	Put(core.IObject) error
	Update(core.IObject) error
	Delete(core.IObject) error
}

func Neo4jInit(host string, username string, password string) {
	config := &gogm.Config{
		IndexStrategy: gogm.VALIDATE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:      50,
		Port:          7687,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Host:          host,
		Username:      username,
		Password:      password,
	}

	err := gogm.Init(config, &ModelGroup{}, &Model{}, &AttributeGroup{}, &Attribute{})
	if err != nil {
		panic(err)
	}
}

func GetSession() *gogm.Session {
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
