package dto

import (
	"encoding/json"
	"errors"
	"time"
)

type TaskDTO struct {
	Name string
	Text string
}

type ErrorDTO struct {
    Message string    `json:"message"`
    Time    time.Time `json:"time"`
}

func (t TaskDTO) ValidateForCreate() error {
	if t.Name == "" {
		return errors.New("title is empty")
	}

	if t.Text == ""{
		return errors.New("book is empty")
	}

	return nil
}

func (e ErrorDTO)ToString() string{
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil{
		panic(err)
	}

	return string(b)
}