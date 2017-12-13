package main

const (
	//MaxResultsPerCall Specifies the number of items to be returned for paginating.
	MaxResultsPerCall = int64(100)

	//SchedulerName The name used as the TaskGroup to show what scheduled the Task.
	SchedulerName = string("svc.customECSscheduler")
)
