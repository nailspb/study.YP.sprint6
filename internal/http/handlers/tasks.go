package handlers

import (
	"encoding/json"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/models"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/helpers/slogHelper"
	"io"
	"log/slog"
	"net/http"
)

type Storage interface {
	GetAllTasks() []models.Task
	AddTask(task models.Task) error
	GetTask(id string) *models.Task
	DeleteTask(id string) bool
}

const (
	ErrorRequest       = "Error on read request"
	ErrorParseJson     = "Error on json parse"
	ErrorSerializeJson = "Error on serialize json"
	ErrorSaveTask      = "Error on save task"
	ErrorWriteResponse = "Error on write response"
	InfoNotFound       = "Task not found"
)

func AddTask(log *slog.Logger, storage Storage) http.HandlerFunc {
	log = log.With(slog.String("Handler", "AddTask()"))
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error(ErrorRequest, slogHelper.GetErrAttr(err))
		} else {
			defer r.Body.Close()
			task := models.Task{}
			err = json.Unmarshal(body, &task)
			if err != nil {
				log.Error(ErrorParseJson, slogHelper.GetErrAttr(err))
			} else {
				if err = storage.AddTask(task); err != nil {
					log.Error(ErrorSaveTask, slogHelper.GetErrAttr(err))
				} else {
					w.WriteHeader(http.StatusCreated)
					return
				}
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetAllTask(log *slog.Logger, storage Storage) http.HandlerFunc {
	log = log.With(slog.String("Handler", "GetAllTask()"))
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(storage.GetAllTasks())
		if err != nil {
			log.Error(ErrorSerializeJson, slogHelper.GetErrAttr(err))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(data)
			if err != nil {
				log.Error(ErrorWriteResponse, slogHelper.GetErrAttr(err))
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetTask(log *slog.Logger, storage Storage) http.HandlerFunc {
	log = log.With(slog.String("Handler", "GetTask()"))
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if task := storage.GetTask(id); task == nil {
			log.Info(InfoNotFound, slog.String("id", id))
		} else {
			data, err := json.Marshal(task)
			if err != nil {
				log.Error(ErrorSerializeJson, slogHelper.GetErrAttr(err))
			} else {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, err = w.Write(data)
				if err != nil {
					log.Error(ErrorWriteResponse, slogHelper.GetErrAttr(err))
				}
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func DeleteTask(log *slog.Logger, storage Storage) http.HandlerFunc {
	log = log.With(slog.String("Handler", "DeleteTask()"))
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if storage.DeleteTask(id) == false {
			log.Info(InfoNotFound, slog.String("id", id))
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
