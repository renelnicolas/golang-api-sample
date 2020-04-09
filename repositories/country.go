package repositories

import (
	"fmt"

	"ohmytech.io/platform/models"
)

// CountryRepository :
type CountryRepository struct {
	// country.id, country.name, country.iso
}

// Count :
func (r CountryRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count  int64
		nvargs []interface{}
		where  = `1=1`
	)

	sSQL := `SELECT COUNT(country.id) AS counter FROM country AS country WHERE ` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("CountryRepository Count QueryRow : %s", err.Error())
	}

	return count, nil
}

// FindAll :
func (r CountryRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	var (
		entities []interface{}
		nvargs   []interface{}
		where    = `1=1`
	)

	sSQL := `
	SELECT
		country.id, country.name, country.iso
	FROM
		country AS country
	WHERE
		` + where + `
	ORDER BY
		country.` + filter.Order + ` ` + filter.Sort

	rows, err := getCon().Query(sSQL, nvargs...)
	if nil != err {
		return nil, fmt.Errorf("CountryRepository FindAll Query : %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			country models.Country
		)

		if err := rows.Scan(&country.ID, &country.Name, &country.Iso); nil != err {
			return nil, fmt.Errorf("CountryRepository FindAll Scan : %s", err.Error())
		}

		entities = append(entities, country)
	}

	if err = rows.Err(); nil != err {
		return nil, fmt.Errorf("CountryRepository FindAll Next : %s", err.Error())
	}

	return entities, nil
}

// FindOne :
func (r CountryRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	return nil, nil
}

// Create :
func (r CountryRepository) Create(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	return nil, nil
}

// Update :
func (r CountryRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	return nil, nil
}

// Delete :
func (r CountryRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	return nil
}
