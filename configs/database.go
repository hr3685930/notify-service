package configs

import (
	"notify-service/pkg/config"
)

// default once
type Database struct {
	MYSQL MYSQL
}

// MYSQL config
type MYSQL struct {
	config.MYSQLDrive
	Host     string `default:"127.0.0.1"`
	Port     string `default:"3306"`
	Database string `default:"notify"`
	Username string `default:"admin"`
	Password string `default:"123456"`
}
