package restclient

import (
	"encoding/json"
	"github.com/juxemburg/truora_server/apierrors"
	"net/http"
)

/*GetJSON returns a json body from a uri endpoint*/
func GetJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return apierrors.NewErrBadGateway("Error while contacting an external server")
	}

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
