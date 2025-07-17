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

var GlobalStore = New()

func New() *TaskStore {
	return &TaskStore{
		tasks: make(map[int]Task),
	}
}

func (ts *TaskStore) GetTask(id int) (Task, bool) {
	task, err := ts.tasks[id]
	return task, err
}

func (ts *TaskStore) CreateTask(task Task) (Task, error) {
	if len(ts.tasks) == 3 {
		return Task{}, response.ErrorTaskLimit
	}
	task.ID = len(ts.tasks) + 1
	ts.tasks[task.ID] = task
	return task, nil
}

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

func (ts *TaskStore) UpdateTaskURL(id int, url string) (Task, error) {
	task, exists := ts.tasks[id]
	if !exists {
		return Task{}, response.TaskNotFound
	}
	task.Status = TaskStatusCompleted
	task.Download_URL = url
	ts.tasks[id] = task
	return task, nil
}
