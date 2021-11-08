package user

import (
	"context"
	"gorm.io/gorm"
	"notify-service/internal/models"
	"notify-service/internal/repo"
)

type DBRepo struct {
	db *gorm.DB
}

func NewUserDBRepo(db *gorm.DB) *DBRepo {
	return &DBRepo{db: db}
}

func (u *DBRepo) GetCurrentUserInfo(ctx context.Context, filter *Filter) (*models.User, error) {
	var adminMember *models.User
	if err := u.db.First(&adminMember, filter.ID).Error; err != nil {
		return adminMember, err
	}
	return adminMember, nil
}

func (u *DBRepo) GetAll(ctx context.Context, a, b string) (*models.UserPaginator, error) {
	var users []*models.User
	query := u.db.Where("username = ?", "asd")
	result, cursor, err := repo.Paginator(a, b, 1, "ASC").Paginate(query, &users)
	if err != nil {
		return nil, err
	}
	return &models.UserPaginator{Users: users, Cursor: cursor}, result.Error
}
