all: build-web

build-web:
	docker build -t harbor.ym/devops/cmdbweb:v0.1.0 -f docker/Dockerfile.web .
	docker push harbor.ym/devops/cmdbweb:v0.1.0
