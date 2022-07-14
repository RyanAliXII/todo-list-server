package model

import "time"

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted int       `json:"isCompleted"`
	UserId      int       `json:"userId"`
	Timestamp   time.Time `json:"timestamp"`
}
type CreateTaskBody struct {
	Title       string `json:"title" validate:"required,max=20,min=3"`
	Description string `json:"description" validate:"required"`
	UserId      int    `json:"userId"`
}
type UpdateTaskBody struct {
	Title       string `json:"title" validate:"required,max=20,min=3"`
	Description string `json:"description" validate:"required"`
}
