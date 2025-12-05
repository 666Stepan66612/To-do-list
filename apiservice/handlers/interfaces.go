package handlers

import "apiservice/models"

// DBClientInterface определяет методы клиента БД
type DBClientInterface interface {
	CreateTask(req *models.CreateTaskRequest, userID int) (*models.Task, error)
	GetAllTasks(userID int) ([]models.Task, error)
	DeleteTask(taskID, userID int) error
	CompleteTask(taskID, userID int) error
	GetCompleted(userID int) ([]models.Task, error)
	GetUncompleted(userID int) ([]models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetTaskByName(name string) (*models.Task, error)
	CreateCollection(req *models.CreateCollectionRequest, userID int) (*models.Collection, error)
	GetCollections(userID int) ([]models.Collection, error)
	DeleteCollection(collectionID, userID int) error
	GetTasksByCollection(collectionID, userID int) ([]models.Task, error)
}

// EventProducerInterface определяет методы продюсера Kafka
type EventProducerInterface interface {
	SendEvent(userID int, username, action, details, status string) error
}
