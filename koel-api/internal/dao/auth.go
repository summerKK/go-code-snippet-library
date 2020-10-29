package dao

import "github.com/summerKK/go-code-snippet-library/koel-api/internal/model"

func (d *Dao) GetAuth(email string) (*model.User, error) {
	user := &model.User{Email: email}

	return user.GetUserByEmail(d.engine)
}
