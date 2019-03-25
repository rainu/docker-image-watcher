package worker

import (
	"github.com/rainu/docker-image-watcher/internal/client"
	"github.com/rainu/docker-image-watcher/internal/database"
	"github.com/rainu/docker-image-watcher/internal/database/model"
	log "github.com/sirupsen/logrus"
)

type ObservationJob struct {
	Observation  model.Observation
	FeedbackChan chan interface{}
}

type observer struct {
	db            database.Repository
	jobs          chan ObservationJob
	dockerClients map[string]client.DockerRegistryClient
}

func NewObserver(jobs chan ObservationJob, db database.Repository, dockerClients map[string]client.DockerRegistryClient) Worker {
	return &observer{
		jobs:          jobs,
		db:            db,
		dockerClients: dockerClients,
	}
}

func (o *observer) Do() {
	log.Info("Start observer...")
	defer log.Info("Stop observer...")

	for {
		job, ok := <-o.jobs
		if !ok {
			return
		}

		log.Infof("Observe %s/%s:%s", job.Observation.Registry, job.Observation.Image, job.Observation.Tag)
		o.observe(job.Observation)
		job.FeedbackChan <- true
	}
}

func (o *observer) observe(observation model.Observation) {
	if dockerClient, ok := o.dockerClients[observation.Registry]; ok {
		log.Infof("Get manifest for %s/%s:%s", observation.Registry, observation.Image, observation.Tag)
		manifest, err := dockerClient.GetManifest(observation.Image, observation.Tag)

		if err != nil {
			log.WithError(err).Errorf("Error while getting Manifest for %s:%s.", observation.Image, observation.Tag)
		} else {
			log.Infof("Got manifest for %s/%s:%s", observation.Registry, observation.Image, observation.Tag)
		}

		err = o.db.UpdateImageHash(
			observation.Registry,
			observation.Image,
			observation.Tag,
			manifest.Config.Digest)

		if err != nil {
			log.WithError(err).Errorf("Unable to update docker image hash for %s/%s:%s.",
				observation.Registry, observation.Image, observation.Tag)
		}
	} else {
		log.Warningf("No docker client found for registry: %s", observation.Registry)
	}

	if err := o.db.TouchObservation(observation.ID); err != nil {
		log.WithError(err).Errorf("Unable to touch observation for %s/%s:%s.",
			observation.Registry, observation.Image, observation.Tag)
	}
}
