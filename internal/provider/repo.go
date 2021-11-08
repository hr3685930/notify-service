package provider

import (
	"notify-service/internal/repo/user"
	"notify-service/pkg/db"
)

func Register() {
	user.Repo = user.NewUserDBRepo(db.Orm)
}
