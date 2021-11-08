package user

import (
	"notify-service/internal/types"
)

type Filter struct {
	ID int64
	Name string
	Pass string
}

func NewFilter(request *types.TestRequest) *Filter {
	return &Filter{
		ID: request.ID,
		Name: request.Name,
		Pass: request.Pass,
	}
}