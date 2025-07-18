package taskstore

import (
	"github.com/bogomollov/15.07.2025/internal/transport/response"
)

const (
	TaskStatusCreated    = "created"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
)

type Task struct {
	ID           int      `json:"id"`
	Status       string   `json:"status" binding:"oneof=created processing completed"`
	Download_URL string   `json:"download_url,omitempty"`
	Links        []string `json:"links,omitempty"`
}

type TaskStore struct {
	tasks map[int]Task
}

var GlobalStore = New() // глобальное хранилище task

func New() *TaskStore {
	return &TaskStore{
		tasks: make(map[int]Task),
	}
}

// получение задачи
func (ts *TaskStore) GetTask(id int) (Task, bool) {
	task, err := ts.tasks[id]
	return task, err
}

// создание задачи
func (ts *TaskStore) CreateTask(task Task) (Task, error) {
	var count int = 0
	for _, t := range ts.tasks {
		if t.Status == TaskStatusProcessing || t.Status == TaskStatusCreated {
			count++
		}
	}
	if count == 3 {
		return Task{}, response.ErrorTaskLimit
	}
	task.ID = len(ts.tasks) + 1
	ts.tasks[task.ID] = task
	return task, nil
}

// обновление задачи + ссылки
func (ts *TaskStore) UpdateTask(id int, links []string) (Task, error) {
	task, exists := ts.tasks[id]
	if !exists {
		return Task{}, response.TaskNotFound
	}
	task.Status = TaskStatusProcessing
	task.Links = append(task.Links, links...)
	ts.tasks[id] = task
	return task, nil
}

// обновление статуса задачи
func (ts *TaskStore) UpdateTaskStatus(id int, status, url string) (Task, error) {
	task, exists := ts.tasks[id]
	if !exists {
		return Task{}, response.TaskNotFound
	}
	task.Status = status
	if url != "" {
		task.Download_URL = url
	}
	ts.tasks[id] = task
	return task, nil
}

// установка URL для скачивания
func (ts *TaskStore) UpdateTaskURL(id int, url string) (Task, error) {
	return ts.UpdateTaskStatus(id, TaskStatusCompleted, url)
}
