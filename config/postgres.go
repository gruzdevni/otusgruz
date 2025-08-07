package config

import "time"

type Postgres struct {
	DSN   string `envconfig:"POSTGRES_DSN" default:"postgres://postgres@localhost:5432/?sslmode=disable"`
	DSNRO string `envconfig:"POSTGRES_DSN_RO" default:"postgres://postgres@localhost:5432/?sslmode=disable"`

	MaxOpenConns    int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns    int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"7"`
	ConnMaxLifetime time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"30m"`
}
