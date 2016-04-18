package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"os"
)

//Specifies the number of items to be returned for paginating.
const numberPerCall = int64(10)

//Returns a []*string for the list of all the Container Instance ARNs for a given cluster
func getInstanceARNs(clusterName string, region string) (output []*string, output_err string) {
	fmt.Println("Function: getInstanceARNs")
	//Change from number array to allow dynamic growing
	var ciARNs []*string
	ciARNs = make([]*string, 0)
	svc := ecs.New(session.New(), aws.NewConfig().WithRegion(region))
	check := aws.String("checkInitial")
	params := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(clusterName),
		MaxResults: aws.Int64(numberPerCall),
	}

	for ok := true; ok; ok = (check != nil) {
		resp, err := svc.ListContainerInstances(params)
		if err != nil {
			return nil, "Error: FIX THESE to be errors not strings" //exit
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
				MaxResults: aws.Int64(numberPerCall),
				NextToken:  aws.String(*resp.NextToken),
			}
			check = resp.NextToken
		} else {
			check = nil
		}
	}
	return ciARNs, ""
}

//Returns the details of all the Container Instances in the given cluster
func describeContainerInstances(clusterName string, region string) (output *ecs.DescribeContainerInstancesOutput, output_err string) {
	fmt.Println("Function: describeContainerInstances - getInstanceARNs")
	instanceARNs, errC := getInstanceARNs(clusterName, region)
	if errC != "" {
		return nil, errC
	} else if len(instanceARNs) == 0 {
		return nil, "Check: No instanceARNs"
	}

	fmt.Println("Function: describeContainerInstances")
	svc := ecs.New(session.New(), aws.NewConfig().WithRegion(region))
	params := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(clusterName),
		ContainerInstances: instanceARNs,
	}
	resp, err := svc.DescribeContainerInstances(params)
	if err != nil {
		return nil, "Error: DescribeContainerInstances"
	}

	return resp, ""
}

//Starts the task on the given Container Instance
func startTask(instanceARN string, cluster string, taskDefinition string, region string) {
	fmt.Printf("Function: startTask: %s, %s\n", instanceARN, taskDefinition)
	//take this out to minimise the number of times this is done as its nor required
	svc := ecs.New(session.New(), aws.NewConfig().WithRegion(region))
	var containers []*string
	containers = make([]*string, 1)
	containers[0] = aws.String(instanceARN)
	params := &ecs.StartTaskInput{
		ContainerInstances: containers,
		TaskDefinition: aws.String(taskDefinition),
		Cluster: aws.String(cluster),
		//StartedBy: aws.String("svc.customECSscheduler")
	}

	resp, err := svc.StartTask(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}

//Determines which Container Instance is currently running the least number of Tasks
func leastTasks(instances *ecs.DescribeContainerInstancesOutput) (instanceARN string, output_err string) {
	fmt.Println("Function: leastTasks\n\n")
	instancesSlice := instances.ContainerInstances
	//sort.Sort(instances.ContainerInstances.RunningTasksCount)

	selected := instancesSlice[0]
	for _, each := range instancesSlice {
		fmt.Printf("Output: ContainerInstanceArn: %s RunningTasksCount: %d\n",
			*each.ContainerInstanceArn, *each.RunningTasksCount)
		if *each.RunningTasksCount < *selected.RunningTasksCount {
			selected = each
		}
	}
    
    //Add better error supporting
	fmt.Printf("\n\nContainer Instance selected: %s\n", *selected.ContainerInstanceArn)
	return *selected.ContainerInstanceArn, "nil"
}

func main() {
	fmt.Println("Starting scheduler...")

    //Add check for commandline arguments
	instances, err := describeContainerInstances(os.Args[1], os.Args[2])
	if err != "" {
		fmt.Println(err)
		os.Exit(1)
	}
	
	instance, err := leastTasks(instances)
	if err != "nil" {
		fmt.Println(err)
		os.Exit(1)
	}
    fmt.Println(instance)
    //for a display method add a show of pending(maybe even to leastTasks to ensure Pendinsgs are counted on display but this is a basic test not sure I care that much)
	startTask(instance, os.Args[1], os.Args[3], os.Args[2])

	fmt.Println("Exiting scheduler...")
}
