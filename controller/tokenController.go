package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenController interface {
	Login(w http.ResponseWriter, r *http.Request) error
	Refresh(w http.ResponseWriter, r *http.Request) error
}

type tokenController struct {
	infoLog, errorLog *log.Logger
}

func NewTokenController(infoLog, errorLog *log.Logger) TokenController {
	return &tokenController{
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

type tokensPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type loginResponse struct {
	TokensPair tokensPair
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (c *tokenController) Login(w http.ResponseWriter, r *http.Request) error {
	tokenPairs, err := newTokenPairs()
	if err != nil {
		return err
	}

	login_response := loginResponse{
		TokensPair: tokenPairs,
	}

	fmt.Printf("HERE: %+v", login_response)

	response, err := json.Marshal(login_response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(response)

	return nil
}

func (c *tokenController) Refresh(w http.ResponseWriter, r *http.Request) error {
	var refreshRequest refreshRequest
	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		// Handle JSON decoding errors
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return err
	}
	defer r.Body.Close()

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	_, err = jwt.Parse(refreshRequest.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("refresh_token_key"), nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	tokenPairs, err := newTokenPairs()
	if err != nil {
		return err
	}

	login_response := loginResponse{
		TokensPair: tokenPairs,
	}

	response, err := json.Marshal(login_response)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(response)

	return nil
}

func newTokenPairs() (tokensPair, error) {
	access_token_claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, access_token_claims)

	signed_access_token, err := access_token.SignedString([]byte("access_token_key"))
	if err != nil {
		return tokensPair{}, err
	}

	refresh_token_claims := MyCustomClaims{
		"bar",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_token_claims)

	signed_refresh_token, err := refresh_token.SignedString([]byte("refresh_token_key"))
	if err != nil {
		return tokensPair{}, err
	}

	return tokensPair{
		AccessToken:  signed_access_token,
		RefreshToken: signed_refresh_token,
	}, nil
}
