package models

// Scheduler :
type Scheduler struct {
	ID    int64  `json:"id"`
	Queue string `json:"queue"`
}

// Analyser :
type Analyser struct {
	ID             int64             `json:"id"`
	Name           NullToEmptyString `json:"name"`
	URL            string            `json:"url"`
	URLHash        string            `json:"url_hash"`
	Enabled        bool              `json:"enabled"`
	TypeOf         string            `json:"type_of"`
	ExternalID     string            `json:"external_id"`
	LastScheduleAt NullToEmptyString `json:"last_schedule_at"`
	Company        *Company          `json:"company"`
}

// Analysers :
type Analysers []Analyser

// AnalyserAlias :
type AnalyserAlias Analyser
