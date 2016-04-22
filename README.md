# ECS Scheduler

## To run:
- Clone repository.
- Get required additional repositories with:
	go get "github.com/aws/aws-sdk-go/aws"
	go get "github.com/aws/aws-sdk-go/aws/session"
	go get "github.com/aws/aws-sdk-go/service/ecs"
	go get "github.com/Sirupsen/logrus"
- Run the application. Basic version such as:
go run src/ecs-scheduler/scheduler.go --task-definition <task_definition> --region <region>  --cluster <cluster_name>
