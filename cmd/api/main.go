package main

import (
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/api/apis"
)

func main() {
	// 启动服务
	apiServer := api.NewBaseServer("0.0.0.0:8081")
	server := apis.NewServer(apiServer)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
