
clean:
	rm ./bin/ecs-scheduler.go

setup:
	GOPATH=$(shell pwd) go get "github.com/aws/aws-sdk-go/aws" "github.com/aws/aws-sdk-go/aws/session" "github.com/aws/aws-sdk-go/service/ecs" "github.com/Sirupsen/logrus"

build:
	GOPATH=$(shell pwd) CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/ecs-scheduler.go ecs-scheduler
