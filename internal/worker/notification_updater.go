package worker

import (
	"github.com/rainu/docker-image-watcher/internal/database"
	log "github.com/sirupsen/logrus"
)

type NotificationUpdate struct {
	ListenerId   uint
	ListenerName string
	Hash         string
}

type notificationUpdater struct {
	db   database.Repository
	jobs chan NotificationUpdate
}

func NewNotificationUpdater(jobs chan NotificationUpdate, db database.Repository) Worker {
	return &notificationUpdater{
		db:   db,
		jobs: jobs,
	}
}

func (n *notificationUpdater) Do() {
	log.Info("Start notification updater...")
	defer log.Info("Stop notification updater...")

	for {
		notificationUpdate, ok := <-n.jobs
		if !ok {
			return
		}

		if err := n.db.UpdateListener(notificationUpdate.ListenerId, notificationUpdate.Hash); err != nil {
			log.Errorf("Unable to save notification state for listener %s. Error: %n",
				notificationUpdate.ListenerName, err)
		}
	}
}
