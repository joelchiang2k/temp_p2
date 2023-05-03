package models

type Token struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Auth_token string `json:"auth_token"`
}
