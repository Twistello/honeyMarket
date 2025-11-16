package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Service interface {

	Register(ctx context.Context, email, password, role string) (*UserResponse, error)

	Authenticate(ctx context.Context, email, password string) (*UserResponse, error)

	GetByID(ctx context.Context, id int64) (*UserResponse, error)

	UpdateRole(ctx context.Context, id int64, newRole string) (*UserResponse, error)

	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, limit, offset int) ([]*UserResponse, error)

}

type UserService struct {
	repo *Repository
}

// эскпортированный конструктор структуры
func NewService(repo *Repository) *UserService {
	return &UserService{repo: repo}
}

// заворачивает user в безопасный userResponse
func toResponse(u *User) *UserResponse {
	user := &UserResponse{u.Id,u.Email,u.Role, u.CreatedAt}
	return user
}

// метод регистрации
func (s *UserService) Register(ctx context.Context, email, password, role string) (*UserResponse, error) {
	// 1. не зареган ли такой пользователь
	existing, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExists
	}
	// 2. хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	// 3. формируем обьект для передачи на слой ниже 
	user := &User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}
	// 4. пользователь улетает в слой репозитория
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	// 5.DTO без пароля возвращается на слой выше
	
	return toResponse(user), nil
}

//метод атуентификации
func (s *UserService) Authenticate(ctx context.Context, email, password string) (*UserResponse, error) {
	
	// 1. Проверка на наличие зареганной почты

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	// 2. сравнение полученного хеша из бд и введенного пароля
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}

	// 3.DTO без пароля возвращается на слой выше
	
	return toResponse(user), nil

}

func (s *UserService) GetByID(ctx context.Context, id int64) (*UserResponse, error) {
	// 1. Достаеп пользователя по id
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check user id: %w", err)
	}

	// 2.DTO без пароля возвращается на слой выше
	
	return toResponse(user), nil
}

func (s *UserService) UpdateRole(ctx context.Context, id int64, newRole string) (*UserResponse, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	user.Role = newRole

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return toResponse(user), nil
}


func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user  by id: %w", err)
	}

	return nil	
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]*UserResponse, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	resp := make([]*UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toResponse(u))
	}
	return resp, nil
}
