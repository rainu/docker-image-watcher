package worker

import (
	"database/sql"
	"github.com/rainu/docker-image-watcher/internal/database"
	"github.com/rainu/docker-image-watcher/internal/database/model"
	log "github.com/sirupsen/logrus"
	"time"
)

type observerWorker struct {
	db              database.Repository
	jobChan         chan model.Observation
	closeChannel    chan interface{}
	interval        time.Duration
	maxObservations int
}

func NewObserverWorker(
	dbRepo database.Repository,
	interval time.Duration,
	jobChan chan model.Observation,
	closeChannel chan interface{}) Worker {

	return &observerWorker{
		db:           dbRepo,
		jobChan:      jobChan,
		closeChannel: closeChannel,
		interval:     interval,
	}
}

func (o *observerWorker) Do() {
	log.Info("Start observation dispatcher...")
	defer log.Info("Stop observation dispatcher...")
	defer close(o.jobChan)

	//first := true

	for {
		select {
		case <-o.closeChannel:
			return
		default:
		}

		//else {
		//	select {
		//	case <-time.After(o.interval):
		//	case <-o.closeChannel:
		//		return
		//	default:
		//	}
		//}

		rows, err := o.db.GetObservations(o.interval)
		if err != nil {
			log.Errorf("Could not get listeners. Error: %v", err)
			continue
		}

		o.processRows(rows, o.jobChan)
	}
}

func (o *observerWorker) processRows(rows *sql.Rows, jobs chan model.Observation) {
	defer rows.Close()

	for rows.Next() {
		select {
		case <-o.closeChannel:
			return
		default:
		}

		observation, err := o.db.NextObservation(rows)
		if err != nil {
			log.Errorf("Could not get observations. Error: %v", err)
			break
		}

		jobs <- observation
	}
}
