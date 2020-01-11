package controllers

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/juxemburg/truora_server/controllers/authentication"
	"github.com/juxemburg/truora_server/controllers/domaininfo"
	"github.com/juxemburg/truora_server/controllers/recentsearch"
)

/*GetRouteConfig ...*/
func GetRouteConfig() *chi.Mux {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api", func(subr chi.Router) {
		routeRegistration(subr, "authentication", authentication.AuthenticationControllerRoutes)
		routeRegistration(subr, "domainInfo", domaininfo.DomainInfoControllerRoutes)
		routeRegistration(subr, "recentsearch", recentsearch.RecentSearchControllerRoutes)
	})

	return r
}

func routeRegistration(r chi.Router, controllerName string, routes map[string]func(chi.Router)) {
	r.Route(fmt.Sprintf(`/%s`, controllerName), func(subr chi.Router) {
		for routeName, routeFn := range routes {
			subr.Route(fmt.Sprintf(`/%s`, routeName), routeFn)
		}
	})
}
