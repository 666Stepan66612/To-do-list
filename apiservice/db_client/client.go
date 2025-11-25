package db_client

import (
	models "apiservice/models"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type DBClient struct {
	BaseURL string
	Client  *http.Client
}

func NewDBClient(baseURL string) *DBClient {
	return &DBClient{
		BaseURL: "http://db-service:8080",
		Client:  &http.Client{},
	}
}

func (c *DBClient) CreateTask(task *models.CreateTaskRequest) (*models.Task, error) {
    jsonData, err := json.Marshal(task)
    if err != nil {
        return nil, err
    }

    resp, err := c.Client.Post(c.BaseURL+"/create", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var createdTask models.Task
    if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
        return nil, err
    }

    return &createdTask, nil
}

func (c *DBClient) GetAllTasks() ([]models.Task, error){
	resp, err := c.Client.Get(c.BaseURL + "/get")

	if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var tasks []models.Task
    if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (c *DBClient) GetCompleted() ([]models.Task, error){
	resp, err := c.Client.Get(c.BaseURL + "/get?complete=true")

	if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var tasks []models.Task
    if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (c *DBClient) GetUncompleted() ([]models.Task, error){
	resp, err := c.Client.Get(c.BaseURL + "/get?complete=false")

	if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var tasks []models.Task
    if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (c *DBClient) DeleteTask(id int) error {
    idStr := strconv.Itoa(id)
    url := c.BaseURL + "/delete/" + idStr
    
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return err
    }

    resp, err := c.Client.Do(req)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    return nil
}

func (c *DBClient) CompleteTask(id int) error {
    idStr := strconv.Itoa(id)
    url := c.BaseURL + "/complete/" + idStr
    
    req, err := http.NewRequest("PUT", url, nil)
    if err != nil {
        return err
    }

    resp, err := c.Client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

func (c *DBClient) GetTaskByID(id int) (*models.Task, error) {
    idStr := strconv.Itoa(id)
    url := c.BaseURL + "/getbyid/" + idStr
    
    resp, err := c.Client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var task models.Task
    if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
        return nil, err
    }

    return &task, nil
}

func (c *DBClient) GetTaskByName(name string) (*models.Task, error) {
    url := c.BaseURL + "/getbyname/" + name
    
    resp, err := c.Client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var task models.Task
    if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
        return nil, err
    }

    return &task, nil
}