package sql

import (
	"errors"
	"gorm.io/gorm"
)

var createProjectError = errors.New("create project failed")
var updateProjectNameError = errors.New("update project name failed")
var projectNotFoundError = errors.New("project not found")
var gettingAllProjectsError = errors.New("getting all projects failed")

type Project struct {
	gorm.Model
	Name   string `gorm:"type:varchar(100);not null" json:"name"`
	User   User   `gorm:"foreignkey:user_id;association_foreignkey:id"` // use UserRefer as foreign key
	UserID uint
}

func (p *Project) Save(db *gorm.DB) error {
	if result := db.Create(&p); result.Error != nil {
		return createProjectError
	}
	return nil
}

func (p *Project) UpdateName(db *gorm.DB) error {
	if result := db.Model(&p).Update("name", p.Name); result.Error != nil {
		return updateProjectNameError
	}
	return nil
}

func (p *Project) Find(db *gorm.DB) error {
	if result := db.First(&p, "id = ?", p.ID); result.Error != nil {
		return projectNotFoundError
	}
	return nil
}

func (p Project) GetAll(db *gorm.DB) ([]Project, error) {
	var projects []Project
	if result := db.Where("user_id = ?", p.UserID).Find(&projects); result.Error != nil {
		return nil, gettingAllProjectsError
	}
	return projects, nil
}
