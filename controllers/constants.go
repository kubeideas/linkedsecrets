package controllers

// Status constants
const (
	JOBSCHEDULED         = "Scheduled"
	JOBSUSPENDED         = "Suspended"
	JOBNOTSCHEDULED      = "NotScheduled"
	JOBFAILPARSESCHEDULE = "FailParseSchedule"
	STATUSSYNCHED        = "Synched"
	STATUSNOTSYNCHED     = "NotSynched"
	KEEPSECRETON         = "ON"
	KEEPSECRETOFF        = "OFF"
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
