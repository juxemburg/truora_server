package entities

import (
	"database/sql"
	"fmt"

	"github.com/juxemburg/truora_server/apierrors"
	"github.com/juxemburg/truora_server/dal/database"
)

/*AppUser represents an app user*/
type AppUser struct {
	ID       int
	Login    string
	password string
}

/*FindUser finds a user in the database with the provided id,
returns null if there is no such user*/
func FindUser(userID int) (*AppUser, error) {
	dbContext := database.GetDBContext()
	statement := fmt.Sprintf(`select * from serverDB.app_users where id = %d`, userID)
	result, dberr := dbContext.DbExtraction(statement, func(rows *sql.Rows) (r interface{}, err error) {
		for rows.Next() {
			var id int
			var login, password string
			if err := rows.Scan(&id, &login, &password); err != nil {
				return nil, apierrors.NewErrSQL(err.Error())
			}
			return &AppUser{ID: id, Login: login, password: password}, nil
		}
		return nil, nil
	})
	user, casted := result.(AppUser)
	if !casted {
		return nil, apierrors.NewErrSQL("Error while retrieving the requested user")
	}
	return &user, dberr
}

	return nil, nil
}
