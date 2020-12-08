package dao

import "github.com/summerKK/go-code-snippet-library/koel-api/internal/model"

func (d *Dao) GetAuth(email string) (*model.UmsAdmin, error) {
	user := &model.UmsAdmin{Username: email}

	return user.GetUserByName(d.engine)
}
