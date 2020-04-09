package repositories

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // comment to jsutify go-lint

	"ohmytech.io/platform/config"
	"ohmytech.io/platform/models"
)

// Repositoryer :
type Repositoryer interface {
	Count(filter models.QueryFilter) (int64, error)
	FindAll(filter models.QueryFilter) ([]interface{}, error)
	FindOne(entity interface{}, filter models.QueryFilter) (interface{}, error)
	Create(entity interface{}, filter models.QueryFilter) (interface{}, error)
	Update(entity interface{}, filter models.QueryFilter) (interface{}, error)
	Delete(entity interface{}, filter models.QueryFilter) error
}

var (
	con *sql.DB
	err error
)

// CreateCon :
func createCon() {
	if nil == con {
		dbConfig := config.GetConfig().Database

		con, err = sql.Open("mysql", dbConfig.Connector)

		// if there is an error opening the connection, handle it
		if nil != err {
			panic(err.Error())
		}
	}
}

func getCon() *sql.DB {
	createCon()
	return con
}
