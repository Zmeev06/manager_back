package models

type User struct {
	Login    string   `json:"login"`
	Password []byte   `json:"password"`
	Admin    bool     `json:"admin"`
	Servers  []string `json:"servers"`
}
