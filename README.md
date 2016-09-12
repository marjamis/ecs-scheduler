# ECS Scheduler

## To build and run:
1. Clone repository.
1. Get required additional libraries with:
	make setup
1. Run the application. Basic version such as:
	go run src/ecs-scheduler/scheduler.go --task-definition <task_definition> --region <region> --cluster <cluster_name>
   and if you're happy with the current state compile a statically compiled version with:
	make build
