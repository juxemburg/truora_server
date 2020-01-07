package authentication

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juxemburg/truora_server/apierrors"
	"github.com/juxemburg/truora_server/controllers/common"
	"github.com/juxemburg/truora_server/dal/entities"
)

/*LoginViewModel ...*/
type LoginViewModel struct {
	Login    string
	Password string
}

/*LoginResponse ...*/
type LoginResponse struct {
	AuthToken string
}

/*LoginService ...*/
func LoginService(viewModel LoginViewModel) (*LoginResponse, error) {
	exists, err := entities.ExistUser(viewModel.Login, viewModel.Password)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, apierrors.NewErrBadRequest("Invalid login credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": viewModel.Login,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(common.APP_KEY))
	if err != nil {
		return nil, apierrors.NewApplicationError("Token generation failed")
	}
	return &LoginResponse{AuthToken: tokenString}, nil
}
