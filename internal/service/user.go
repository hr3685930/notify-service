package service

import (
	"context"
	"notify-service/internal/models"
	"notify-service/internal/repo/user"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetAll(ctx context.Context, a, b string) (*models.UserPaginator, error) {

	//可以使用并发原语
	return user.Repo.GetAll(ctx, a, b)
}
