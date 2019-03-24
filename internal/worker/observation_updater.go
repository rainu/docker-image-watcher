package worker

import (
	"github.com/rainu/docker-image-watcher/internal/database"
	log "github.com/sirupsen/logrus"
)

type ObservationUpdate struct {
	Registry string
	Image    string
	Tag      string
	Hash     string
}

type observationUpdater struct {
	db   database.Repository
	jobs chan ObservationUpdate
}

func NewObservationUpdater(jobs chan ObservationUpdate, db database.Repository) Worker {
	return &observationUpdater{
		jobs: jobs,
		db:   db,
	}
}

func (o *observationUpdater) Do() {
	log.Info("Start observation updater...")
	defer log.Info("Stop observation updater...")

	for {
		observationUpdate, ok := <-o.jobs
		if !ok {
			return
		}

		err := o.db.UpdateImageHash(
			observationUpdate.Registry,
			observationUpdate.Image,
			observationUpdate.Tag,
			observationUpdate.Hash)

		if err != nil {
			log.Errorf("Unable to update docker image hash for %s/%s:%s. Error: %v",
				observationUpdate.Registry, observationUpdate.Image, observationUpdate.Tag, err)
		}
	}
}
