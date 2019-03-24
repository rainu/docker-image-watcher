package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/rainu/docker-image-watcher/internal/database"
	"github.com/rainu/docker-image-watcher/internal/database/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type registryHandler struct {
	Repository database.Repository
}

type addObservationRequest struct {
	Registry string `json:"registry"`
	Image    string `json:"image"`
	Tag      string `json:"tag"`

	Trigger addObservationTrigger `json:"trigger"`
}

type addObservationTrigger struct {
	Name   string            `json:"name"`
	Method string            `json:"method"`
	Url    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

func NewAddObservationHandler(repo database.Repository) http.Handler {
	return &registryHandler{
		Repository: repo,
	}
}

func (r *registryHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	parsedBody := &addObservationRequest{
		Registry: "docker.io",
		Tag:      "latest",
		Trigger: addObservationTrigger{
			Method: "GET",
		},
	}
	if err := json.NewDecoder(request.Body).Decode(parsedBody); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if parsedBody.Image == "" || parsedBody.Trigger.Url == "" || parsedBody.Trigger.Name == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	observation := model.Observation{
		Registry: parsedBody.Registry,
		Image:    parsedBody.Image,
		Tag:      parsedBody.Tag,
		Listener: []model.Listener{{
			Name:   parsedBody.Trigger.Name,
			Method: parsedBody.Trigger.Method,
			Url:    parsedBody.Trigger.Url,
		}},
	}

	if parsedBody.Trigger.Header != nil {
		rawHeader, err := json.Marshal(parsedBody.Trigger.Header)
		if err != nil {
			log.Errorf("Error while marshall trigger-header: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		observation.Listener[0].Header = rawHeader
	}

	if _, err := base64.StdEncoding.Decode(observation.Listener[0].Body, []byte(parsedBody.Trigger.Body)); err != nil {
		log.Errorf("Error while decode base64 trigger-body: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := r.Repository.AddObservation(observation); err != nil {
		log.Errorf("Error while persisting observation: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}
