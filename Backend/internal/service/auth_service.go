// auth_service - creating new auth service , login and Register the user

package service

import (
	"github.com/Suthar345Piyush/invoicego/internal/config"
	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/Suthar345Piyush/invoicego/internal/util"
)

type AuthService struct {
	userService *UserService
	jwtConfig   *config.JWTConfig
}

// function for new auth service and return auth service

func NewAuthService(userService *UserService, jwtConfig *config.JWTConfig) *AuthService {
	return &AuthService{
		userService: userService,
		jwtConfig:   jwtConfig,
	}

}

// user registration , taking register request as input and returns login response as output

func (s *AuthService) Register(req *domain.RegisterRequest) (*domain.LoginResponse, error) {

	// firstly validating the input

	if err := util.ValidateStruct(req); err != nil {
		return nil, domain.ErrInvalidInput
	}

	// after this creating a user

	user, err := s.userService.CreateUser(req)
	if err != nil {
		return nil, err
	}

	// generating tokens

	// access token

	accessToken, err := util.GenerateAccessToken(
		user.ID,
		user.Email,
		s.jwtConfig.Secret,
		s.jwtConfig.AccessExpiry,
	)

	if err != nil {
		return nil, err
	}

	// refresh token

	refreshToken, err := util.GenerateRefreshToken(
		user.ID,
		s.jwtConfig.Secret,
		s.jwtConfig.RefreshExpiry,
	)

	if err != nil {
		return nil, err
	}

	// at last returning the login response

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil

}

// login process same taking login request and returning login response

func (s *AuthService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {

}
