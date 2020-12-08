package service

import (
	"errors"

	"github.com/summerKK/go-code-snippet-library/koel-api/internal/dto"
	"github.com/summerKK/go-code-snippet-library/koel-api/pkg/security"
)

func (s *Service) CheckAuth(param *dto.UserRequest) error {
	var err error
	user, err := s.dao.GetAuth(param.UserName)
	if err != nil {
		return err
	}

	if security.VerifyPassword(user.Password, param.Password) {
		return nil
	}

	return errors.New("check auth failed.")
}
