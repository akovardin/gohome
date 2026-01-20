package models

import "gorm.io/gorm"

type Package struct {
	gorm.Model

	Title   string
	Info    string
	Repo    string
	Package string
	Active  bool
}
