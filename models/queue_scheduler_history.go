package models

// QueueSchedulerHistory :
type QueueSchedulerHistory struct {
	ID                 int64             `json:"id"`
	UUID               string            `json:"uuid"`
	Status             int               `json:"status"`
	ScheduledStartDate NullToEmptyString `json:"scheduled_start_date"`
	ScheduledEndDate   NullToEmptyString `json:"scheduled_end_date"`
	QueueSchedulerID   int64             `json:"queue_scheduler_id"`
}

// QueuesSchedulerHistory :
type QueuesSchedulerHistory []QueueSchedulerHistory

// QueueSchedulerHistoryAlias :
type QueueSchedulerHistoryAlias QueueSchedulerHistory
