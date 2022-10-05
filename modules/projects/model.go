package projects

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ID   int    `json:"id" gorm:"primaryKey" default:"uuid_generate_v4()"`
	Name string `json:"name"`

	JavaScript string `json:"javascript"`
	HTML       string `json:"html"`
	Go         string `json:"go"`
}

func NewProject(name string) *Project {
	return &Project{
		Name: name,
	}
}
