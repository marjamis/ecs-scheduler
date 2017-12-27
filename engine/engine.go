package engine

import (
	"ecs-scheduler/action"
	"ecs-scheduler/schedulers"
	"ecs-scheduler/state"

	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"

	log "github.com/sirupsen/logrus"
)

var (
	debug           bool
	leastTasksSched bool
	cluster         string
	region          string
	taskDefinition  string
)

func init() {
	//Flags
	flag.BoolVar(&debug, "debug", false, "Sets the debug level of output.")
	flag.BoolVar(&leastTasksSched, "leastTasks", true, "Use this for the LeastTasks running schedule.")

	//Settings
	flag.StringVar(&cluster, "cluster", "", "Name of the cluster to schedule against.")
	flag.StringVar(&region, "region", "", "Region that the cluster is in.")
	flag.StringVar(&taskDefinition, "task-definition", "", "The Task Definition to be used when scheduling the Task.")

	flag.Parse()
}

//ConnectToECS Simple function to set ECS Client settings, largely incase I change/need new clients.
func ConnectToECS(region *string) (svc *ecs.ECS) {
	return ecs.New(session.New(), aws.NewConfig().WithRegion(*region))
}

//Run will run the engines CLI component
func Run() int {
	if cluster == "" || region == "" || taskDefinition == "" {
		log.WithFields(log.Fields{
			"function": "engine.Run",
		}).Error("Error: Insufficient command-line options have been supplied. Use --help to see the required options.")
		return ExitInvalidCLIOptions
	}

	if debug == true {
		log.SetLevel(log.DebugLevel)
	}

	svc := ConnectToECS(&region)

	instances, err := state.DescribeContainerInstances(&cluster, svc)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "engine.Run",
		}).Error(err)
		return ExitStateError
	}

	instance := schedulers.LeastTasks(instances)
	if &instance == nil {
		log.WithFields(log.Fields{
			"function": "engine.Run",
		}).Error("No valid Container Instance returned to start task on")
		return ExitNoValidContainerInstance
	}

	//Selection of which scheduler to be used based off the flag that was passed in. Default is leastTasks.
	var runTaskError error
	switch {
	//Room to move to add additional schedules in the future.
	case leastTasksSched == true:
		runTaskError = action.StartTask(instance, &cluster, svc, taskDefinition)
	}

	if runTaskError != nil {
		log.WithFields(log.Fields{
			"function": "engine.Run",
		}).Error(runTaskError)
		return ExitStartTaskFailure
	}

	return ExitSuccess
}
