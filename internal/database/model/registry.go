package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Observation struct {
	gorm.Model

	Registry string `gorm:"index:idx_docker_image"`
	Image    string `gorm:"index:idx_docker_image"`
	Tag      string `gorm:"index:idx_docker_image"`
	LastHash string

	Listener []Listener `gorm:"foreignkey:ObservationID"`
}

type Listener struct {
	gorm.Model

	ObservationID    uint   `gorm:"index"`
	Name             string `gorm:"index"`
	Method           string
	Url              string
	Header           []byte
	Body             []byte
	LastNotification *time.Time
	LastHash         string
}
