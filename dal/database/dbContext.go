package database

import (
	"database/sql"
	"log"

	/*Database dirver import*/
	_ "github.com/lib/pq"
)

/*DbContext ...*/
type DbContext struct {
	driverName     string
	datasourceName string
}

/*DbConnection get the connection to the database*/
func (context DbContext) DbConnection() (*sql.DB, error) {
	db, err := sql.Open(context.driverName, context.datasourceName)
	if err != nil {
		log.Println("error connecting to the database: ", err)
	}
	return db, err
}

func newDbContext(datasourceName string) *DbContext {
	return &DbContext{
		driverName:     "postgres",
		datasourceName: datasourceName,
	}
}

var context *DbContext = nil

/*GetDBContext Gets the current dbContext for the configured environmet*/
func GetDBContext() *DbContext {
	if context == nil {
		context = newDbContext("postgresql://maxroach@localhost:26257/serverDB?sslmode=disable")
	}
	return context
}
