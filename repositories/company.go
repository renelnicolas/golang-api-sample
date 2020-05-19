package repositories

import (
	"fmt"
	"strings"

	"github.com/omt/go-rest-api/helper"
	"ohmytech.io/platform/models"
)

// CompanyRepository :
type CompanyRepository struct {
	// company.id, company.country_id, company.name, company.enabled, company.website, company.email, company.address, company.zip_code, company.city, company.phone, company.vat, company.rcs, company.external_id, company.created_at, company.updated_at
}

// Count :
func (r CompanyRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count  int64
		nvargs []interface{}
		where  = `1=1`
	)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND company.name REGEXP '` + searchFilter + `'` // TODO : Becarefull, must be sanitize
	}

	sSQL := `SELECT COUNT(company.id) AS counter FROM company AS company WHERE ` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("CompanyRepository Count QueryRow : %s", err.Error())
	}

	return count, nil
}

// FindAll :
func (r CompanyRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	var (
		entities []interface{}
		nvargs   []interface{}
		where    = `1=1`
	)

	nvargs = append(nvargs, filter.Offset)
	nvargs = append(nvargs, filter.Limit)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND company.name REGEXP '` + searchFilter + `'` // TODO : Becarefull, must be sanitize
	}

	sSQL := `
	SELECT
		company.id, company.name, company.enabled, company.website, company.contact_email, company.address, company.zip_code, company.city, company.phone, company.vat, company.rcs, company.external_id,
		country.id, country.name, country.iso
	FROM
		company AS company
		LEFT JOIN country AS country ON(country.id = company.country_id)
	WHERE
		` + where + `
	ORDER BY
		company.` + filter.Order + ` ` + filter.Sort + `
	LIMIT ?,?`

	rows, err := getCon().Query(sSQL, nvargs...)
	if nil != err {
		return nil, fmt.Errorf("CompanyRepository FindAll Query : %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			company models.Company
			country models.Country
		)

		if err := rows.Scan(&company.ID, &company.Name, &company.Enabled, &company.Website, &company.ContactEmail, &company.Address, &company.ZipCode, &company.City, &company.Phone, &company.VAT, &company.RCS, &company.ExternalID,
			&country.ID, &country.Name, &country.Iso); nil != err {
			return nil, fmt.Errorf("CompanyRepository FindAll Scan : %s", err.Error())
		}

		company.Country = &country

		entities = append(entities, company)
	}

	if err = rows.Err(); nil != err {
		return nil, fmt.Errorf("CompanyRepository FindAll Next : %s", err.Error())
	}

	return entities, nil
}

// FindOne :
func (r CompanyRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	var (
		country models.Country
		nvargs  []interface{}
	)

	company := entity.(models.Company)

	nvargs = append(nvargs, company.ID)

	sSQL := `
	SELECT
		company.id, company.name, company.enabled, company.website, company.contact_email, company.address, company.zip_code, company.city, company.phone, company.vat, company.rcs, company.external_id,
		country.id, country.name, country.iso
	FROM
		company AS company
		LEFT JOIN country AS country ON(country.id = company.country_id)
	WHERE
		company.id = ?`

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&company.ID, &company.Name, &company.Enabled, &company.Website, &company.ContactEmail, &company.Address, &company.ZipCode, &company.City, &company.Phone, &company.VAT, &company.RCS, &company.ExternalID,
		&country.ID, &country.Name, &country.Iso)
	if nil != err {
		return nil, fmt.Errorf("CompanyRepository FindOne QueryRow : %s", err.Error())
	}

	company.Country = &country

	return company, nil
}

// Create :
func (r CompanyRepository) Create(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	// TODO

	return nil, nil
}

// Update :
func (r CompanyRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	// TODO

	return nil, nil
}

// Delete :
func (r CompanyRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	// TODO

	return nil
}
