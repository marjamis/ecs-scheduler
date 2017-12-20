package action

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"

	log "github.com/sirupsen/logrus"
)

//StartTask Starts the task on the given Container Instance
func StartTask(instanceARN *string, cluster *string, svc ecsiface.ECSAPI, taskDefinition string) (erro error) {
	log.WithFields(log.Fields{
		"arn":             instanceARN,
		"task-definition": taskDefinition,
	}).Info("Function: startTask")

	var containers []*string
	containers = make([]*string, 1)
	containers[0] = aws.String(*instanceARN)

	params := &ecs.StartTaskInput{
		ContainerInstances: containers,
		TaskDefinition:     aws.String(taskDefinition),
		Cluster:            aws.String(*cluster),
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
