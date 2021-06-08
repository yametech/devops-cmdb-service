#all: build-web build-api

build-web:
	docker build -t harbor.ym/devops/cmdbweb:v0.1.0 -f docker/Dockerfile.web .
	docker push harbor.ym/devops/cmdbweb:v0.1.0


build-api:
	docker build -t harbor.ym/devops/cmdb-api:v0.1.0 -f docker/Dockerfile.api .
	docker push harbor.ym/devops/cmdb-api:v0.1.0