package models

import (
	"encoding/json"
	"fmt"
)

// QueueScheduler :
type QueueScheduler struct {
	ID                    int64                  `json:"id"`
	Name                  string                 `json:"name"`
	URL                   string                 `json:"url"`
	URLHash               string                 `json:"url_hash"`
	Enabled               bool                   `json:"enabled"`
	ExternalID            string                 `json:"external_id"`
	Config                MapStringInterface     `json:"config"`
	LastScheduleAt        NullToEmptyString      `json:"last_schedule_at"`
	QueueSchedulerHistory *QueueSchedulerHistory `json:"queue_schedule_history"`
	QueueType             *QueueType             `json:"queue_type"`
	Company               *Company               `json:"company"`
}

// QueuesScheduler :
type QueuesScheduler []QueueScheduler

// QueueSchedulerAlias :
type QueueSchedulerAlias QueueScheduler

// UnmarshalJSON :
func (e *QueueScheduler) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	a := QueueSchedulerAlias{}

	var (
		id         float64
		externalID string
		name       string
		URL        string
		URLHash    string
		config     map[string]interface{}
		vqueueType map[string]interface{}
		vcompany   map[string]interface{}
		queueType  QueueType
		company    Company
		b          []byte
		err        error
	)

	// if nil != v["id"] {
	// 	id = v["id"].(float64)
	// }

	// if nil != v["externalId"] {
	// 	externalID = v["externalId"].(string)
	// }

	if nil != v["queue_type"] {
		vqueueType = v["queue_type"].(map[string]interface{})
	}

	if nil != v["company"] {
		vcompany = v["company"].(map[string]interface{})
	}

	id, ok := v["id"].(float64)
	externalID, ok = v["externalId"].(string)
	name, ok = v["name"].(string)
	URL, ok = v["url"].(string)
	URLHash, ok = v["url_hash"].(string)
	enabled := v["enabled"].(bool)
	vconfig, ok := v["config"].(string)

	if !ok {
		config = v["config"].(map[string]interface{})
	} else {
		if err := json.Unmarshal([]byte(vconfig), &config); err != nil {
			fmt.Printf("QueueScheduler UnmarshalJSON - Unmarshal %v\n", err)
			return err
		}
	}

	if b, err = json.Marshal(vqueueType); err != nil {
		fmt.Printf("Marshal QueueScheduler QueueType error: %v", err)
		return err
	}

	if err := json.Unmarshal(b, &queueType); err != nil {
		fmt.Printf("Unmarshal QueueScheduler QueueType error: %v", err)
		return err
	}

	if b, err = json.Marshal(vcompany); err != nil {
		fmt.Printf("Marshal QueueScheduler Company error: %v", err)
		return err
	}

	if err := json.Unmarshal(b, &company); err != nil {
		fmt.Printf("Unmarshal QueueScheduler Company error: %v", err)
		return err
	}

	a.ID = int64(id)
	a.Name = name
	a.URL = URL
	a.URLHash = URLHash
	a.Enabled = enabled
	a.ExternalID = externalID
	a.Config = config
	a.QueueType = &queueType
	a.Company = &company
	a.LastScheduleAt = NullToEmptyString("")

	*e = QueueScheduler(a)

	return nil
}
