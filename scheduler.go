package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	log "github.com/sirupsen/logrus"
)

//Returns a []*string for the list of all the Container Instance ARNs for a given cluster
func getInstanceARNs(clusterName string, svc ecsiface.ECSAPI) (output []*string, erro error) {
	log.WithFields(log.Fields{
		"function": "getInstanceARNs",
	}).Info("start")

	var ciARNs []*string
	ciARNs = make([]*string, 0)

	nextToken := aws.String("nextToken")
	params := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int64(MaxResultsPerCall),
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
				MaxResults: aws.Int64(MaxResultsPerCall),
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
func describeContainerInstances(clusterName string, svc ecsiface.ECSAPI) (output *ecs.DescribeContainerInstancesOutput, erro error) {
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
		StartedBy:          aws.String(SchedulerName),
	}

	resp, err := svc.StartTask(params)
	if err != nil {
		return err
	}

	if len(resp.Failures) != 0 {
		log.WithFields(log.Fields{
			"response": resp,
		}).Error("Response from creating Task")
		return errors.New("Failures listed in response")
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
