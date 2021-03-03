package main

import (
	"github.com/yametech/devops-cmdb-service/pkg/api"
	"github.com/yametech/devops-cmdb-service/pkg/api/web"
)

func main() {
	apiServer := api.NewBaseServer("0.0.0.0:8080")
	server := web.NewServer(apiServer)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
