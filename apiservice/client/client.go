package client

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
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *DBClient) CreateTask(task *models.CreateTaskRequest, userID int) (*models.Task, error) {
	jsonData, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + "/create?user_id=" + strconv.Itoa(userID)
	resp, err := c.Client.Post(url, "application/json", bytes.NewBuffer(jsonData))
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

func (c *DBClient) GetAllTasks(userID int) ([]models.Task, error) {
	resp, err := c.Client.Get(c.BaseURL + "/get?user_id=" + strconv.Itoa(userID))

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

func (c *DBClient) GetCompleted(userID int) ([]models.Task, error) {
	resp, err := c.Client.Get(c.BaseURL + "/get?complete=true&user_id=" + strconv.Itoa(userID))

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

func (c *DBClient) GetUncompleted(userID int) ([]models.Task, error) {
	resp, err := c.Client.Get(c.BaseURL + "/get?complete=false&user_id=" + strconv.Itoa(userID))

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

func (c *DBClient) DeleteTask(id, userID int) error {
	idStr := strconv.Itoa(id)
	url := c.BaseURL + "/delete/" + idStr + "?user_id=" + strconv.Itoa(userID)

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

func (c *DBClient) CompleteTask(id, userID int) error {
	idStr := strconv.Itoa(id)
	url := c.BaseURL + "/complete/" + idStr + "?user_id=" + strconv.Itoa(userID)

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
