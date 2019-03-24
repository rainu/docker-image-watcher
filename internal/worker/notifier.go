package worker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/rainu/docker-image-watcher/internal/database"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type NotificationJob struct {
	Listener     database.OverdueListener
	FeedbackChan chan interface{}
}

type notifier struct {
	db         database.Repository
	jobs       chan NotificationJob
	httpClient *http.Client
}

func NewNotifier(jobs chan NotificationJob, db database.Repository, httpClient *http.Client) Worker {
	return &notifier{
		jobs:       jobs,
		db:         db,
		httpClient: httpClient,
	}
}

func (n *notifier) Do() {
	log.Info("Start notifier...")
	defer log.Info("Stop notifier...")

	for {
		job, ok := <-n.jobs
		if !ok {
			return
		}

		n.notify(job.Listener)
		job.FeedbackChan <- true
	}
}

func (n *notifier) notify(listener database.OverdueListener) {
	request, err := http.NewRequest(listener.Method, listener.Url, nil)
	if err != nil {
		log.Errorf("Unable to create request for listener %s. Error: %n",
			listener.Name, err)
		return
	}

	if len(listener.Body) > 0 {
		request.Body = ioutil.NopCloser(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(listener.Body)))
	}

	if len(listener.Header) > 0 {
		header := make(map[string]string)
		if err := json.Unmarshal(listener.Header, &header); err != nil {
			log.Errorf("Unable to unmarshal header for listener %s. Error: %n",
				listener.Name, err)
			return
		}

		for key, value := range header {
			request.Header.Set(key, value)
		}
	}

	response, err := n.httpClient.Do(request)
	if err != nil {
		log.Errorf("Unable to notify listener %s. Error: %n",
			listener.Name, err)
		return
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	if err := n.db.UpdateListener(listener.ID, listener.ObservationHash); err != nil {
		log.Errorf("Unable to save notification state for listener %s. Error: %n",
			listener.Name, err)
	}
}
