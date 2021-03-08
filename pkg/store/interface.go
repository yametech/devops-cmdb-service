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
