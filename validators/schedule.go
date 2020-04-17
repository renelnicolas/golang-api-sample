package validators

import (
	"encoding/json"
	"errors"
	"fmt"

	"ohmytech.io/platform/models"
)

// ValidateScheduleQueue :
func ValidateScheduleQueue(filter models.QueryFilter, body []byte, validateContent bool) (_entity interface{}, err error) {
	extras, ok := filter.Extras.(models.Scheduler)
	if !ok {
		return nil, errors.New("Cannot convert to 'Scheduler'")
	}

	switch extras.Queue {
	case "parsing":
		fallthrough
	case "analyser":
		e := models.Analyser{
			TypeOf: extras.Queue,
		}

		err = json.Unmarshal(body, &e)
		if err != nil {
			return nil, err
		}

		err = validateAnalyser(e)
		if err != nil {
			return nil, err
		}

		return e, nil
	}

	return nil, fmt.Errorf("Cannot find Queue :%s", extras.Queue)
}

func validateAnalyser(entity models.Analyser) error {
	if "" == entity.URL {
		return errors.New("URL field cannot be empty")
	}

	return nil
}
