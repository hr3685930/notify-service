package user

import (
	"context"
	"notify-service/internal/models"
)

var Repo Repository

type Repository interface {
	GetCurrentUserInfo(ctx context.Context, filter *Filter) (*models.User, error)
	GetAll(ctx context.Context, a, b string) (*models.UserPaginator, error)
}
