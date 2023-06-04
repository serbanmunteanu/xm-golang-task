package di

import (
	companyRepo "github.com/serbanmunteanu/xm-golang-task/company"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/kafka"
	userRepo "github.com/serbanmunteanu/xm-golang-task/user"
)

type DI struct {
	UserRepository    userRepo.IRepository
	CompanyRepository companyRepo.IRepository
	Jwt               jwt.Jwt
	Producer          kafka.Producer
}
