package service

import (
	"context"

	"github.com/summerKK/go-code-snippet-library/koel-api/global"
	"github.com/summerKK/go-code-snippet-library/koel-api/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
		dao: dao.New(global.DBEngine),
	}
}
