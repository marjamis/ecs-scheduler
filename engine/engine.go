package engine

import (
	"ecs-scheduler/action"
	"ecs-scheduler/schedulers"
	"ecs-scheduler/state"

	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"

	log "github.com/sirupsen/logrus"
)

var (
	debug          bool
	scheduler      string
	cluster        string
	region         string
	taskDefinition string
	statel         = "default"
)

func init() {
	//Flags
	flag.BoolVar(&debug, "debug", false, "Sets the debug level of output.")

	//Settings
	flag.StringVar(&scheduler, "scheduler", "leastTasks", "Use this for the LeastTasks running schedule.")
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
	return selectProcess(svc)
}

func selectProcess(svc ecsiface.ECSAPI) int {
	var instances *ecs.DescribeContainerInstancesOutput
	var err error
	switch statel {
	default:
		instances, err = state.DescribeContainerInstances(&cluster, svc)
		if err != nil {
			log.WithFields(log.Fields{
				"function": "engine.selectProcess",
			}).Error(err)
			return ExitStateError
		}
	}

	var instance *string
	switch scheduler {
	case "leastTasks":
		fallthrough
	default:
		instance = schedulers.LeastTasks(instances)
		if instance == nil {
			log.WithFields(log.Fields{
				"function": "engine.selectProcess",
			}).Error("No valid Container Instance returned to start task on")
			return ExitNoValidContainerInstance
		}
	}

	log.Warn(instance)
	startTaskError := action.StartTask(instance, &cluster, svc, taskDefinition)
	if startTaskError != nil {
		log.WithFields(log.Fields{
			"function": "engine.selectProcess",
		}).Error(startTaskError)
		return ExitStartTaskFailure
	}

	return ExitSuccess
}
