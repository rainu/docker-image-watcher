package worker

import (
	"github.com/rainu/docker-image-watcher/internal/client"
	"github.com/rainu/docker-image-watcher/internal/database/model"
	log "github.com/sirupsen/logrus"
)

type observer struct {
	jobs          chan model.Observation
	updateChan    chan ObservationUpdate
	dockerClients map[string]client.DockerRegistryClient
}

func NewObserver(
	jobs chan model.Observation,
	updateChan chan ObservationUpdate,
	dockerClients map[string]client.DockerRegistryClient) Worker {

	return &observer{
		jobs:          jobs,
		updateChan:    updateChan,
		dockerClients: dockerClients,
	}
}

func (o *observer) Do() {
	log.Info("Start observer...")
	defer log.Info("Stop observer...")

	for {
		observation, ok := <-o.jobs
		if !ok {
			return
		}

		log.Infof("Observe %s/%s:%s", observation.Registry, observation.Image, observation.Tag)
		o.observe(observation)
	}
}

func (o *observer) observe(observation model.Observation) {
	if dockerClient, ok := o.dockerClients[observation.Registry]; ok {
		log.Infof("Get manifest for %s/%s:%s", observation.Registry, observation.Image, observation.Tag)
		manifest, err := dockerClient.GetManifest(observation.Image, observation.Tag)

		if err != nil {
			log.Errorf("Error while getting Manifest for %s:%s. Error: %v", observation.Image, observation.Tag, err)
		} else {
			log.Infof("Got manifest for %s/%s:%s", observation.Registry, observation.Image, observation.Tag)
		}

		o.updateChan <- ObservationUpdate{
			Registry: observation.Registry,
			Image:    observation.Image,
			Tag:      observation.Tag,
			Hash:     manifest.Config.Digest,
		}
	} else {
		log.Warningf("No docker client found for registry: %s", observation.Registry)
	}

	//if err := o.db.TouchObservation(observation.ID); err != nil {
	//	log.Errorf("Unable to touch observation for %s/%s:%s. Error: %v",
	//		observation.Registry, observation.Image, observation.Tag, err)
	//}
}
