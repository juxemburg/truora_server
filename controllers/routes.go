package controllers

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/juxemburg/truora_server/controllers/authentication"
)

/*GetRouteConfig ...*/
func GetRouteConfig() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api", func(subr chi.Router) {
		routeRegistration(subr, "authentication", authentication.AuthenticationControllerRoutes)
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
