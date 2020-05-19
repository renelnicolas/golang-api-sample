package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// QueueType :
type QueueType struct {
	ID         int64              `json:"id"`
	Name       string             `json:"name"`
	Config     MapStringInterface `json:"config"`
	Enabled    bool               `json:"enabled"`
	ExternalID string             `json:"external_id"`
}

// QueuesType :
type QueuesType []QueueType

// QueueTypeAlias :
type QueueTypeAlias QueueType

// UnmarshalJSON :
func (e *QueueType) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	a := QueueTypeAlias{}

	var (
		id         float64
		externalID string
		name       string
		config     map[string]interface{}
		err        error
		ok         bool
	)

	if nil != v["id"] {
		id, ok = v["id"].(float64)

		if !ok {
			vid := v["id"].(string)

			if id, err = strconv.ParseFloat(vid, 64); err != nil {
				fmt.Printf("QueueType UnmarshalJSON - Unmarshal %v\n", err)
				return err
			}
		}
	}

	if nil != v["externalId"] {
		externalID = v["externalId"].(string)
	}

	if nil != v["name"] {
		name = v["name"].(string)
	}

	if nil != v["config"] {
		vconfig, ok := v["config"].(string)

		if !ok {
			config = v["config"].(map[string]interface{})
		} else {
			if "" != vconfig {
				if err := json.Unmarshal([]byte(vconfig), &config); err != nil {
					fmt.Printf("QueueType UnmarshalJSON - Unmarshal %v\n", err)
					return err
				}
			}
		}
	}

	enabled := v["enabled"].(bool)

	a.ID = int64(id)
	a.Name = name
	a.Enabled = enabled
	a.ExternalID = externalID
	a.Config = config

	*e = QueueType(a)

	return nil
}
