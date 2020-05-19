package validators

import (
	"encoding/json"
	"errors"

	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
)

// ValidateQueueScheduler : QueueScheduler object validator from POST form
func ValidateQueueScheduler(filter models.QueryFilter, body []byte) (entity *models.QueueScheduler, err error) {
	err = json.Unmarshal(body, &entity)
	if err != nil {
		return nil, err
	}

	err = validateQueueSchedulerFields(*entity)
	if err != nil {
		return nil, err
	}

	err = validateQueueType(entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// validateQueueSchedulerFields : check if all necessary field is defined and not empty
func validateQueueSchedulerFields(entity models.QueueScheduler) error {
	if "" == entity.URL {
		return errors.New("URL field cannot be empty")
	}

	if nil == entity.QueueType {
		return errors.New("QueueType field cannot be null")
	}

	if 0 == len(entity.Config) {
		return errors.New("Config field cannot be empty")
	}

	return nil
}

// validateQueueType: check if current QueueType exist
func validateQueueType(entity *models.QueueScheduler) error {
	repo := repositories.QueueTypeRepository{}

	r, err := repo.FindOne(*entity.QueueType, models.QueryFilter{})
	if err != nil {
		return err
	}

	queueType := r.(models.QueueType)

	entity.QueueType = &queueType

	return nil
}
