package store

import (
	"fmt"
	dsl "github.com/mindstand/go-cypherdsl"
	"github.com/yametech/devops-cmdb-service/pkg/core"
	"github.com/yametech/devops-cmdb-service/pkg/gogm"
	"reflect"
)

type IStore interface {
	List(string) ([]core.IObject, error)
	DeepSearch(obj core.IObject, edge string) ([]core.IObject, error)
	Get(string) (core.IObject, error)
	Put(core.IObject) error
	Update(core.IObject) error
	Delete(core.IObject) error
}

func Neo4jInit(neo4jUri, neo4jUsername, neo4jPassword string) {
	//return &Neo4jDomain{Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))}
	config := &gogm.Config{
		IndexStrategy: gogm.VALIDATE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:      200,
		Port:          7687,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Host:          neo4jUri,
		Username:      neo4jUsername,
		Password:      neo4jPassword,
		LogLevel:      "DEBUG",
	}

	err := gogm.Init(config,
		&ModelGroup{}, &Model{}, &AttributeGroup{}, &Attribute{}, &ModelRelation{},
		&Resource{}, &AttributeGroupIns{}, &AttributeIns{}, &RelationshipModel{},
	)
	if err != nil {
		panic(err)
	}
}

type INeo4j interface {
	Get(string) (interface{}, error)
	List(string) ([]interface{}, error)
	Save(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}

type Neo4jDomain struct {
	//read  *utils.GenericPool
	//write *utils.GenericPool
}

func (domain *Neo4jDomain) GetSession(readonly bool) *gogm.Session {
	session, err := gogm.NewSession(readonly)
	if err != nil {
		panic(err)
	}
	return session
}

func (domain *Neo4jDomain) Get(respObj interface{}, key string, value interface{}) error {
	rft := reflect.TypeOf(respObj)
	params, _ := dsl.ParamsFromMap(map[string]interface{}{key: value})
	cypher, _ := dsl.QB().
		Match(dsl.Path().V(dsl.V{Name: "a", Type: rft.Elem().Name(), Params: params}).Build()).
		Return(false, dsl.ReturnPart{Name: "a"}).
		ToCypher()

	fmt.Println(cypher)
	session := domain.GetSession(true)
	defer session.Close()
	return session.Query(cypher, nil, respObj)
}

func (domain *Neo4jDomain) List(respObj interface{}) error {
	session := domain.GetSession(true)
	defer session.Close()
	return session.LoadAll(respObj)
}

func (domain *Neo4jDomain) Save(respObj interface{}) error {
	session := domain.GetSession(false)
	defer session.Close()
	return session.Save(respObj)
}

func (domain *Neo4jDomain) Update(respObj interface{}) error {
	session := domain.GetSession(false)
	defer session.Close()
	return session.Save(respObj)
}

func (domain *Neo4jDomain) Delete(respObj interface{}) error {
	session := domain.GetSession(false)
	defer session.Close()
	return session.Delete(respObj)
}
