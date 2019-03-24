package config

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"time"
)

type Config struct {
	databaseConfig

	BindPort int `arg:"--bind-port,env:BIND_PORT,help:The HTTP server bind port."`

	ObservationInterval  time.Duration `arg:"--observation-interval,env:OBSERVATION_INTERVAL,help:The interval for checking the observations."`
	ObservationLimit     int           `arg:"--observation-limit,env:OBSERVATION_LIMIT,help:How many observations should be check simultaneously."`
	NotificationInterval time.Duration `arg:"--notification-interval,env:NOTIFICATION_INTERVAL,help:The interval for notify the overdue listeners."`
	NotificationLimit    int           `arg:"--notification-limit,env:NOTIFICATION_LIMIT,help:How many listeners should be notify simultaneously."`
	NotificationTimeout  time.Duration `arg:"--notification-timeout,env:NOTIFICATION_TIMEOUT,help:The connection timeout for each listener."`
}

type databaseConfig struct {
	DatabaseSSL      string `arg:"--database-ssl,env:DATABASE_SSL,help:The database ssl mode: enable/disable."`
	DatabaseHost     string `arg:"--database-host,env:DATABASE_HOST,help:The database host."`
	DatabasePort     int    `arg:"--database-port,env:DATABASE_PORT,help:The database port."`
	DatabaseSchema   string `arg:"--database-schema,env:DATABASE_SCHEMA,help:The database schema."`
	DatabaseUser     string `arg:"--database-user,env:DATABASE_USER,help:The database user."`
	DatabasePassword string `arg:"--database-password,env:DATABASE_PASSWORD,help:The database password."`
}

func NewConfig() *Config {
	cfg := &Config{
		BindPort: 8080,
		databaseConfig: databaseConfig{
			DatabaseSSL:      "disable",
			DatabaseHost:     "localhost",
			DatabasePort:     5432,
			DatabaseSchema:   "postgres",
			DatabaseUser:     "postgres",
			DatabasePassword: "postgres",
		},

		ObservationInterval:  1 * time.Minute,
		ObservationLimit:     10,
		NotificationInterval: 1 * time.Minute,
		NotificationLimit:    60,
		NotificationTimeout:  10 * time.Second,
	}

	arg.Parse(cfg)

	return cfg
}

func (c *databaseConfig) DatabaseInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseSchema,
		c.DatabaseSSL,
	)
}
