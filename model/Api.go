package model

import (
	"time"
)

type DeleteResponse struct {
	IgnoredKeys []string `json:"ignored"`
}

type SecretsVersionKeys struct {
	VersionCreatedAt time.Time `json:"timestamp"`
	ID               string    `json:"id"`
	Keys             []string  `json:"keys"`
}
