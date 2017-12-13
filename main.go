package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

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
		log.Error("Error: Insufficient command-line options have been supplied. Use --help to see the required options.")
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
