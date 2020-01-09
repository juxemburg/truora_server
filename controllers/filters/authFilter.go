package filters

import (
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/juxemburg/truora_server/controllers/common"
)

/*AuthFilter filters HTTP request for authentication*/
func AuthFilter(next http.Handler) http.Handler {
	if len(common.AppKey) == 0 {
		log.Fatal("the APP_KEY const have not been setted")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(common.AppKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}
