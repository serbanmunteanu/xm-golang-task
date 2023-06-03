package di

import (
	companyRepo "github.com/serbanmunteanu/xm-golang-task/company"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	userRepo "github.com/serbanmunteanu/xm-golang-task/user"
)

type DI struct {
	UserRepository    userRepo.UserRepository
	CompanyRepository companyRepo.CompanyRepository
	Jwt               jwt.Jwt
}
