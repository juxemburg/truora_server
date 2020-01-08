package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/juxemburg/truora_server/apierrors"
	/*Database dirver import*/
	_ "github.com/lib/pq"
)

/*DbContext ...*/
type DbContext struct {
	driverName     string
	datasourceName string
}

func (context DbContext) dbConnection() (*sql.DB, error) {
	db, err := sql.Open(context.driverName, context.datasourceName)
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return nil, apierrors.NewErrSQL(err.Error())
	}
	return db, nil
}

/*DbExtraction retrieves an interface based on a SQL statement, alongside an extraction function */
func (context DbContext) DbExtraction(statement string, allowsNull bool, extractionFn func(rows *sql.Rows) (interface{}, error)) (interface{}, error) {
	db, dberr := context.dbConnection()
	defer db.Close()
	if dberr != nil {
		return nil, dberr
	}

	rows, rowErr := db.Query(statement)

	if rowErr != nil {
		fmt.Println(rowErr.Error())
		return nil, apierrors.NewErrSQL(rowErr.Error())
	}
	defer rows.Close()
	result, err := extractionFn(rows)
	if err != nil {
		return nil, apierrors.NewErrSQL(rowErr.Error())
	}
	if result == nil && !allowsNull {
		return nil, apierrors.NewErrNotFound("The requested resource was not found")
	}

	return result, nil
}


/*DbExecution executes multiple SQL statements in the database */
func (context DbContext) DbExecution(statements []string) error {
	db, dberr := context.dbConnection()
	defer db.Close()
	if dberr != nil {
		return dberr
	}
	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
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
