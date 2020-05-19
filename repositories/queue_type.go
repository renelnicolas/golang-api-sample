package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/omt/go-rest-api/helper"
	"ohmytech.io/platform/models"
)

// QueueTypeRepository :
type QueueTypeRepository struct {
	// id, name, enabled, config, external_id
}

// Count :
func (r QueueTypeRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count  int64
		nvargs []interface{}
		where  = `1=1`
	)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND queue_type.name REGEXP '` + searchFilter + `'` // TODO : Becarefull, must be sanitize
	}

	sSQL := `SELECT COUNT(queue_type.id) AS counter FROM queue_type AS queue_type WHERE ` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("QueueTypeRepository Count QueryRow : %s", err.Error())
	}

	return count, nil
}

// FindAll :
func (r QueueTypeRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	var (
		entities []interface{}
		nvargs   []interface{}
		where    = `1=1`
	)

	nvargs = append(nvargs, filter.Offset)
	nvargs = append(nvargs, filter.Limit)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND queue_type.name REGEXP '` + searchFilter + `'` // TODO : Becarefull, must be sanitize
	}

	sSQL := `
	SELECT
		queue_type.id, queue_type.name, queue_type.enabled, queue_type.config, queue_type.external_id
	FROM
		queue_type AS queue_type
	WHERE
		` + where + `
	ORDER BY
		queue_type.` + filter.Order + ` ` + filter.Sort + `
	LIMIT ?,?`

	rows, err := getCon().Query(sSQL, nvargs...)
	if nil != err {
		return nil, fmt.Errorf("QueueTypeRepository FindAll Query : %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			_entity models.QueueType
		)

		if err := rows.Scan(&_entity.ID, &_entity.Name, &_entity.Enabled, &_entity.Config, &_entity.ExternalID); nil != err {
			return nil, fmt.Errorf("QueueTypeRepository FindAll Scan : %s", err.Error())
		}

		entities = append(entities, _entity)
	}

	if err = rows.Err(); nil != err {
		return nil, fmt.Errorf("QueueTypeRepository FindAll Next : %s", err.Error())
	}

	return entities, nil
}

// FindOne :
func (r QueueTypeRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	_entity, ok := entity.(models.QueueType)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueType'")
	}

	sSQL := `
	SELECT
		queue_type.id, queue_type.name, queue_type.enabled, queue_type.config, queue_type.external_id
	FROM
		queue_type AS queue_type
	WHERE
		queue_type.id = ?`

	err := getCon().QueryRow(sSQL, _entity.ID).Scan(&_entity.ID, &_entity.Name, &_entity.Enabled, &_entity.Config, &_entity.ExternalID)
	if nil != err {
		return nil, fmt.Errorf("QueueTypeRepository FindOne QueryRow : %s", err.Error())
	}

	return _entity, nil
}

// Create :
func (r QueueTypeRepository) Create(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	_entity, ok := entity.(models.QueueType)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueType'")
	}

	_entity.ExternalID = uuid.New().String()

	sSQL := `INSERT INTO queue_type SET name=?, enabled=?, config=?, external_id=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return nil, fmt.Errorf("QueueTypeRepository Prepare : %s", err.Error())
	}

	res, err := stmt.Exec(_entity.Name, _entity.Enabled, _entity.Config, _entity.ExternalID)
	if err != nil {
		return nil, fmt.Errorf("QueueTypeRepository Exec : %s", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("QueueTypeRepository LastInsertId : %s", err.Error())
	}

	_entity.ID = id

	return _entity, nil
}

// Update :
func (r QueueTypeRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	_entity, ok := entity.(models.QueueType)
	if !ok {
		return nil, errors.New("Cannot convert to 'QueueType'")
	}

	var (
		nvargs []interface{}
	)

	nvargs = append(nvargs, _entity.Name)
	nvargs = append(nvargs, _entity.Enabled)
	nvargs = append(nvargs, _entity.Config)
	nvargs = append(nvargs, _entity.ID)

	sSQL := `UPDATE queue_type SET name=?, enabled=?, config=? WHERE id=?`

	stmt, err := getCon().Prepare(sSQL)
	if err != nil {
		return nil, fmt.Errorf("QueueTypeRepository Update Prepare : %s", err.Error())
	}

	_, err = stmt.Exec(nvargs...)
	if err != nil {
		return nil, fmt.Errorf("QueueTypeRepository Update Exec : %s", err.Error())
	}

	return _entity, nil
}

// Delete :
func (r QueueTypeRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	_, ok := entity.(models.QueueType)
	if !ok {
		return errors.New("Cannot convert to 'QueueType'")
	}

	return nil
}
