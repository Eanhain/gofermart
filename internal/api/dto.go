package dto

//go:generate easyjson -all .

//easyjson:json
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//easyjson:json
type UserArray []User

func ToDTO(username string, password string) User {
	return User{
		Login:    username,
		Password: password}
}
