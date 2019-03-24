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

type notifier struct {
	jobs       chan database.OverdueListener
	updateChan chan NotificationUpdate
	httpClient *http.Client
}

func NewNotifier(
	jobs chan database.OverdueListener,
	updateChan chan NotificationUpdate,
	httpClient *http.Client) Worker {

	return &notifier{
		jobs:       jobs,
		updateChan: updateChan,
		httpClient: httpClient,
	}
}

func (n *notifier) Do() {
	log.Info("Start notifier...")
	defer log.Info("Stop notifier...")

	for {
		listener, ok := <-n.jobs
		if !ok {
			return
		}

		log.Infof("Notify listener %s", listener.Name)
		n.notify(listener)
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

	n.updateChan <- NotificationUpdate{
		ListenerId:   listener.ID,
		ListenerName: listener.Name,
		Hash:         listener.ObservationHash,
	}
}
