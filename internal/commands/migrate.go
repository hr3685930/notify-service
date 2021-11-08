package commands

import (
	"github.com/urfave/cli"
	"notify-service/internal/models"
	"notify-service/pkg/db"
)

func Migrate(c *cli.Context) {
	println("migrate start")
	_ = db.Orm.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&models.User{})
	println("migrate end")
}
