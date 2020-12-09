package service

import (
	"context"

	"github.com/summerKK/go-code-snippet-library/koel-api/internal/dao"
)

type service struct {
	ctx context.Context
	dao *dao.Dao
}
