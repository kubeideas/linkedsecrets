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
	AZURE  = "Azure"
)

// Provider data format
const (
	JSONFORMAT  = "JSON"
	PLAINFORMAT = "PLAIN"
)
