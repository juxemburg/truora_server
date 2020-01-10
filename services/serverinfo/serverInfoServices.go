package serverinfo

import (
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/juxemburg/truora_server/apierrors"
	"github.com/juxemburg/truora_server/dal/entities"
	"github.com/juxemburg/truora_server/services/htmlinfo"
	"github.com/juxemburg/truora_server/services/restclient"
)

const (
	unavailableGradeValue = "—"
)

var (
	serverGrades = [6]string{"M", "T", "A-F", " A-", "A", "A+"}
)

type domainInfo struct {
	Host            string
	Port            int
	Protocol        string
	IsPublic        bool
	Status          string
	StartTime       int64
	TestTime        int64
	EngineVersion   string
	CriteriaVersion string
	Endpoints       []struct {
		IPAddress         string
		ServerName        string
		StatusMessage     string
		Grade             string
		GradeTrustIgnored string
		HasWarnings       bool
		IsExceptional     bool
		Progress          int
		Duration          int
		Delegation        int
	}
}

func (di domainInfo) toDomainInfo(prevDomainInfo *entities.DomainInfo) *DomainInfoViewModel {
	var serverinfo []*ServerViewModel
	var currentMaxGrade = unavailableGradeValue
	var pageInfo *htmlinfo.PageInfo
	previousSslGrade := "unknown"
	serversChanged := false

	if prevDomainInfo == nil {
		pageInfo = htmlinfo.GetHTMLPageInfo(di.Host)
	} else {
		pageInfo = &htmlinfo.PageInfo{Title: prevDomainInfo.Title, IconURL: prevDomainInfo.Logo}
		serversChanged = prevDomainInfo.ServersChanged
	}

	for _, server := range di.Endpoints {
		serverinfo = append(serverinfo, &ServerViewModel{
			Address:  server.IPAddress,
			SslGrade: server.Grade,
			Country:  "",
			Owner:    "",
		})
		currentMaxGrade = maxGrade(currentMaxGrade, server.Grade)
	}
	newDomaininfo := &DomainInfoViewModel{
		Host:             di.Host,
		ServersChanged:   false,
		SslGrade:         currentMaxGrade,
		PreviousSslGrade: previousSslGrade,
		Logo:             pageInfo.IconURL,
		Title:            pageInfo.Title,
		IsDown:           di.Status == "ERROR",
		Servers:          serverinfo,
	}
	newDomainEntity := newDomaininfo.toEntity()
	if prevDomainInfo != nil {
		serversChanged = !reflect.DeepEqual(prevDomainInfo, newDomainEntity)
	}
	newDomaininfo.ServersChanged = serversChanged
	return newDomaininfo
}

func maxGrade(g1 string, g2 string) string {
	for _, grade := range serverGrades {
		if g1 == grade || g2 == grade {
			return grade
		}
	}
	return unavailableGradeValue
}

/*ServerViewModel ...*/
type ServerViewModel struct {
	Address  string `json:"address"`
	SslGrade string `json:"sslGrade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

func (vw ServerViewModel) toEntity(host string) *entities.DomainServer {
	return &entities.DomainServer{
		Address:  vw.Address,
		SslGrade: vw.SslGrade,
		Country:  vw.Country,
		Owner:    vw.Owner,
	}
}

/*DomainInfoViewModel ...*/
type DomainInfoViewModel struct {
	Host             string             `json:"host"`
	ServersChanged   bool               `json:"serversChanged"`
	SslGrade         string             `json:"sslGrade"`
	PreviousSslGrade string             `json:"previousSslGrade"`
	Logo             string             `json:"logo"`
	Title            string             `json:"title"`
	IsDown           bool               `json:"isDown"`
	Servers          []*ServerViewModel `json:"servers"`
}

func (vw DomainInfoViewModel) toEntity() *entities.DomainInfo {
	var servers []*entities.DomainServer
	for _, vwServer := range vw.Servers {
		servers = append(servers, vwServer.toEntity(vw.Host))
	}
	return &entities.DomainInfo{
		Host:           vw.Host,
		ServersChanged: vw.ServersChanged,
		SslGrade:       vw.SslGrade,
		Logo:           vw.Logo,
		Title:          vw.Title,
		IsDown:         vw.IsDown,
		Servers:        servers,
	}
}

func toViewModel(e *entities.DomainInfo) *DomainInfoViewModel {
	var servers []*ServerViewModel
	for _, server := range e.Servers {
		servers = append(servers,
			&ServerViewModel{
				Address:  server.Address,
				SslGrade: server.SslGrade,
				Country:  server.Country,
				Owner:    server.Owner,
			})
	}
	return &DomainInfoViewModel{
		Host:             e.Host,
		ServersChanged:   e.ServersChanged,
		SslGrade:         e.SslGrade,
		PreviousSslGrade: "unkown",
		Logo:             e.Logo,
		Title:            e.Title,
		IsDown:           e.IsDown,
		Servers:          servers,
	}
}

/*GetDomainInfo gets the information of a given domain*/
func GetDomainInfo(domainName string) (*DomainInfoViewModel, error) {
	domain, err := processDomainURL(domainName)
	if err != nil {
		return nil, err
	}
	//Query previous domain info first
	prevDomainInfo, err := entities.GetDomainInfo(domain)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf(`https://api.ssllabs.com/api/v3/analyze?host=%v`, domain)
	var result domainInfo
	err = restclient.GetJSON(uri, &result)
	if err != nil {
		if prevDomainInfo != nil {
			return toViewModel(prevDomainInfo), nil
		}
		return nil, err
	}

	domainInfo := result.toDomainInfo(prevDomainInfo)
	err = entities.InsertDomainInfo(domainInfo.toEntity())
	if err != nil {
		log.Println(err.Error())
		return nil, apierrors.NewApplicationError("Error while inserting into the database")
	}
	searchErr := entities.InsertRecentSearch(domain)
	if searchErr != nil {
		log.Println(err.Error())
	}
	return domainInfo, err
}

func processDomainURL(domainName string) (string, error) {
	uriRegexp := regexp.MustCompile(`^[https:\/\/]?[www\.]?[A-Za-z0-9\.]+`)
	if !uriRegexp.MatchString(domainName) {
		return "", apierrors.NewErrBadRequest("Invalid domain name 😒")
	}

	var re = regexp.MustCompile(`^(https:\/\/)|(www\.)`)
	return re.ReplaceAllString(domainName, ``), nil
}
