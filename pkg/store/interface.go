package store

import (
	"fmt"
	dsl "github.com/mindstand/go-cypherdsl"
	"github.com/mindstand/gogm"
	"github.com/yametech/devops-cmdb-service/pkg/core"
	"reflect"
	"time"
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
		PoolSize:      500,
		Port:          7687,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Host:          host,
		Username:      username,
		Password:      password,
	}

	err := gogm.Init(config,
		&ModelGroup{}, &Model{}, &AttributeGroup{}, &Attribute{}, &ModelRelation{},
		&Resource{}, &AttributeGroupIns{}, &AttributeIns{}, &RelationshipModel{},
	)
	if err != nil {
		panic(err)
	}
}

func GetSession(readonly bool) *gogm.Session {
	//param is readonly, we're going to make stuff so we're going to do read write
	sess, err := gogm.NewSession(readonly)
	//sess, err := gogm.NewSessionWithConfig(gogm.SessionConfig{DatabaseName:"cmdb"})
	if err != nil {
		panic(err)
	}

	//close the session
	defer func() {
		start := time.Now()
		if e := sess.Close(); e != nil {
			fmt.Println("close the session err：", e)
		}
		end := time.Now()
		latency := end.Sub(start)
		fmt.Println("close the session finished cost:", latency)
	}()

	return sess
}

type INeo4j interface {
	Get(string) (interface{}, error)
	List(string) ([]interface{}, error)
	Save(interface{}) error
	Update(interface{}) error
	Delete(interface{}) error
}

type Neo4jDomain struct {
	// neo4j node id
	Id string
	// gogm中间件主键
	Uuid string
	// cmdb主键
	Uid string
}

func (domain *Neo4jDomain) Get(respObj interface{}, key string, value interface{}) error {
	rft := reflect.TypeOf(respObj)
	params, _ := dsl.ParamsFromMap(map[string]interface{}{key: value})
	cypher, _ := dsl.QB().
		Match(dsl.Path().V(dsl.V{Name: "a", Type: rft.Elem().Name(), Params: params}).Build()).
		Return(false, dsl.ReturnPart{Name: "a"}).
		ToCypher()

	fmt.Println(cypher)
	return GetSession(true).Query(cypher, nil, respObj)
}

func (domain *Neo4jDomain) List(respObj interface{}) error {
	return GetSession(true).LoadAll(respObj)
}

func (domain *Neo4jDomain) Save(respObj interface{}) error {
	return GetSession(false).Save(respObj)
}

func (domain *Neo4jDomain) Update(respObj interface{}) error {
	return GetSession(false).Save(respObj)
}

func (domain *Neo4jDomain) Delete(respObj interface{}) error {
	return GetSession(false).Delete(respObj)
}
