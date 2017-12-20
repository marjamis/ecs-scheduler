package state

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"

	log "github.com/sirupsen/logrus"
)

//Returns a []*string for the list of all the Container Instance ARNs for a given cluster
func getInstanceARNs(clusterName *string, svc ecsiface.ECSAPI) (output []*string, erro error) {
	log.WithFields(log.Fields{
		"function": "getInstanceARNs",
	}).Info("start")

	var ciARNs []*string
	ciARNs = make([]*string, 0)

	nextToken := aws.String("nextToken")
	params := &ecs.ListContainerInstancesInput{
		Cluster:    aws.String(*clusterName),
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
				Cluster:    aws.String(*clusterName),
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

//DescribeContainerInstances Returns the details of all the Container Instances in the given cluster
func DescribeContainerInstances(clusterName *string, svc ecsiface.ECSAPI) (output *ecs.DescribeContainerInstancesOutput, erro error) {
	log.Info("Function: describeContainerInstances - getInstanceARNs")
	instanceARNs, err := getInstanceARNs(clusterName, svc)
	if err != nil {
		return nil, err
	}

	log.Info("Function: describeContainerInstances")
	params := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(*clusterName),
		ContainerInstances: instanceARNs,
	}
	resp, err := svc.DescribeContainerInstances(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
