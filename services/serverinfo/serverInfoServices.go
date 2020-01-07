package serverinfo

import (
	"fmt"

	"github.com/juxemburg/truora_server/services/restclient"
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

func (di domainInfo) toDomainInfo() *DomainInfo {
	var serverinfo []ServerInfo

	for _, server := range di.Endpoints {
		serverinfo = append(serverinfo, ServerInfo{
			Address:  server.IPAddress,
			SslGrade: server.Grade,
			Country:  "",
			Owner:    "",
		})
	}
	return &DomainInfo{
		ServersChanged:   false,
		SslGrade:         "",
		PreviousSslGrade: "",
		Logo:             "",
		Title:            "",
		IsDown:           di.Status == "ERROR",
		Servers:          serverinfo,
	}
}

/*ServerInfo ...*/
type ServerInfo struct {
	Address  string
	SslGrade string
	Country  string
	Owner    string
}

/*DomainInfo ...*/
type DomainInfo struct {
	ServersChanged   bool
	SslGrade         string
	PreviousSslGrade string
	Logo             string
	Title            string
	IsDown           bool
	Servers          []ServerInfo
}

/*GetDomainInfo gets the information of a given domain*/
func GetDomainInfo(domainName string) (*DomainInfo, error) {
	uri := fmt.Sprintf(`https://api.ssllabs.com/api/v3/analyze?host=%v`, domainName)
	var result domainInfo
	err := restclient.GetJSON(uri, &result)
	if err != nil {
		return nil, err
	}
	return result.toDomainInfo(), err
}
