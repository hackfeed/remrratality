package userrepo

import "github.com/hackfeed/remrratality/backend/internal/domain"

type UserRepository interface {
	AddUser(string, string) (domain.User, error)
	GetUser(string) (domain.User, error)
	UpdateUser(string, domain.User) error
}
