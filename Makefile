GOOS=$(shell go env | grep GOOS | awk -F "=" '{print $$NF}' | awk -F "\"" '{print $$2}')
GOARCH=$(shell go env | grep GOARCH | awk -F "=" '{print $$NF}' | awk -F "\"" '{print $$2}')
NAME=$(shell pwd | awk -F "/" '{print $$NF}')
SERVICE=bin/$(NAME)
IMAGE=registry.cn-hangzhou.aliyuncs.com/qingzhen/$(NAME)
TAG?=latest
ENTRY=main.go

default: build
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(SERVICE) $(ENTRY)
linux:
	GOOS=linux GOARCH=amd64  go build -o $(SERVICE)-linux $(ENTRY)
darwin:
	GOOS=darwin GOARCH=amd64  go build -o $(SERVICE) $(ENTRY)
windows:
	GOOS=windows GOARCH=amd64  go build -o $(SERVICE).exe $(ENTRY)
image: linux
	# docker build --build-arg http_proxy=http://192.168.50.162:1087 --build-arg https_proxy=192.168.50.162:1087 --no-cache -t $(IMAGE):$(TAG) .
	docker build --no-cache -t $(IMAGE):$(TAG) .
push: image
	docker push $(IMAGE):$(TAG)
clean:
	rm -rf $(SERVICE)
run: build
	./$(SERVICE)