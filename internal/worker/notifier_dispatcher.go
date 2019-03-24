package worker

import (
	"database/sql"
	"github.com/rainu/docker-image-watcher/internal/database"
	log "github.com/sirupsen/logrus"
	"time"
)

type notifyWorker struct {
	db           database.Repository
	jobChan      chan database.OverdueListener
	closeChannel chan interface{}
	interval     time.Duration
}

func NewNotifyWorker(
	dbRepo database.Repository,
	interval time.Duration,
	jobs chan database.OverdueListener,
	closeChannel chan interface{}) Worker {

	return &notifyWorker{
		db:           dbRepo,
		jobChan:      jobs,
		closeChannel: closeChannel,
		interval:     interval,
	}
}

func (n *notifyWorker) Do() {
	log.Info("Start notify dispatcher...")
	defer log.Info("Stop notify dispatcher...")
	defer close(n.jobChan)

	first := true

	for {
		if first {
			first = false

			select {
			case <-n.closeChannel:
				return
			default:
			}
		} else {
			select {
			case <-time.After(n.interval):
			case <-n.closeChannel:
				return
			default:
			}
		}

		rows, err := n.db.GetOverdueNotifications()
		if err != nil {
			log.Errorf("Could not get listeners. Error: %v", err)
			continue
		}

		n.processRows(rows, n.jobChan)
	}
}

func (n *notifyWorker) processRows(rows *sql.Rows, jobs chan database.OverdueListener) {
	defer rows.Close()

	for rows.Next() {
		select {
		case <-n.closeChannel:
			return
		default:
		}

		listener, err := n.db.NextNotification(rows)
		if err != nil {
			log.Errorf("Could not get listeners. Error: %v", err)
			break
		}

		jobs <- listener
	}
}
