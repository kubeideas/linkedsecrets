package controllers

// Status constants
const (
	JOBACTIVE            = "Active"
	JOBSCHEDULED         = "Scheduled"
	JOBSUSPENDED         = "Suspended"
	JOBNOTSCHEDULED      = "NotScheduled"
	JOBFAILPARSESCHEDULE = "FailParseSchedule"
	STATUSSYNCHED        = "Synched"
	STATUSNOTSYNCHED     = "NotSynched"
)

// Providers constants
const (
	GOOGLE = "Google"
	AWS    = "AWS"
)

// Provider data format
const (
	JSONFORMAT  = "JSON"
	PLAINFORMAT = "PLAIN"
)
