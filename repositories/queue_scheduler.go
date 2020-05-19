package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/omt/go-rest-api/helper"
	"ohmytech.io/platform/helpers"
	"ohmytech.io/platform/models"
)

// QueueSchedulerRepository :
type QueueSchedulerRepository struct {
	// id, company_id, queue_type_id, name, url, url_hash, enabled, config, external_id, last_schedule_at
}

// Count :
func (r QueueSchedulerRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count  int64
		nvargs []interface{}
		where  = `1=1`
	)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND (queue_scheduler.name REGEXP '` + searchFilter + `' OR queue_scheduler.url REGEXP '` + searchFilter + `')` // TODO : Becarefull, must be sanitize
	}

	sSQL := `SELECT COUNT(queue_scheduler.id) AS counter FROM queue_scheduler AS queue_scheduler WHERE ` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("QueueSchedulerRepository Count QueryRow : %s", err.Error())
	}

	return count, nil
}

// FindAll :
func (r QueueSchedulerRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	var (
		entities []interface{}
		nvargs   []interface{}
		where    = `company.id=?`
	)

	nvargs = append(nvargs, filter.User.Company.ID)
	nvargs = append(nvargs, filter.Offset)
	nvargs = append(nvargs, filter.Limit)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND (queue_scheduler.name REGEXP '` + searchFilter + `' OR queue_scheduler.url REGEXP '` + searchFilter + `')` // TODO : Becarefull, must be sanitize
	}

	sSQL := `
	SELECT
		queue_scheduler.id, queue_scheduler.name, queue_scheduler.url, queue_scheduler.url_hash, queue_scheduler.enabled, queue_scheduler.config, queue_scheduler.external_id, queue_scheduler.last_schedule_at,
		queue_type.id, queue_type.name, queue_type.config, queue_type.external_id,
		company.id, company.name, company.enabled, company.external_id,
		country.id, country.name, country.iso
	FROM
		queue_scheduler AS queue_scheduler
		INNER JOIN queue_type AS queue_type ON (queue_type.id = queue_scheduler.queue_type_id)
		INNER JOIN company AS company ON (company.id = queue_scheduler.company_id)
		LEFT JOIN country AS country ON (country.id = company.country_id)
	WHERE
		` + where + `
	ORDER BY
		queue_scheduler.` + filter.Order + ` ` + filter.Sort + `
	LIMIT ?,?`

	rows, err := getCon().Query(sSQL, nvargs...)
	if nil != err {
		return nil, fmt.Errorf("QueueSchedulerRepository FindAll Query : %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			queueScheduler models.QueueScheduler
			queueType      models.QueueType
			company        models.Company
			country        models.Country
		)

		if err := rows.Scan(&queueScheduler.ID, &queueScheduler.Name, &queueScheduler.URL, &queueScheduler.URLHash, &queueScheduler.Enabled, &queueScheduler.Config, &queueScheduler.ExternalID, &queueScheduler.LastScheduleAt,
			&queueType.ID, &queueType.Name, &queueType.Config, &queueType.ExternalID,
			&company.ID, &company.Name, &company.Enabled, &company.ExternalID,
			&country.ID, &country.Name, &country.Iso); nil != err {
			return nil, fmt.Errorf("QueueSchedulerRepository FindAll Scan : %s", err.Error())
		}

		company.Country = &country

		queueScheduler.QueueType = &queueType
		queueScheduler.Company = &company

		entities = append(entities, queueScheduler)
	}

	if err = rows.Err(); nil != err {
		return nil, fmt.Errorf("QueueSchedulerRepository FindAll Next : %s", err.Error())
	}

	return entities, nil
}

// FindOne :
func (r QueueSchedulerRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	queueScheduler, ok := entity.(models.QueueScheduler)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueScheduler'")
	}

	var (
		queueType models.QueueType
		company   models.Company
		country   models.Country
		nvargs    []interface{}
		where     = `queue_scheduler.id = ? AND company.id = ?`
	)

	nvargs = append(nvargs, queueScheduler.ID)
	nvargs = append(nvargs, filter.User.Company.ID)

	if nil != queueScheduler.QueueType {
		nvargs = append(nvargs, queueScheduler.QueueType.ID)
		where += ` AND queue_type.id = ?`
	}

	sSQL := `
	SELECT
		queue_scheduler.id, queue_scheduler.name, queue_scheduler.url, queue_scheduler.url_hash, queue_scheduler.enabled, queue_scheduler.config, queue_scheduler.external_id, queue_scheduler.last_schedule_at,
		queue_type.id, queue_type.name, queue_type.config, queue_type.external_id,
		company.id, company.name, company.enabled, company.external_id,
		country.id, country.name, country.iso
	FROM
		queue_scheduler AS queue_scheduler
		INNER JOIN queue_type AS queue_type ON (queue_type.id = queue_scheduler.queue_type_id)
		INNER JOIN company AS company ON (company.id = queue_scheduler.company_id)
		LEFT JOIN country AS country ON (country.id = company.country_id)
	WHERE
		` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&queueScheduler.ID, &queueScheduler.Name, &queueScheduler.URL, &queueScheduler.URLHash, &queueScheduler.Enabled, &queueScheduler.Config, &queueScheduler.ExternalID, &queueScheduler.LastScheduleAt,
		&queueType.ID, &queueType.Name, &queueType.Config, &queueType.ExternalID,
		&company.ID, &company.Name, &company.Enabled, &company.ExternalID,
		&country.ID, &country.Name, &country.Iso)
	if nil != err {
		return nil, fmt.Errorf("QueueSchedulerRepository FindOne QueryRow : %s", err.Error())
	}

	company.Country = &country

	queueScheduler.QueueType = &queueType
	queueScheduler.Company = &company

	return queueScheduler, nil
}

// Create :
func (r QueueSchedulerRepository) Create(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	queueScheduler, ok := entity.(models.QueueScheduler)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueScheduler'")
	}

	queueScheduler.ExternalID = uuid.New().String()
	queueScheduler.URLHash = helpers.Hash(string(queueScheduler.URL))

	sSQL := `INSERT INTO queue_scheduler SET name=?, company_id=?, queue_type_id=?, url=?, url_hash=?, enabled=?, config=?, external_id=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerRepository Prepare : %s", err.Error())
	}

	res, err := stmt.Exec(queueScheduler.Name, queueScheduler.Company.ID, queueScheduler.QueueType.ID, queueScheduler.URL, queueScheduler.URLHash, queueScheduler.Enabled, queueScheduler.Config, queueScheduler.ExternalID)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerRepository Exec : %s", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerRepository LastInsertId : %s", err.Error())
	}

	queueScheduler.ID = id

	return queueScheduler, nil
}

// Update :
func (r QueueSchedulerRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	_entity, ok := entity.(models.QueueScheduler)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueScheduler'")
	}

	var (
		nvargs []interface{}
	)

	_entity.URLHash = helpers.Hash(string(_entity.URL))

	nvargs = append(nvargs, _entity.Name)
	nvargs = append(nvargs, _entity.QueueType.ID)
	nvargs = append(nvargs, _entity.URL)
	nvargs = append(nvargs, _entity.URLHash)
	nvargs = append(nvargs, _entity.Enabled)
	nvargs = append(nvargs, _entity.Config)
	nvargs = append(nvargs, _entity.ID)

	sSQL := `UPDATE queue_scheduler SET name=?, queue_type_id=?, url=?, url_hash=?, enabled=?, config=? WHERE id=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerRepository Update Prepare : %s", err.Error())
	}

	_, err = stmt.Exec(nvargs...)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerRepository Update Exec : %s", err.Error())
	}

	return _entity, nil
}

// Delete :
func (r QueueSchedulerRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	_, ok := entity.(models.QueueScheduler)
	if !ok {
		return errors.New("Cannot convert to 'QueueScheduler'")
	}

	return nil
}
