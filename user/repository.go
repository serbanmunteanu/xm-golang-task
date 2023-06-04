package user

type IRepository interface {
	Create(user *Model) error
	Read(email string) (*Model, error)
}
