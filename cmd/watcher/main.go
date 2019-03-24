package main

import (
	"context"
	"fmt"
	"github.com/rainu/docker-image-watcher/internal/client"
	"github.com/rainu/docker-image-watcher/internal/config"
	"github.com/rainu/docker-image-watcher/internal/database"
	"github.com/rainu/docker-image-watcher/internal/database/model"
	"github.com/rainu/docker-image-watcher/internal/server"
	"github.com/rainu/docker-image-watcher/internal/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var closeChan chan interface{}
var cfg *config.Config

func main() {
	if cfg == nil {
		cfg = config.NewConfig()
	}

	//database
	dbRepo := database.NewRepository(cfg.DatabaseInfo())
	defer dbRepo.Close()

	httpServer := startServer(dbRepo, cfg)
	defer func() {
		//gracefully shutdown the httpServer
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		httpServer.Shutdown(ctx)
	}()

	//clients
	registryHttpClient := http.DefaultClient
	notifierHttpClient := &http.Client{
		Timeout: cfg.NotificationTimeout,
	}
	registryClients := make(map[string]client.DockerRegistryClient)
	registryClients["docker.io"] = client.NewDockerRegistryClient(registryHttpClient)

	//wait for interruption
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	closeChan = make(chan interface{})

	//worker
	startWorker(dbRepo, registryClients, notifierHttpClient, closeChan)

	<-stop
	close(closeChan)
}

func startServer(dbRepo database.Repository, cfg *config.Config) *http.Server {
	router := server.NewRouter(dbRepo)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%v", cfg.BindPort), Handler: router}

	go func() {
		httpServer.ListenAndServe()
	}()

	return httpServer
}

func startWorker(
	dbRepo database.Repository,
	clients map[string]client.DockerRegistryClient,
	httpClient *http.Client,
	closeChan chan interface{}) {

	notifierJobs := make(chan database.OverdueListener, cfg.NotificationLimit)
	notifierUpdateJobs := make(chan worker.NotificationUpdate)
	go worker.NewNotifier(notifierJobs, notifierUpdateJobs, httpClient).Do()
	go worker.NewNotificationUpdater(notifierUpdateJobs, dbRepo).Do()

	observationJobs := make(chan model.Observation, cfg.ObservationLimit)
	observationUpdateJobs := make(chan worker.ObservationUpdate)
	go worker.NewObserver(observationJobs, observationUpdateJobs, clients).Do()
	go worker.NewObservationUpdater(observationUpdateJobs, dbRepo).Do()

	go worker.NewObserverWorker(dbRepo, cfg.ObservationInterval, observationJobs, closeChan).Do()
	go worker.NewNotifyWorker(dbRepo, cfg.NotificationInterval, notifierJobs, closeChan).Do()
}
