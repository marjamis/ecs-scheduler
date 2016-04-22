package main

import (
	"errors"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"os"
)

//Specifies the number of items to be returned for paginating.
const MAX_RESULTS_PER_CALL = int64(10)
const SCHEDULER_NAME = string("svc.customECSscheduler")

//Returns a []*string for the list of all the Container Instance ARNs for a given cluster
func getInstanceARNs(clusterName string, svc *ecs.ECS) (output []*string, erro error) {
	log.Info("Function: getInstanceARNs")

	var ciARNs []*string
	ciARNs = make([]*string, 0)

	nextToken := aws.String("nextToken")
	params := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int64(MAX_RESULTS_PER_CALL),
	}

	for ok := true; ok; ok = (nextToken != nil) {
		resp, err := svc.ListContainerInstances(params)
		if err != nil {
			return nil, err
		} else if len(resp.ContainerInstanceArns) == 0 {
			return nil, errors.New("Function: getInstanceARNs: No Container Instances in Cluster")
		}

		var tempciARNs []*string
		tempciARNs = make([]*string, len(resp.ContainerInstanceArns))
		count := 0
		for _, each := range resp.ContainerInstanceArns {
			tempciARNs[count] = aws.String(*each)
			count = count + 1
		}

		ciARNs = append(ciARNs, tempciARNs...)

		if resp.NextToken != nil {
			params = &ecs.ListContainerInstancesInput{
				Cluster:    aws.String(clusterName),
				MaxResults: aws.Int64(MAX_RESULTS_PER_CALL),
				NextToken:  aws.String(*resp.NextToken),
			}
			nextToken = resp.NextToken
		} else {
			nextToken = nil
		}
	}
	return ciARNs, nil
}

//Returns the details of all the Container Instances in the given cluster
func describeContainerInstances(clusterName string, svc *ecs.ECS) (output *ecs.DescribeContainerInstancesOutput, erro error) {
	log.Info("Function: describeContainerInstances - getInstanceARNs")
	instanceARNs, err := getInstanceARNs(clusterName, svc)
	if err != nil {
		return nil, err
	}

	log.Info("Function: describeContainerInstances")
	params := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: instanceARNs,
	}
	resp, err := svc.DescribeContainerInstances(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//Starts the task on the given Container Instance
func startTask(instanceARN string, cluster string, svc *ecs.ECS, taskDefinition string) (erro error) {
	log.WithFields(log.Fields{
		"arn":             instanceARN,
		"task-definition": taskDefinition,
	}).Info("Function: startTask")

	var containers []*string
	containers = make([]*string, 1)
	containers[0] = aws.String(instanceARN)

	params := &ecs.StartTaskInput{
		ContainerInstances: containers,
		TaskDefinition:     aws.String(taskDefinition),
		Cluster:            aws.String(cluster),
		StartedBy:          aws.String(SCHEDULER_NAME),
	}

	resp, err := svc.StartTask(params)
	if err != nil {
		return err
	}

	if len(resp.Failures) != 0 {
		log.WithFields(log.Fields{
			"response": resp,
		}).Error("Response from creating Task")
		return errors.New("Failures listed in response.")
	}

	log.WithFields(log.Fields{
		"response": resp,
	}).Info("Response from creating Task")

	return nil
}

//Determines which Container Instance is currently running the least number of Tasks
func leastTasks(instances *ecs.DescribeContainerInstancesOutput) (instanceARN string, erro error) {
	log.Info("Function: leastTasks")
	instancesSlice := instances.ContainerInstances

	selected := instancesSlice[0]
	for _, each := range instancesSlice {
		log.WithFields(log.Fields{
			"arn":          *each.ContainerInstanceArn,
			"runningTasks": *each.RunningTasksCount,
			"pendingTasks": *each.PendingTasksCount,
		}).Debug("Checking details of Container Instance")
		if (*each.RunningTasksCount + *each.PendingTasksCount) < *selected.RunningTasksCount {
			selected = each
		}
	}

	log.WithFields(log.Fields{
		"arn": *selected.ContainerInstanceArn,
	}).Info("Container Instance selected")
	return *selected.ContainerInstanceArn, nil
}

func connectToECS(region string) (svc *ecs.ECS) {
	return ecs.New(session.New(), aws.NewConfig().WithRegion(region))
}

func main() {
	//Flags
	// - List of all the flags for which scheduler that is to be used.
	leastTasksSched := flag.Bool("leastTasks", true, "Use this for the LeastTasks running schedule.")

	//Flags for how the application runs
	debug := flag.Bool("debug", false, "Sets the debug level of output.")

	//Settings
	cluster := flag.String("cluster", "", "Name of the cluster to schedule against.")
	region := flag.String("region", "", "Region that the cluster is in.")
	taskDefinition := flag.String("task-definition", "", "The Task Definition to be used when scheduling the Task.")

	flag.Parse()

	if *cluster == "" || *region == "" || *taskDefinition == "" {
		log.Error("Error: Insufficient command-line options haven't been supplied. Use --help to see required options.")
		os.Exit(1)
	}

	if *debug == true {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting scheduler...")

	svc := connectToECS(*region)

	instances, err := describeContainerInstances(*cluster, svc)
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

	instance, err := leastTasks(instances)
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}

	//Selection of which scheduler to be used based off the flag that was passed in. Default is leastTasks.
	var runTaskError error
	switch {
	//Room to move to add additional schedules in the future.
	case *leastTasksSched == true:
		runTaskError = startTask(instance, *cluster, svc, *taskDefinition)
	}

	if runTaskError != nil {
		log.Error(runTaskError)
		os.Exit(4)
	}

	log.Info("Exitiing scheduler...")
}
