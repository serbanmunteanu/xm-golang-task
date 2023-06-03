package di

import (
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/user/repository"
)

type DI struct {
	UserRepository repository.UserRepository
	Jwt            jwt.Jwt
}
