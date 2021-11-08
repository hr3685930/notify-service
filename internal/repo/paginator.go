package repo

import "github.com/pilagod/gorm-cursor-paginator/v2/paginator"

func Paginator(after, before string, limit int, order paginator.Order) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			Rules: []paginator.Rule{
				{
					Key: "ID",
				},
				{
					Key:     "CreatedAt",
					Order:   paginator.ASC,
					SQLRepr: "user.created_at",
				},
			},
			Limit:  limit,
			Order:  order,
			Before: before,
			After:  after,
		},
	}
	return paginator.New(opts...)
}
