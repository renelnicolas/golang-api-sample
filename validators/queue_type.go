package validators

import (
	"encoding/json"
	"errors"

	"ohmytech.io/platform/models"
)

// ValidateQueueType : QueueType object validator from POST form
func ValidateQueueType(filter models.QueryFilter, body []byte) (entity *models.QueueType, err error) {
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return nil, err
	}

	err = validateQueueTypeFields(*entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// validateQueueTypeFields : check if all necessary field is defined and not empty
func validateQueueTypeFields(entity models.QueueType) error {
	if "" == entity.Name {
		return errors.New("Name field cannot be empty")
	}

	if 0 == len(entity.Config) {
		return errors.New("Config field cannot be empty")
	}

	return nil
}
