package memory

import (
	"sync"

	"github.com/Nurdaulet-no/auth-svc/internal/domain"
	"github.com/Nurdaulet-no/auth-svc/internal/repository"
)

type UserRepo struct{
	mu sync.RWMutex
	byID map[string]domain.User
	byEmail map[string]string
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		byID: make(map[string]domain.User),
		byEmail: make(map[string]string),
	}
}

var _ repository.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[u.Email]; exists {
		return domain.ErrEmailTaken
	}

	r.byID[u.ID] = u
	r.byEmail[u.Email] = u.ID

	return nil
}

func (r *UserRepo) FindByEmail(email string) (domain.User, error){
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[email]

	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}

	return r.byID[id], nil
}

func (r *UserRepo) FindByID(id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byID[id]
	if !ok{
		return  domain.User{}, domain.ErrUserNotFound
	}

	return u, nil
}