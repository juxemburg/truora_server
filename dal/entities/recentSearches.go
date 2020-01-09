package entities

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/juxemburg/truora_server/apierrors"
	"github.com/juxemburg/truora_server/controllers/common"
	"github.com/juxemburg/truora_server/dal/database"
)

/*RecentSearch ...*/
type RecentSearch struct {
	HostID      string
	LastVisited time.Time
}

/*RecentSearchMetadata ...*/
type RecentSearchMetadata struct {
	HostID      string
	LastVisited time.Time
	LogoURL     string
	PageTitle   string
}

/*InsertRecentSearch inserts, or updates if exists, a recently searched host into the database*/
func InsertRecentSearch(hostID string) error {
	dbContext := database.GetDBContext()
	currentTime := time.Now()
	var statements = []string{
		fmt.Sprintf(`
			DELETE 
			FROM serverDB.recent_searches 
			WHERE  hostId = '%v'`, hostID),
		fmt.Sprintf(`INSERT INTO serverDB.recent_searches (hostId, last_visited) 
			VALUES('%v', '%v')`, hostID, currentTime.Format(common.DateDbFormat)),
	}
	return dbContext.DbExecution(statements)
}

/*GetRecentSearches returns a list with all the recent searches, along with some metadata*/
func GetRecentSearches() ([]*RecentSearchMetadata, error) {
	dbContext := database.GetDBContext()
	statement := `SELECT s.hostId, s.last_visited, d.logoUrl, d.pageTitle
	FROM serverDB.recent_searches s, serverDB.domain_info d
	WHERE s.hostId = d.host
	ORDER BY s.last_visited desc;`
	result, dberr := dbContext.DbExtraction(statement, false, func(rows *sql.Rows) (interface{}, error) {
		var searches []*RecentSearchMetadata
		for rows.Next() {
			var hostID, logoURL, pageTitle string
			var lastVisited time.Time
			if err := rows.Scan(&hostID, &lastVisited, &logoURL, &pageTitle); err != nil {
				return nil, apierrors.NewErrSQL(err.Error())
			}
			searches = append(searches,
				&RecentSearchMetadata{
					HostID:      hostID,
					LastVisited: lastVisited,
					LogoURL:     logoURL,
					PageTitle:   pageTitle,
				})
		}
		return searches, nil
	})
	if dberr != nil {
		return nil, dberr
	}
	searches, casted := result.([]*RecentSearchMetadata)
	if !casted {
		return nil, apierrors.NewErrSQL("Error while retrieving recent searches")
	}
	return searches, dberr
}
