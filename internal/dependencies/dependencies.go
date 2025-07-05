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
	NoteRepository repository.NoteRepository
	UserService    service.UserService
	NoteService    service.NoteService
	TokenManager   auth.TokenManager
}

func New(
	config *config.Config,
	userRepository repository.UserRepository,
	noteRepository repository.NoteRepository,
	userService service.UserService,
	noteService service.NoteService,
	manager auth.TokenManager) *Dependencies {
	return &Dependencies{
		Config:         config,
		UserRepository: userRepository,
		NoteRepository: noteRepository,
		UserService:    userService,
		NoteService:    noteService,
		TokenManager:   manager,
	}
}
