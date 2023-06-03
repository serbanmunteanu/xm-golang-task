package repository

import "github.com/serbanmunteanu/xm-golang-task/user"

type UserRepository interface {
	Create(user *user.UserDbModel) error
	Read(email string) (*user.UserDbModel, error)
}
