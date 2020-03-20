package common

type UserInfo struct {
	User     string `json:"user"`
	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
