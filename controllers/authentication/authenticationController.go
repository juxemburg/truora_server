package authentication

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/juxemburg/truora_server/controllers/filters"
	"github.com/juxemburg/truora_server/controllers/httpresult"
	"github.com/juxemburg/truora_server/services/authentication"
)

func logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout action ok"))
}
func login(w http.ResponseWriter, r *http.Request) {
	httpresult.HandleRequestResponse(w, r, func() (interface{}, error) {
		var viewModel authentication.LoginViewModel
		bodyErr := httpresult.HandleRequestBody(w, r, &viewModel)
		if bodyErr != nil {
			return nil, bodyErr
		}

		return authentication.LoginService(viewModel)
	})
}

/*AuthenticationControllerRoutes ...*/
var AuthenticationControllerRoutes = map[string]func(chi.Router){
	"login": func(r chi.Router) {
		r.Post("/", login) // POST /login
	},
	"logout": func(r chi.Router) {
		r.Use(filters.AuthFilter)
		r.Get("/", logout) // POST /logout
	},
}
