package store

import "gorm.io/gorm"

// type Service struct {
// 	gorm.Model
// 	Name        string   `gorm:"type:varchar(255);unique_index"`
// 	Image       string   `gorm:"type:varchar(255)"`
// 	Command     string   `gorm:"type:varchar(255)"`
// 	Restart     string   `gorm:"type:varchar(255)"`
// 	WorkingDir  string   `gorm:"type:varchar(255)"`
// 	Environment []string `gorm:"type:varchar(255)"`
// 	Ports       []string `gorm:"type:varchar(255)"`
// }

type Project struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(255);unique_index"`
	Repository  string   `gorm:"type:varchar(255)"`
	Branch      string   `gorm:"type:varchar(255)"`
	Image       string   `gorm:"type:varchar(255)"`
	Environment []string `gorm:"type:varchar(255)"`
	Command     string   `gorm:"type:varchar(255)"`
	Restart     string   `gorm:"type:varchar(255)"`
	WorkingDir  string   `gorm:"type:varchar(255)"`
	Ports       []string `gorm:"type:varchar(255)"`
}
