package user

type User struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string
}
