package main

type AuthData struct {
	Username     string
	PasswordHash string //bcrypt hash of password
}

