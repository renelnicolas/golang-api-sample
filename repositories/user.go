package repositories

import (
	"fmt"
	"strings"

	"github.com/omt/go-rest-api/helper"
	"ohmytech.io/platform/models"
)

// UserRepository :
type UserRepository struct {
	// user.id, user.company_id, user.country_id, user.email, user.password, user.firstname, user.lastname, user.phone, user.enabled, user.roles, user.registred_company, user.registred_ip, user.registred_at,user.external_id
}

// Count :
func (r UserRepository) Count(filter models.QueryFilter) (int64, error) {
	var (
		count  int64
		nvargs []interface{}
		where  = ``
	)

	nvargs = append(nvargs, filter.User.Company.ID)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND (user.firstname REGEXP '` + searchFilter + `' OR user.lastname REGEXP '` + searchFilter + `' OR user.email REGEXP '` + searchFilter + `')` // TODO : Becarefull, must be sanitize
	}

	sSQL := `SELECT COUNT(user.id) AS counter FROM user AS user WHERE user.company_id = ?` + where

	err := getCon().QueryRow(sSQL, nvargs...).Scan(&count)
	if nil != err {
		return 0, fmt.Errorf("UserRepository Count QueryRow : %s", err.Error())
	}

	return count, nil
}

// FindAll :
func (r UserRepository) FindAll(filter models.QueryFilter) ([]interface{}, error) {
	var (
		entities []interface{}
		nvargs   []interface{}
		where    = ``
	)

	nvargs = append(nvargs, filter.User.Company.ID)
	nvargs = append(nvargs, filter.Offset)
	nvargs = append(nvargs, filter.Limit)

	searchFilter := strings.TrimFunc(filter.Search, helper.TrimWhitespaceFn)

	if "" != searchFilter {
		where += ` AND (user.firstname REGEXP '` + searchFilter + `' OR user.lastname REGEXP '` + searchFilter + `' OR user.email REGEXP '` + searchFilter + `')` // TODO : Becarefull, must be sanitize
	}

	sSQL := `
	SELECT
		user.id, user.email, user.firstname, user.lastname, user.phone, user.enabled, user.roles, user.registred_company, user.registred_ip, user.registred_at, user.external_id,
		company.id, company.name, company.enabled, company.external_id,
		country.id, country.name, country.iso
	FROM
		user AS user
		LEFT JOIN company AS company ON(company.id = user.company_id)
		LEFT JOIN country AS country ON(country.id = user.country_id)
	WHERE
		user.company_id = ?` + where + `
	ORDER BY
		user.` + filter.Order + ` ` + filter.Sort + `
	LIMIT ?,?`

	rows, err := getCon().Query(sSQL, nvargs...)
	if nil != err {
		return nil, fmt.Errorf("UserRepository FindAll Query : %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			user    models.User
			company models.Company
			country models.Country
		)

		if err := rows.Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname, &user.Phone, &user.Enabled, &user.Roles, &user.CompanyRegistered, &user.RegistredIP, &user.RegistredAt, &user.ExternalID,
			&company.ID, &company.Name, &company.Enabled, &company.ExternalID,
			&country.ID, &country.Name, &country.Iso); nil != err {
			return nil, fmt.Errorf("UserRepository FindAll Scan : %s", err.Error())
		}

		user.Company = &company
		user.Country = &country

		entities = append(entities, user)
	}

	if err = rows.Err(); nil != err {
		return nil, fmt.Errorf("UserRepository FindAll Next : %s", err.Error())
	}

	return entities, nil
}

// FindOne :
func (r UserRepository) FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	var (
		company models.Company
		country models.Country
		nvargs  []interface{}
	)

	user := entity.(models.User)

	nvargs = append(nvargs, filter.User.Company.ID)
	nvargs = append(nvargs, user.ID)

	sSQL := `
	SELECT
		user.id, user.email, user.firstname, user.lastname, user.phone, user.enabled, user.roles, user.registred_company, user.registred_ip, user.registred_at, user.external_id,
		company.id, company.name, company.enabled, company.external_id,
		country.id, country.name, country.iso
	FROM
		user AS user
		LEFT JOIN company AS company ON(company.id = user.company_id)
		LEFT JOIN country AS country ON(country.id = user.country_id)
	WHERE
		user.company_id = ? AND user.id = ?`

	// usr.salt, usr.password, usr.password_requested_at, usr.confirmation_token
	err := getCon().QueryRow(sSQL, nvargs...).Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname, &user.Phone, &user.Enabled, &user.Roles, &user.CompanyRegistered, &user.RegistredIP, &user.RegistredAt, &user.ExternalID,
		&company.ID, &company.Name, &company.Enabled, &company.ExternalID,
		&country.ID, &country.Name, &country.Iso)
	if nil != err {
		return nil, fmt.Errorf("UserRepository FindOne QueryRow : : %s", err.Error())
	}

	user.Company = &company
	user.Country = &country

	return user, nil
}

// FindOneByEmail :
func (r UserRepository) FindOneByEmail(entity *models.User) error {
	var (
		company  models.Company
		ccountry models.Country
		ucountry models.Country
	)

	sSQL := `
	SELECT
		user.id, user.email, user.password, user.firstname, user.lastname, user.phone, user.enabled, user.roles, user.registred_company, user.registred_ip, user.registred_at, user.external_id,
		company.id, company.name, company.enabled, company.external_id,
		country.id, country.name, country.iso,
		ccountry.id, ccountry.name, ccountry.iso
	FROM
		user AS user
		INNER JOIN company AS company ON (company.id = user.company_id)
		LEFT JOIN country AS country ON (country.id = user.country_id)
		LEFT JOIN country AS ccountry ON (ccountry.id = company.country_id)
	WHERE
		user.email=? AND user.enabled=1 AND company.enabled=1`

	err := getCon().QueryRow(sSQL, entity.Email).Scan(&entity.ID, &entity.Email, &entity.Password, &entity.Firstname, &entity.Lastname, &entity.Phone, &entity.Enabled, &entity.Roles, &entity.CompanyRegistered, &entity.RegistredIP, &entity.RegistredAt, &entity.ExternalID,
		&company.ID, &company.Name, &company.Enabled, &company.ExternalID,
		&ucountry.ID, &ucountry.Name, &ucountry.Iso,
		&ccountry.ID, &ccountry.Name, &ccountry.Iso)
	if nil != err {
		return fmt.Errorf("UserRepository FindOneByEmail QueryRow : %s", err.Error())
	}

	company.Country = &ccountry

	entity.Company = &company
	entity.Country = &ucountry

	return nil
}

// Create :
func (r UserRepository) Create(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	// TODO

	return nil, nil
}

// Update :
func (r UserRepository) Update(entity interface{}, filter models.QueryFilter) (interface{}, error) {
	// TODO

	return nil, nil
}

// Delete :
func (r UserRepository) Delete(entity interface{}, filter models.QueryFilter) error {
	// TODO

	return nil
}
