package utils

import "github.com/mindstand/gogm"

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
