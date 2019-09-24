package main

import (
	"golang.org/x/crypto/bcrypt"
)

type AuthData struct {
	Username     string
	IsTrial      bool   //true : no password, only username. trial users like discord
	PasswordHash []byte //bcrypt hash of password
}

func (a *AuthData) HashPassword(plaintextPassword string) {
	a.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.MinCost)
}

func (a AuthData) IsMatch(plai)
