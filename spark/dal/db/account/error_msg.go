package account

import "errors"

var (
	DbUserExists = errors.New("用户已存在")
)
