module github.com/yametech/devops-cmdb-service

go 1.16

require (
	github.com/adam-hanna/arrayOperations v0.2.6
	github.com/cornelk/hashmap v1.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/google/uuid v1.2.0
	github.com/mindstand/go-cypherdsl v0.2.0
	//github.com/mindstand/gogm v1.5.1
	github.com/neo4j/neo4j-go-driver/v4 v4.2.4
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.5.1
	github.com/yametech/go_insect v0.0.0-20210324065405-897e12c643e9
)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
