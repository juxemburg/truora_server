package authentication

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/juxemburg/truora_server/controllers/httpresult"
	"github.com/juxemburg/truora_server/dal/entities"
)

type loginViewModel struct {
	Login    string
	Password string
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout action ok"))
}
func login(w http.ResponseWriter, r *http.Request) {
	httpresult.HandleRequestResponse(w, r, func() (v interface{}, err error) {
		var viewModel loginViewModel
		bodyErr := httpresult.HandleRequestBody(w, r, &viewModel)
		if bodyErr != nil {
			return nil, bodyErr
		}
		exists, err := entities.ExistUser(viewModel.Login, viewModel.Password)
		return exists, nil
	})
}

/*AuthenticationControllerRoutes ...*/
var AuthenticationControllerRoutes = map[string]func(chi.Router){
	"login": func(r chi.Router) {
		r.Post("/", login) // POST /login
	},
	"logout": func(r chi.Router) {
		r.Get("/", logout) // POST /logout
	},
}
