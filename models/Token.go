package models

type Token struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	AuthToken string `json:"AuthToken"`
}