package main

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

type AuthData struct {
	Username     string
	IsTrial      bool   //true : no password, only username. trial users like discord
	PasswordHash []byte //bcrypt hash of password
}

//NewAuthData returns an authdata object. if plaintextpassword is "" create a trial account like discord
func NewAuthData(uname string, plaintextPassword string) AuthData {
	a := AuthData{
		Username: uname,
	}
	if plaintextPassword == "" {
		a.IsTrial = true
	} else {
		a.IsTrial = false
		a.HashPassword(plaintextPassword)
	}
	return a
}

func (a *AuthData) HashPassword(plaintextPassword string) {
	a.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.MinCost)
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// IsMatch verifies is passowrds match, returns http cookie with JWT(30 min) if match = true
func (a AuthData) IsMatch(plaintextPassword string) (match bool, cookie http.Cookie) {
	//TODO trial should generate JWT always
	if a.IsTrial {
		match = false
		return
	}
	if err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(plaintextPassword)); err != nil {
		match = false
		return
	}
	match = true

	//Get a JWT token for the user for 30 minutes
	cookie.Name = "jwt-token"
	cookie.Expires = time.Now().Add(30 * time.Minute)
	claims := &claims{
		Username: a.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: cookie.Expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	if cookie.Value, err = token.SignedString(jwtKey); err != nil {
		log.Fatal("JWT signing error")
	}
	return
}
