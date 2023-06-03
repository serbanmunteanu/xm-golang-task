package user

type UserRepository interface {
	Create(user *Model) error
	Read(email string) (*Model, error)
}
