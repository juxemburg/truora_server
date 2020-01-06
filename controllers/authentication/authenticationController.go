package authentication

import (
	"github.com/go-chi/chi"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Loggin action ok"))
}
func logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout action ok"))
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
