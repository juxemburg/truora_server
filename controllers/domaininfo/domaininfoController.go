package domaininfo

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/juxemburg/truora_server/apierrors"
	"github.com/juxemburg/truora_server/controllers/filters"
	"github.com/juxemburg/truora_server/controllers/httpresult"
	"github.com/juxemburg/truora_server/services/serverinfo"
)

func getDomainInfo(w http.ResponseWriter, r *http.Request) {
	httpresult.HandleRequestResponse(w, r, func() (interface{}, error) {
		keys, ok := r.URL.Query()["domain"]
		if !ok || len(keys[0]) < 1 {
			return nil, apierrors.NewErrBadRequest("Url Param 'domain' is missing")
		}

		return serverinfo.GetDomainInfo(keys[0])
	})
}

/*DomainInfoControllerRoutes ...*/
var DomainInfoControllerRoutes = map[string]func(chi.Router){
	"host": func(r chi.Router) {
		r.Use(filters.AuthFilter)
		r.Get("/", getDomainInfo) // GET /host?domain='www.google.com'
	},
}
