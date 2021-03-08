package utils

import (
//dsl "github.com/mindstand/go-cypherdsl"
//"github.com/mindstand/gogm"
//"reflect"
)

//func GetSession() *gogm.Session {
//	//param is readonly, we're going to make stuff so we're going to do read write
//	sess, err := gogm.NewSession(false)
//	//sess, err := gogm.NewSessionWithConfig(gogm.SessionConfig{DatabaseName:"cmdb"})
//	if err != nil {
//		panic(err)
//	}
//
//	//close the session
//	defer sess.Close()
//
//	return sess
//}

//func Neo4jInit(host string, username string, password string)  {
//	config := &gogm.Config{
//		IndexStrategy: gogm.VALIDATE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
//		PoolSize:      50,
//		Port:          7687,
//		IsCluster:     false, //tells it whether or not to use `bolt+routing`
//		Host:          host,
//		Username:  username,
//		Password:  password,
//	}
//
//	err := gogm.Init(config, &store.ModelGroup{}, &store.Model{}, &store.AttributeGroup{}, &store.Attribute{})
//	if err != nil {
//		panic(err)
//	}
//}

func cypherMachine(obj interface{}) {

	//t := reflect.TypeOf(obj)
	//path := dsl.Path().
	//	P().
	//	V(dsl.V{Name: "n"})
	//builder := dsl.QB().
	//	Match(path.Build())
	//cyp, _ = dsl.Path().V(dsl.V{Name: "n", Type: t.Name()}).ToCypher()
}
