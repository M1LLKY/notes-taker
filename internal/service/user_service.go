package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"notes-taker/internal/config"
	"notes-taker/internal/logger"
	"notes-taker/internal/models"
	"notes-taker/internal/repository"
	"notes-taker/pkg/auth"
	"strconv"
	"strings"
)

type UserService struct {
	UserRepository repository.UserRepository
	TokenManager   auth.TokenManager
	Config         *config.Config
}

func NewUserService(userRepository repository.UserRepository, manager auth.TokenManager, config *config.Config) UserService {
	return UserService{
		UserRepository: userRepository,
		TokenManager:   manager,
		Config:         config,
	}
}

func (s *UserService) SignIn(ctx context.Context, input SignIn) (*models.AuthResponse, error) {
	logger.Get().Info(ctx, "Регистрация нового пользователя", logrus.Fields{
		"username": strings.ToLower(input.Username),
	})

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при хэшировании пароля", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	userID, err := s.UserRepository.CreateUser(ctx, strings.ToLower(input.Username), string(hash))
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при создании пользователя", logrus.Fields{
			"username": strings.ToLower(input.Username),
			"error":    err.Error(),
		})
		return nil, err
	}

	jwt, err := s.TokenManager.NewJWT(strconv.Itoa(userID), s.Config.JwtTTLDuration)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при генерации JWT", logrus.Fields{
			"username": strings.ToLower(input.Username),
			"error":    err.Error(),
		})
		return nil, err
	}

	logger.Get().Info(ctx, "Пользователь успешно зарегестрирован", logrus.Fields{
		"username": strings.ToLower(input.Username),
	})

	return &models.AuthResponse{
		Data: &models.AuthData{
			Username:    strings.ToLower(input.Username),
			AccessToken: jwt,
		},
	}, nil
}
