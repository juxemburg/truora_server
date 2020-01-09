package entities

import (
	"database/sql"
	"fmt"

	"github.com/juxemburg/truora_server/apierrors"
	"github.com/juxemburg/truora_server/dal/database"
)

/*DomainServer ...*/
type DomainServer struct {
	Address  string //Primary key
	SslGrade string
	Country  string
	Owner    string
}

/*DomainInfo ...*/
type DomainInfo struct {
	Host           string //Primary key
	ServersChanged bool
	SslGrade       string
	Logo           string
	Title          string
	IsDown         bool
	Servers        []*DomainServer
}

/*GetDomainInfo get a DomainInfo item with its servers from the database*/
func GetDomainInfo(hostName string) (*DomainInfo, error) {
	dbContext := database.GetDBContext()
	statement := fmt.Sprintf(`select * from serverDB.domain_info where host = '%v'`, hostName)
	result, dberr := dbContext.DbExtraction(statement, true, func(rows *sql.Rows) (r interface{}, err error) {
		for rows.Next() {
			var host, sslGrade, logo, title string
			var serversChanged, isDown bool
			if err := rows.Scan(&host, &serversChanged, &sslGrade, &logo, &title, &isDown); err != nil {
				return nil, apierrors.NewErrSQL(err.Error())
			}
			return &DomainInfo{
				Host:           host,
				ServersChanged: serversChanged,
				SslGrade:       sslGrade,
				Logo:           logo,
				Title:          title,
				IsDown:         isDown,
				Servers:        []*DomainServer{},
			}, nil
		}
		return nil, nil
	})

	if result == nil {
		return nil, dberr
	}

	domainInfo, casted := result.(*DomainInfo)
	if !casted {
		return nil, apierrors.NewErrSQL("Error while retrieving the requested user")
	}
	domainServers, serverErr := getDomainServers(hostName)
	if serverErr != nil {
		return nil, serverErr
	}

	domainInfo.Servers = domainServers
	return domainInfo, dberr
}

func getDomainServers(hostName string) ([]*DomainServer, error) {
	dbContext := database.GetDBContext()
	statement := fmt.Sprintf(`select * from serverDB.domain_server where hostId = '%v'`, hostName)

	result, dberr := dbContext.DbExtraction(statement, false, func(rows *sql.Rows) (interface{}, error) {
		var servers []*DomainServer
		for rows.Next() {
			var address, hostID, sslGrade, country, owner string
			if err := rows.Scan(&address, &hostID, &sslGrade, &country, &owner); err != nil {
				return nil, apierrors.NewErrSQL(err.Error())
			}
			servers = append(servers, &DomainServer{Address: address, SslGrade: sslGrade, Country: country, Owner: owner})
		}
		return servers, nil
	})
	if dberr != nil {
		return nil, dberr
	}
	servers, casted := result.([]*DomainServer)
	if !casted {
		return nil, apierrors.NewErrSQL("Error while retrieving the domain servers")
	}
	return servers, dberr
}

/*InsertDomainInfo inserts a domain info, alongside its servers, into the database*/
func InsertDomainInfo(domainInfo *DomainInfo) error {
	dbContext := database.GetDBContext()
	var statements = []string{
		fmt.Sprintf(`
			DELETE 
			FROM serverDB.domain_server 
			WHERE  hostId = '%v'`, domainInfo.Host),
		fmt.Sprintf(`
			DELETE 
			FROM serverDB.domain_info 
			WHERE  host = '%v'`, domainInfo.Host),
		fmt.Sprintf(
			`INSERT INTO serverDB.domain_info (
			host,
			serversChanged,
			sslGrade,
			logoUrl,
			pageTitle,
			isDown
		) 
		VALUES
		(
			'%v',
			%v,
			'%v',
			'%v',
			'%v',
			%v
		)`,
			domainInfo.Host,
			domainInfo.ServersChanged,
			domainInfo.SslGrade,
			domainInfo.Logo,
			domainInfo.Title,
			domainInfo.IsDown),
	}

	for _, server := range domainInfo.Servers {
		statements = append(statements,
			fmt.Sprintf(
				`INSERT INTO serverDB.domain_server (
				ipAddress,
				hostId,
				sslGrade,
				country,
				owner
			) 
			VALUES
			(
				'%v',
				'%v',
				'%v',
				'%v',
				'%v'
			)`,
				server.Address,
				domainInfo.Host,
				server.SslGrade,
				server.Country,
				server.Owner),
		)
	}

	return dbContext.DbExecution(statements)
}
