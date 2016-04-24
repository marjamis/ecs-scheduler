# ECS Scheduler

## To run:
1. Clone repository.
1. Get required additional repositories with:
go get "github.com/aws/aws-sdk-go/aws" "github.com/aws/aws-sdk-go/aws/session" "github.com/aws/aws-sdk-go/service/ecs" "github.com/Sirupsen/logrus"
1. Run the application. Basic version such as:
	go run src/ecs-scheduler/scheduler.go --task-definition task_definition --region region  --cluster cluster_name
