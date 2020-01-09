package recentsearch

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/juxemburg/truora_server/controllers/filters"
	"github.com/juxemburg/truora_server/controllers/httpresult"
	"github.com/juxemburg/truora_server/dal/entities"
)

func getRecentSearches(w http.ResponseWriter, r *http.Request) {
	httpresult.HandleRequestResponse(w, r, func() (interface{}, error) {
		return entities.GetRecentSearches()
	})
}

/*RecentSearchControllerRoutes ...*/
var RecentSearchControllerRoutes = map[string]func(chi.Router){
	"": func(r chi.Router) {
		r.Use(filters.AuthFilter)
		r.Get("/", getRecentSearches) // GET /
	},
}
