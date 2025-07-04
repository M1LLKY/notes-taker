package dependencies

import (
	"notes-taker/internal/config"
	"notes-taker/internal/repository"
	"notes-taker/internal/service"
	"notes-taker/pkg/auth"
)

type Dependencies struct {
	Config         *config.Config
	UserRepository repository.UserRepository
	//TaskRepository repository.TaskRepository
	UserService service.UserService
	//TaskService    service.TaskService
	TokenManager auth.TokenManager
}

func New(
	config *config.Config,
	userRepository repository.UserRepository,
	//taskRepository repository.TaskRepository,
	userService service.UserService,
	//taskService service.TaskService,
	manager auth.TokenManager) *Dependencies {
	return &Dependencies{
		Config:         config,
		UserRepository: userRepository,
		//TaskRepository: taskRepository,
		UserService: userService,
		//TaskService:    taskService,
		TokenManager: manager,
	}
}
