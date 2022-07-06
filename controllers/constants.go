package controllers

// Cronjob Status
const (
	JOBSCHEDULED         = "Scheduled"
	JOBSUSPENDED         = "Suspended"
	JOBNOTSCHEDULED      = "NotScheduled"
	JOBFAILPARSESCHEDULE = "FailParseSchedule"
)

// secret synch status
const (
	STATUSSYNCHED    = "Synched"
	STATUSNOTSYNCHED = "NotSynched"
)

// Cloud constants
const (
	GOOGLE = "Google"
	AWS    = "AWS"
	AZURE  = "Azure"
	IBM    = "IBM"
)

// Cloud Secret data format
const (
	JSONFORMAT  = "JSON"
	PLAINFORMAT = "PLAIN"
)

//LINKEDSECRETFINALIZER identify linkendsecret to be intercept before delete
const LINKEDSECRETFINALIZER = "cronjob.finalizers.linkedsecrets.kubeidea.io"
