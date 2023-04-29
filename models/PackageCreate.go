package models

type PackageCreate struct {
	Name string `json:"Name"`
	Version string `json:"Version"`
	Content string `json:"Content"`
	URL string `json:"URL"`
	ID uint `json:"id" gorm:"primary_key"`
}

