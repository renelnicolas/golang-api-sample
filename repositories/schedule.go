package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"ohmytech.io/platform/helpers"
	"ohmytech.io/platform/models"
)

// ScheduleRepository :
type ScheduleRepository struct {
	// analyser.id, analyser.company_id, analyser.name, analyser.url, analyser.url_hash, analyser.enabled, analyser.external_id, analyser.last_schedule_at
}

// Count :
func (r ScheduleRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count int64
	)

	return count, nil
}

// FindAll :
func (r ScheduleRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	// TODO

	return nil, nil
}

// FindOne :
func (r ScheduleRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	extras, ok := filter.Extras.(models.Scheduler)
	if !ok {
		return nil, errors.New("Cannot convert to 'Scheduler'")
	}

	switch extras.Queue {
	case "parsing":
		fallthrough
	case "analyser":
		return findOneAnalyser(entity, filter)
	}

	return nil, fmt.Errorf("Cannot find ScheduleRepository Queue :%s", extras.Queue)
}

// Create :
func (r ScheduleRepository) Create(entity interface{}, filter models.QueryFilter) (_entity interface{}, err error) {
	extras, ok := filter.Extras.(models.Scheduler)
	if !ok {
		return nil, errors.New("Cannot convert to 'Scheduler'")
	}

	switch extras.Queue {
	case "parsing":
		fallthrough
	case "analyser":
		return createAnalyser(entity, filter)
	}

	return nil, fmt.Errorf("Cannot find ScheduleRepository Queue :%s", extras.Queue)
}

// Update :
func (r ScheduleRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	// TODO

	return nil, nil
}

// Delete :
func (r ScheduleRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	// TODO

	return nil
}

func createAnalyser(entity interface{}, filter models.QueryFilter) (result models.Analyser, err error) {
	result, ok := entity.(models.Analyser)
	if !ok {
		return result, errors.New("Cannot convert Extras to 'Analyser'")
	}

	result.Company = filter.User.Company
	result.ExternalID = uuid.New().String()
	result.URLHash = helpers.Hash(string(result.URL))

	sSQL := `INSERT INTO analyser SET name=?, enabled=?, url=?, url_hash=?, external_id=?, type_of=?, company_id=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return result, fmt.Errorf("createAnalyser Prepare : %s", err.Error())
	}

	res, err := stmt.Exec(result.Name, result.Enabled, result.URL, result.URLHash, result.ExternalID, result.TypeOf, result.Company.ID)
	if err != nil {
		return result, fmt.Errorf("createAnalyser Exec : %s", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return result, fmt.Errorf("createAnalyser LastInsertId : %s", err.Error())
	}

	result.ID = id

	return result, nil
}

func findOneAnalyser(entity interface{}, filter models.QueryFilter) (result models.Analyser, err error) {
	scheduler, ok := entity.(models.Scheduler)
	if !ok {
		return result, errors.New("Cannot convert Extras to 'Scheduler'")
	}

	var (
		company models.Company
		country models.Country
	)

	sSQL := `
	SELECT
		analyser.id, analyser.name, analyser.url, analyser.url_hash, analyser.enabled, analyser.external_id, analyser.type_of, analyser.last_schedule_at,
		company.id, company.name, company.enabled, company.external_id,
		country.id, country.name, country.iso
	FROM
		analyser AS analyser
		INNER JOIN company AS company ON (company.id = analyser.company_id)
		LEFT JOIN country AS country ON (country.id = company.country_id)
	WHERE
		analyser.id = ? AND analyser.type_of = ? AND analyser.company_id = ?
	`

	err = getCon().QueryRow(sSQL, scheduler.ID, scheduler.Queue, filter.User.Company.ID).Scan(&result.ID, &result.Name, &result.URL, &result.URLHash, &result.Enabled, &result.ExternalID, &result.TypeOf, &result.LastScheduleAt,
		&company.ID, &company.Name, &company.Enabled, &company.ExternalID,
		&country.ID, &country.Name, &country.Iso)
	if nil != err {
		return result, fmt.Errorf("findOneAnalyser QueryRow : %s", err.Error())
	}

	company.Country = &country

	result.Company = &company

	return result, nil
}
