package hash

import (
	dto "github.com/Eanhain/gofermart/internal/api"
	"github.com/Eanhain/gofermart/internal/domain"
	"github.com/alexedwards/argon2id"
)

func CreateUserHash(log domain.Logger, user dto.UserInput) dto.User {
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		log.Warnln("Cannot create hash", user.Login, err)
		return dto.User{}
	}
	return dto.User{Login: user.Login, Hash: hash}
}

func VerifyUserHash(log domain.Logger, user dto.UserInput, tUser dto.User) bool {
	match, err := argon2id.ComparePasswordAndHash(user.Password, tUser.Hash)
	if err != nil {
		log.Warnln("Cannot verify hash", user.Login, err)
	}
	log.Infoln("User", user.Login, "match:", match)
	return match
}
