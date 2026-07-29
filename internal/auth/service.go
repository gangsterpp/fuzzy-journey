package auth

import (
	"context"
	"fmt"

	"github.com/gangsterpp/fuzzy-journey/internal/response"
	user "github.com/gangsterpp/fuzzy-journey/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	Generate(userID string) (string, error)
}
type AuthService interface {
	Register(ctx context.Context, request *AuthModel) (*user.AuthResponse, error)
	Login(ctx context.Context, request *AuthModel) (*user.AuthResponse, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type Service struct {
	authrepository Repository
	tokenManager   TokenManager
}

func (s *Service) Register(ctx context.Context, request *AuthModel) (*user.AuthResponse, error) {

	u, error := s.authrepository.FindUserByEmail(ctx, request.Email)
	if error != nil {
		return nil, error
	}
	if u != nil {
		return nil, response.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	passwordHash := string(hash)

	r, error := s.authrepository.RegisterUser(ctx, request.Email, passwordHash)
	if error != nil {
		return nil, error
	}
	tokenString, err := s.tokenManager.Generate(
		r.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &user.AuthResponse{
		User:  user.UserToUserRespose(r),
		Token: tokenString,
	}, nil
}

func (s *Service) Login(ctx context.Context, request *AuthModel) (*user.AuthResponse, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	passwordHash := string(hash)

	u, error := s.authrepository.FindUserByEmail(ctx, request.Email)
	if error != nil {
		return nil, error
	}

	if u.PasswordHash != passwordHash {
		return nil, response.ErrInvalidCredentials
	}
	tokenString, err := s.tokenManager.Generate(
		u.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &user.AuthResponse{
		User:  user.UserToUserRespose(u),
		Token: tokenString,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id string) (bool, error) {

	error := s.authrepository.DeleteByID(ctx, id)
	if error != nil {
		return false, error
	}
	return true, nil
}

var _ AuthService = (*Service)(nil)

func CreateAuthService(authrepository Repository, tokenManager TokenManager) AuthService {
	return &Service{authrepository: authrepository, tokenManager: tokenManager}
}
