.PHONY: docker
# make docker 启动
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rm -f d17/webook:v0.0.1
	@docker build -t d17/webook:v0.0.1 .