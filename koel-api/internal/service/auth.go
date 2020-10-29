package service

import (
	"errors"

	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/security"
)

type UserRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (s *Service) CheckAuth(param *UserRequest) error {
	var err error
	user, err := s.dao.GetAuth(param.Email)
	if err != nil {
		return err
	}

	if security.VerifyPassword(user.Password, param.Password) {
		return nil
	}

	return errors.New("check auth failed.")
}
