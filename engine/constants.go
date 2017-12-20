package engine

const (
	//ExitSuccess A Task was successfully started
	ExitSuccess = 0

	//ExitInvalidCLIOptions Arguments for the CLI were missing
	ExitInvalidCLIOptions = 1

	//ExitStateError Failure with an API call in determining state
	ExitStateError = 2

	//ExitNoValidContainerInstance No valid Container Instance was returned as a placement instance
	ExitNoValidContainerInstance = 3

	//ExitStartTaskFailure StartTask API call failed
	ExitStartTaskFailure = 4
)
