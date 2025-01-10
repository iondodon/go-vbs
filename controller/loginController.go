package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginController interface {
	Login(w http.ResponseWriter, r *http.Request) error
}

type loginController struct {
	infoLog, errorLog *log.Logger
}

func NewLoginController(infoLog, errorLog *log.Logger) LoginController {
	return &loginController{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func (c *loginController) Login(w http.ResponseWriter, r *http.Request) error {
	claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte("signingkeysigningkeysigni"))
	if err != nil {
		return err
	}

	w.Write([]byte(ss))

	return nil
}
