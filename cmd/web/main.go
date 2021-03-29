package main

import (
	"flag"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/api/web"
	"github.com/yametech/devops-cmdb-service/pkg/common"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	"github.com/yametech/go_insect"
)

func main() {
	var host, username, password = "localhost", "neo4j", "123456"
	var etcdAddress, insectServerName = "http://10.200.65.207:2379", "cmdb"
	flag.StringVar(&host, "host", "localhost", "-host xxxx")
	flag.StringVar(&username, "username", "neo4j", "-username xxxx")
	flag.StringVar(&password, "password", "123456", "-password xxxx")
	flag.StringVar(&etcdAddress, "etcdAddress", etcdAddress, "-etcdAddress xxxx")
	flag.StringVar(&insectServerName, "insectServerName", insectServerName, "-insectServerName xxxx")
	flag.StringVar(&common.LdapAuthPassword, "LdapAuthPassword", "", "-LdapAuthPassword xxxx")
	flag.Parse()

	fmt.Println("Neo4jInit....start")
	store.Neo4jInit(host, username, password)
	fmt.Println("Neo4jInit....end")

	go_insect.GlobalEtcdAddress = etcdAddress // etcd host
	//go_insect.GlobalEtcdTTL = 60                          // ttl
	//go_insect.INSECT_SERVER_URL = "10.1.150.90"               // register server host
	//go_insect.INSECT_SERVER_URL = "10.1.170.82"               // register server host
	go_insect.INSECT_SERVER_PORT = 8080
	go_insect.INSECT_SERVER_NAME = insectServerName // register server name
	go go_insect.EtcdProxy()

	apiServer := api.NewBaseServer("0.0.0.0:8080")
	server := web.NewServer(apiServer)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
