package models

type AuthInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

