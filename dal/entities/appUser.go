package entities

import (
	"fmt"

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
	db, dberr := dbContext.DbConnection()
	if dberr != nil {
		return nil, dberr
	}

	statement := fmt.Sprintf(`select * from serverDB.app_users where id = %d`, userID)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var login, password string
		if err := rows.Scan(&id, &login, &password); err != nil {
			return nil, err
		}
		return &AppUser{ID: id, Login: login, password: password}, nil
	}

	return nil, nil
}
