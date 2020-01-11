package whois

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/juxemburg/truora_server/apierrors"
)

/*Details contains the whois command info from an ip address*/
type Details struct {
	Country string
	Owner   string
}

/*GetWhoIsCommandDetails ...*/
func GetWhoIsCommandDetails(ipAddress string) *Details {
	whois, err := runWhoisCommand("-v", ipAddress)
	if err != nil {
		return &Details{Country: "unknown", Owner: "unknown"}
	}
	cmdResult := string(whois)
	country := extractInfoFromCommandResult("Admin Country", cmdResult)
	city := extractInfoFromCommandResult("Admin City", cmdResult)
	owner := extractInfoFromCommandResult("Admin Organization", cmdResult)

	return &Details{Country: fmt.Sprintf("%v, %v",city, country), Owner: owner}
}

func extractInfoFromCommandResult(description string, cmdResult string) string {
	command := fmt.Sprintf(`(%v: )(\w+)`, description)
	regexp := regexp.MustCompile(command)
	if !regexp.MatchString(cmdResult) {
		return "Unknown"
	}
	results := regexp.FindStringSubmatch(cmdResult)

	if len(results) < 3 {
		return "Unknown"
	}

	return results[2]

}

func runWhoisCommand(args ...string) ([]byte, error) {

	var cmd *exec.Cmd
	if runtime.GOOS != "windows" {
		// this operation is not yet supported on linux/unix systems
		return nil, apierrors.NewApplicationError("unsupported operation")
	}

	cmd = exec.Command("whois", args...)
	out, cmdErr := cmd.CombinedOutput()
	if cmdErr != nil {
		return nil, apierrors.NewApplicationError(cmdErr.Error())
	}

	return out, nil
}
