package main

import (
	"context"
	"fmt"
	"github.com/rainu/docker-image-watcher/internal/client"
	"github.com/rainu/docker-image-watcher/internal/config"
	"github.com/rainu/docker-image-watcher/internal/database"
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

	notificationJobs := make(chan worker.NotificationJob, cfg.NotificationLimit)
	for i := 0; i < cfg.NotificationLimit; i++ {
		go worker.NewNotifier(notificationJobs, dbRepo, httpClient).Do()
	}

	observationJobs := make(chan worker.ObservationJob, cfg.ObservationLimit)
	for i := 0; i < cfg.ObservationLimit; i++ {
		go worker.NewObserver(observationJobs, dbRepo, clients).Do()
	}

	go worker.NewObserverWorker(dbRepo, cfg.ObservationInterval, cfg.ObservationDispatchInterval, observationJobs, closeChan).Do()
	go worker.NewNotifyWorker(dbRepo, cfg.NotificationDispatchInterval, notificationJobs, closeChan).Do()
}
