package common

import (
	"flag"
	"fmt"
	"github.com/yametech/devops-cmdb-service/pkg/store"
	insect "github.com/yametech/go_insect"
	"log"
	"os"
)

func WebConfig() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	// 获取配置信息
	var host, username, password = "localhost", "neo4j", "123456"
	var etcdAddress, insectServerName = "http://0.0.0.0:2379", "cmdb"
	flag.StringVar(&host, "host", "10.200.10.51", "-host xxxx")
	flag.StringVar(&username, "username", "neo4j", "-username xxxx")
	flag.StringVar(&password, "password", "test123qwe", "-password xxxx")
	flag.StringVar(&etcdAddress, "etcdAddress", etcdAddress, "-etcdAddress xxxx")
	flag.StringVar(&insectServerName, "insectServerName", insectServerName, "-insectServerName xxxx")
	flag.StringVar(&LdapAuthPassword, "LdapAuthPassword", "", "-LdapAuthPassword xxxx")
	flag.Parse()

	// 数据库连接初始化
	fmt.Println("Neo4jInit....start")
	store.Neo4jInit(host, username, password)
	fmt.Println("Neo4jInit....end")

	// 接入网关
	insect.GlobalEtcdAddress = etcdAddress // etcd host
	//insect.GlobalEtcdTTL = 60                          // ttl
	//insect.INSECT_SERVER_URL = "10.1.150.90"               // register server host
	//insect.INSECT_SERVER_URL = "10.1.170.82"               // register server host
	insect.INSECT_SERVER_PORT = 8080
	insect.INSECT_SERVER_NAME = insectServerName // register server name
	go insect.EtcdProxy()
}

func ApiConfig() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	// 获取配置信息
	var host, username, password = "localhost", "neo4j", "123456"
	if os.Getenv("host") != "" {
		host = os.Getenv("host")
	}
	if os.Getenv("account") != "" {
		username = os.Getenv("account")
	}
	if os.Getenv("password") != "" {
		password = os.Getenv("password")
	}

	// 数据库连接初始化
	fmt.Println("Neo4jInit....start")
	store.Neo4jInit(host, username, password)
	fmt.Println("Neo4jInit....end")
}
