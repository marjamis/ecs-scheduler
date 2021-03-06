package schedulers

import (
	"github.com/aws/aws-sdk-go/service/ecs"

	log "github.com/sirupsen/logrus"
)

//LeastTasks Determines which Container Instance is currently running the least number of Tasks
func LeastTasks(instances *ecs.DescribeContainerInstancesOutput) *string {
	log.WithFields(log.Fields{
		"function": "schedulers.LeastTasks",
	}).Info("Starting")
	var instancesSlice []*ecs.ContainerInstance
	if instances != nil {
		instancesSlice = instances.ContainerInstances
	} else {
		return nil
	}

	var selected *ecs.ContainerInstance
	if instancesSlice != nil {
		selected = instancesSlice[0]
	} else {
		return nil
	}

	for _, each := range instancesSlice {
		log.WithFields(log.Fields{
			"function":     "schedulers.LeastTasks",
			"arn":          *each.ContainerInstanceArn,
			"runningTasks": *each.RunningTasksCount,
			"pendingTasks": *each.PendingTasksCount,
		}).Debug("Checking details of Container Instance")
		if (*each.RunningTasksCount + *each.PendingTasksCount) < (*selected.RunningTasksCount + +*selected.PendingTasksCount) {
			selected = each
		}
	}

	log.WithFields(log.Fields{
		"function": "schedulers.LeastTasks",
		"arn":      *selected.ContainerInstanceArn,
	}).Info("Container Instance selected")
	log.Info(*selected.ContainerInstanceArn)
	return selected.ContainerInstanceArn
}
