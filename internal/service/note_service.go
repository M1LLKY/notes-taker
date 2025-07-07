package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"notes-taker/internal/logger"
	"notes-taker/internal/mapper"
	"notes-taker/internal/models"
	"notes-taker/internal/repository"
	"notes-taker/internal/repository/postgres"
)

type NoteService struct {
	NoteRepository repository.NoteRepository
}

func NewNoteService(noteRepository repository.NoteRepository) NoteService {
	return NoteService{
		NoteRepository: noteRepository,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, input *CreateNote) (*models.CreateNoteResponse, error) {
	logger.Get().Info(ctx, "Создание новой заметки", logrus.Fields{
		"title":   input.Title,
		"content": input.Content,
	})

	noteID, err := s.NoteRepository.CreateNote(ctx, input.Title, input.Content, input.UserID)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при создании заметки", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Get().Info(ctx, "Заметка успешно создана", logrus.Fields{
		"note_id": noteID,
	})

	return &models.CreateNoteResponse{Data: &models.CreateNoteData{NoteID: noteID}}, nil
}

func (s *NoteService) GetAllNotes(ctx context.Context, userID int) (*models.NoteListResponse, error) {
	logger.Get().Info(ctx, "Получение всех заметок пользователя", nil)

	notes, err := s.NoteRepository.GetAllNotes(ctx, userID)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при получении заметок", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	noteDTOs := make([]models.NoteDTO, len(notes))
	for i, e := range notes {
		noteDTOs[i] = mapper.MapNoteDTOFromNoteDb(e)
	}

	return &models.NoteListResponse{Data: &noteDTOs}, nil
}

func (s *NoteService) GetNoteByID(ctx context.Context, userID, noteID int) (*models.NoteResponse, error) {
	logger.Get().Info(ctx, "Получение заметки по ID", logrus.Fields{
		"note_id": noteID,
	})

	note, err := s.NoteRepository.GetNoteByID(ctx, noteID)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при получении заметки", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	if note.UserID != userID {
		logger.Get().Error(ctx, "Попытка доступа к чужой заметке", logrus.Fields{
			"note_id":  noteID,
			"owner_id": note.UserID,
		})
		return nil, ErrNoteForbidden
	}

	noteDTO := mapper.MapNoteDTOFromNoteDb(*note)
	return &models.NoteResponse{Data: &noteDTO}, nil
}

func (s *NoteService) UpdateNoteByID(ctx context.Context, input *UpdateNote) (*models.OperationResultResponse, error) {
	note, err := s.NoteRepository.GetNoteByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, postgres.ErrNoteNotFound) {
			logger.Get().Info(ctx, "Заметка для обновления не найдена", logrus.Fields{
				"note_id": input.ID,
			})
			return nil, ErrNoteNotFound
		}
		logger.Get().Error(ctx, "Ошибка при получении заметки перед обновлением", logrus.Fields{
			"note_id": input.ID,
			"error":   err.Error(),
		})
		return nil, err
	}

	if note.UserID != input.UserID {
		logger.Get().Info(ctx, "Попытка обновить чужую заметку", logrus.Fields{
			"note_id":  input.ID,
			"owner_id": note.UserID,
		})
		return nil, ErrNoteForbidden
	}

	err = s.NoteRepository.UpdateNoteByID(ctx, input.Title, input.Content, input.ID)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при обновлении заметки", logrus.Fields{
			"note_id": input.ID,
			"error":   err.Error(),
		})
		return nil, err
	}

	return &models.OperationResultResponse{Data: &models.OperationResultData{Success: true}}, nil
}

func (s *NoteService) DeleteNoteByID(ctx context.Context, userID, noteID int) (*models.OperationResultResponse, error) {
	note, err := s.NoteRepository.GetNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, postgres.ErrNoteNotFound) {
			logger.Get().Error(ctx, "Заметка для удаления не найдена", logrus.Fields{
				"note_id": noteID,
			})
			return nil, ErrNoteNotFound
		}
		logger.Get().Error(ctx, "Ошибка при получении заметки перед удалением", logrus.Fields{
			"note_id": noteID,
			"error":   err.Error(),
		})
		return nil, err
	}

	if note.UserID != userID {
		logger.Get().Error(ctx, "Попытка удалить чужую заметку", logrus.Fields{
			"note_id":  noteID,
			"owner_id": note.UserID,
		})
		return nil, ErrNoteForbidden
	}

	logger.Get().Info(ctx, "Удаление заметки по ID", logrus.Fields{
		"note_id": noteID,
	})

	err = s.NoteRepository.DeleteNoteByID(ctx, noteID)
	if err != nil {
		logger.Get().Error(ctx, "Ошибка при удалении заметки", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.Get().Info(ctx, "Заметка успешно удалена", logrus.Fields{
		"note_id": noteID,
	})

	return &models.OperationResultResponse{Data: &models.OperationResultData{Success: true}}, nil
}
