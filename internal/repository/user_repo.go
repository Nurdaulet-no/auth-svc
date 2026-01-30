package repository


import "github.com/Nurdaulet-no/auth-svc/internal/domain"

type UserRepository interface {
	Create(u domain.User) error
	FindByEmail(email string) (domain.User, error)
	FindByID(id string) (domain.User, error)
}
