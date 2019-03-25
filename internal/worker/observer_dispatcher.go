package worker

import (
	"database/sql"
	"github.com/rainu/docker-image-watcher/internal/database"
	log "github.com/sirupsen/logrus"
	"time"
)

type observerWorker struct {
	db               database.Repository
	jobChan          chan ObservationJob
	closeChannel     chan interface{}
	lookupInterval   time.Duration
	dispatchInterval time.Duration
	maxObservations  int
}

func NewObserverWorker(
	dbRepo database.Repository,
	lookupInterval time.Duration,
	dispatchInterval time.Duration,
	jobChan chan ObservationJob,
	closeChannel chan interface{}) Worker {

	return &observerWorker{
		db:               dbRepo,
		jobChan:          jobChan,
		closeChannel:     closeChannel,
		lookupInterval:   lookupInterval,
		dispatchInterval: dispatchInterval,
	}
}

func (o *observerWorker) Do() {
	log.Info("Start observation dispatcher...")
	defer log.Info("Stop observation dispatcher...")
	defer close(o.jobChan)

	first := true

	for {
		if first {
			first = false

			select {
			case <-o.closeChannel:
				return
			default:
			}
		} else {
			select {
			case <-time.After(o.dispatchInterval):
			case <-o.closeChannel:
				return
			}
		}

		rows, err := o.db.GetObservations(o.lookupInterval)
		if err != nil {
			log.WithError(err).Error("Could not get listeners.")
			continue
		}

		o.processRows(rows, o.jobChan)
	}
}

func (o *observerWorker) processRows(rows *sql.Rows, jobs chan ObservationJob) {
	defer rows.Close()

	feedbackChannels := make([]chan interface{}, 0)

	for rows.Next() {
		select {
		case <-o.closeChannel:
			return
		default:
		}

		observation, err := o.db.NextObservation(rows)
		if err != nil {
			log.WithError(err).Error("Could not get observations.")
			break
		}

		feedbackChannels = append(feedbackChannels, make(chan interface{}))
		jobs <- ObservationJob{
			Observation:  observation,
			FeedbackChan: feedbackChannels[len(feedbackChannels)-1],
		}
	}

	//wait for feedbacks
	for _, feedbackChan := range feedbackChannels {
		select {
		case <-feedbackChan:
		case <-o.closeChannel:
			return
		}
	}
}
