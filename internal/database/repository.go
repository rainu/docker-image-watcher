package database

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/github"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/rainu/docker-image-watcher/internal/database/model"
	"io"
	"time"
)

type OverdueListener struct {
	model.Listener

	ObservationHash string
}

type Repository interface {
	io.Closer

	AddObservation(observation model.Observation) error
	GetObservations(since time.Duration) (*sql.Rows, error)
	NextObservation(rows *sql.Rows) (model.Observation, error)

	UpdateImageHash(registry, image, tag, hash string) error
	TouchObservation(id uint) error

	GetOverdueNotifications() (*sql.Rows, error)
	NextNotification(rows *sql.Rows) (OverdueListener, error)
	UpdateListener(id uint, hash string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(dbConfig string) Repository {
	db, err := gorm.Open("postgres", dbConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to connect database"))
	}

	// Migrate the schema
	db.AutoMigrate(&model.Listener{})
	db.AutoMigrate(&model.Observation{})

	return &repository{
		db,
	}
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) AddObservation(observation model.Observation) error {
	var foundObservation model.Observation
	r.db.
		Where("registry = ? AND image = ? AND tag = ?",
			observation.Registry,
			observation.Image,
			observation.Tag).
		Find(&foundObservation)

	if foundObservation.Image != observation.Image {
		return r.db.Create(&observation).Error
	}

	tx := r.db.Begin()

	//to delete it permanently we use "Unscoped()"
	if tx.Unscoped().Delete(model.Listener{}, "name = ? AND observation_id = ?", observation.Listener[0].Name, foundObservation.ID).Error != nil {
		return tx.Rollback().Error
	}

	observation.Listener[0].ObservationID = foundObservation.ID
	if tx.Create(&observation.Listener[0]).Error != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
}

func (r *repository) GetObservations(since time.Duration) (*sql.Rows, error) {
	return r.db.
		Model(&model.Observation{}).
		Where(fmt.Sprintf("created_at = updated_at OR updated_at < NOW() + interval '%f second'", since.Seconds())).
		Order("updated_at ASC").
		Rows()
}

func (r *repository) NextObservation(rows *sql.Rows) (model.Observation, error) {
	var result model.Observation
	return result, r.db.ScanRows(rows, &result)
}

func (r *repository) UpdateImageHash(registry, image, tag, hash string) error {
	return r.db.
		Model(model.Observation{}).
		Where("registry = ? AND image = ? AND tag = ?",
			registry,
			image,
			tag).
		Update("last_hash", hash).
		Error
}

func (r *repository) TouchObservation(id uint) error {
	return r.db.
		Model(model.Observation{}).
		Where("id = ?", id).
		Update("updated_at", time.Now()).
		Error
}

func (r *repository) GetOverdueNotifications() (*sql.Rows, error) {
	return r.db.
		Raw(`SELECT l.*, o.last_hash AS observation_hash
   			FROM observations o, listeners l
   			WHERE l.observation_id = o.id AND (
     			(l.last_hash <> o.last_hash) OR
     			(l.last_hash IS NULL AND o.last_hash IS NOT NULL)
   			)`).
		Rows()
}

func (r *repository) NextNotification(rows *sql.Rows) (OverdueListener, error) {
	var result OverdueListener
	return result, r.db.ScanRows(rows, &result)
}

func (r *repository) UpdateListener(id uint, hash string) error {
	return r.db.
		Model(model.Listener{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_notification": time.Now(),
			"last_hash":         hash,
		}).
		Error
}
