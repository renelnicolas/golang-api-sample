package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/omt/go-rest-api/helper"
	"ohmytech.io/platform/models"
)

// QueueSchedulerHistoryRepository :
type QueueSchedulerHistoryRepository struct {
	// id, queue_scheduler_id, uuid, status, scheduled_start_date, scheduled_end_date
}

// Count :
func (r QueueSchedulerHistoryRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count int64
	)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	sSQL := `SELECT COUNT(queue_scheduler_history.id) AS counter FROM queue_scheduler_history AS queue_scheduler_history WHERE queue_scheduler_id=?`

	err := getCon().QueryRow(sSQL, searchFilter).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("QueueSchedulerHistoryRepository Count QueryRow : %s", err.Error())
	}

	return count, nil
}

// FindAll :
func (r QueueSchedulerHistoryRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	var (
		entities []interface{}
		nvargs   []interface{}
		where    = `queue_scheduler.company_id=? AND queue_scheduler.id=?`
	)

	nvargs = append(nvargs, filter.User.Company.ID)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	nvargs = append(nvargs, searchFilter)
	nvargs = append(nvargs, filter.Offset)
	nvargs = append(nvargs, filter.Limit)

	sSQL := `
	SELECT
		queue_scheduler_history.id, queue_scheduler_history.uuid, queue_scheduler_history.status, queue_scheduler_history.scheduled_start_date, queue_scheduler_history.scheduled_end_date,
		queue_scheduler.id, queue_scheduler.name, queue_scheduler.url, queue_scheduler.url_hash, queue_scheduler.enabled, queue_scheduler.config, queue_scheduler.external_id
	FROM
		queue_scheduler_history AS queue_scheduler_history
		INNER JOIN queue_scheduler AS queue_scheduler ON (queue_scheduler.id = queue_scheduler_history.queue_scheduler_id)
	WHERE
		` + where + `
	ORDER BY
		queue_scheduler.` + filter.Order + ` ` + filter.Sort + `
	LIMIT ?,?`

	rows, err := getCon().Query(sSQL, nvargs...)
	if nil != err {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository FindAll Query : %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			queueSchedulerHistory models.QueueSchedulerHistory
			queueScheduler        models.QueueScheduler
		)

		if err := rows.Scan(&queueSchedulerHistory.ID, &queueSchedulerHistory.UUID, &queueSchedulerHistory.Status, &queueSchedulerHistory.ScheduledStartDate, &queueSchedulerHistory.ScheduledEndDate,
			&queueScheduler.ID, &queueScheduler.Name, &queueScheduler.URL, &queueScheduler.URLHash, &queueScheduler.Enabled, &queueScheduler.Config, &queueScheduler.ExternalID); nil != err {
			return nil, fmt.Errorf("QueueSchedulerHistoryRepository FindAll Scan : %s", err.Error())
		}

		// queueSchedulerHistory.QueueScheduler = &queueScheduler
		queueSchedulerHistory.QueueSchedulerID = queueScheduler.ID

		entities = append(entities, queueSchedulerHistory)
	}

	if err = rows.Err(); nil != err {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository FindAll Next : %s", err.Error())
	}

	return entities, nil
}

// FindOne :
func (r QueueSchedulerHistoryRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	queueSchedulerHistory, ok := entity.(models.QueueSchedulerHistory)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueSchedulerHistory'")
	}

	var (
		queueScheduler models.QueueScheduler
		nvargs         []interface{}
		where          = `queue_scheduler.company_id=? AND queue_scheduler_history.id=?`
	)

	nvargs = append(nvargs, filter.User.Company.ID)
	nvargs = append(nvargs, queueSchedulerHistory.ID)

	sSQL := `
	SELECT
		queue_scheduler_history.id, queue_scheduler_history.uuid, queue_scheduler_history.status, queue_scheduler_history.scheduled_start_date, queue_scheduler_history.scheduled_end_date,
		queue_scheduler.id, queue_scheduler.name, queue_scheduler.url, queue_scheduler.url_hash, queue_scheduler.enabled, queue_scheduler.config, queue_scheduler.external_id
	FROM
		queue_scheduler_history AS queue_scheduler_history
		INNER JOIN queue_scheduler AS queue_scheduler ON (queue_scheduler.id = queue_scheduler_history.queue_scheduler_id)
	WHERE
		` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&queueSchedulerHistory.ID, &queueSchedulerHistory.UUID, &queueSchedulerHistory.Status, &queueSchedulerHistory.ScheduledStartDate, &queueSchedulerHistory.ScheduledEndDate,
		&queueScheduler.ID, &queueScheduler.Name, &queueScheduler.URL, &queueScheduler.URLHash, &queueScheduler.Enabled, &queueScheduler.Config, &queueScheduler.ExternalID)
	if nil != err {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository FindOne QueryRow : %s", err.Error())
	}

	// queueSchedulerHistory.QueueScheduler = &queueScheduler
	queueSchedulerHistory.QueueSchedulerID = queueScheduler.ID

	return queueScheduler, nil
}

// Create :
func (r QueueSchedulerHistoryRepository) Create(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	queueSchedulerHistory, ok := entity.(models.QueueSchedulerHistory)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueSchedulerHistory'")
	}

	queueSchedulerHistory.UUID = uuid.New().String()

	sSQL := `INSERT INTO queue_scheduler_history SET queue_scheduler_id=?, uuid=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository Prepare : %s", err.Error())
	}

	res, err := stmt.Exec(queueSchedulerHistory.QueueSchedulerID, queueSchedulerHistory.UUID)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository Exec : %s", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository LastInsertId : %s", err.Error())
	}

	queueSchedulerHistory.ID = id

	return queueSchedulerHistory, nil
}

// Update :
func (r QueueSchedulerHistoryRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	queueSchedulerHistory, ok := entity.(models.QueueSchedulerHistory)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueSchedulerHistory'")
	}

	var (
		nvargs []interface{}
	)

	nvargs = append(nvargs, queueSchedulerHistory.Status)
	nvargs = append(nvargs, queueSchedulerHistory.ID)

	sSQL := `UPDATE queue_scheduler SET status=?, scheduled_end_date=NOW() WHERE id=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository Update Prepare : %s", err.Error())
	}

	_, err = stmt.Exec(nvargs...)
	if err != nil {
		return nil, fmt.Errorf("QueueSchedulerHistoryRepository Update Exec : %s", err.Error())
	}

	return queueSchedulerHistory, nil
}

// Delete :
func (r QueueSchedulerHistoryRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	_, ok := entity.(models.QueueSchedulerHistory)
	if !ok {
		return errors.New("Cannot convert to 'QueueSchedulerHistory'")
	}

	return nil
}
