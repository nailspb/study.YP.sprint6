package mem

import (
	"cmp"
	"errors"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/models"
	"slices"
)

type Storage struct {
	data map[string]models.Task
}

func New() *Storage {
	s := &Storage{}
	s.data = map[string]models.Task{
		"1": {
			ID:          "1",
			Description: "Сделать финальное задание темы REST API",
			Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
			Applications: []string{
				"VS Code",
				"Terminal",
				"git",
			},
		},
		"2": {
			ID:          "2",
			Description: "Протестировать финальное задание с помощью Postmen",
			Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
			Applications: []string{
				"VS Code",
				"Terminal",
				"git",
				"Postman",
			},
		},
	}
	return s
}

func (s *Storage) GetAllTasks() []models.Task {
	tasks := make([]models.Task, 0, len(s.data))
	for _, t := range s.data {
		tasks = append(tasks, t)
	}
	slices.SortFunc(tasks, func(a, b models.Task) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return tasks
}

func (s *Storage) GetTask(id string) *models.Task {
	if val, ok := s.data[id]; ok {
		return &val
	}
	return nil
}

func (s *Storage) DeleteTask(id string) bool {
	if _, ok := s.data[id]; ok {
		delete(s.data, id)
		return true
	}
	return false
}

func (s *Storage) AddTask(task models.Task) error {
	if _, ok := s.data[task.ID]; ok {
		return errors.New("task already exists")
	}
	s.data[task.ID] = task
	return nil
}
