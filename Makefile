clean:
	rm ./bin/ecs-scheduler.go

setup:
	go get "github.com/aws/aws-sdk-go/aws" "github.com/aws/aws-sdk-go/aws/session" "github.com/aws/aws-sdk-go/service/ecs" "github.com/Sirupsen/logrus"

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/ecs-scheduler.go ecs-scheduler
