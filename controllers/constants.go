package controllers

// Status constants
const (
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
	AZURE  = "Azure"
	IBM    = "IBM"
)

// Provider data format
const (
	JSONFORMAT  = "JSON"
	PLAINFORMAT = "PLAIN"
)
