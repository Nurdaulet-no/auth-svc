package usecase

import (
	"time"

	"github.com/Nurdaulet-no/auth-svc/internal/domain"
	"github.com/Nurdaulet-no/auth-svc/internal/repository"
	"github.com/Nurdaulet-no/auth-svc/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users repository.UserRepository
	jwt *jwt.Manager
	idGen func() string
	now func() time.Time
}

func NewAuthService(users repository.UserRepository, jwtManager *jwt.Manager, idGen func() string) *AuthService {
	return &AuthService{
		users: users,
		jwt: jwtManager,
		idGen: idGen,
		now: time.Now,
	}
}

func (s *AuthService) Register(email, password string) (domain.User, error){
	_, err := s.users.FindByEmail(email)
	if err == nil {
		return  domain.User{}, domain.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte (password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	u := domain.User{
		ID: s.idGen(),
		Email: email,
		PasswordHash: string(hash),
		CreatedAt: s.now(),
	}

	if err := s.users.Create(u); err != nil {
		return domain.User{}, err
	}

	return u, nil
} 

func (s *AuthService) Login (email, password string) (string, error){
	u, err := s.users.FindByEmail(email)
	if err != nil {
		return "", domain.ErrInvalidCreadiantals
	}

	if err := bcrypt .CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCreadiantals
	}

	return s.jwt.Issue(u.ID)
}

func (s *AuthService) Me (userId string) (domain.User, error){
	return s.users.FindByID(userId)
}

