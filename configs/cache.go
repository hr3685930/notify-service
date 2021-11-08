package configs

import "notify-service/pkg/config"

// Cache default once
type Cache struct {
	Redis Redis
	Sync  Sync
}

type Redis struct {
	config.RedisDrive
	Host     string `default:"127.0.0.1"`
	Port     string `default:"6379"`
	Database string `default:"0"`
	Auth     string `default:""`
}

type Sync struct {
	config.SyncDrive
}
