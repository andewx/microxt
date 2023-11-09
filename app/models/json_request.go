package models

import (
	"encoding/json"
)

type JsonRequest struct {
	RouteKey   string            `json:"routekey"`
	SessionKey int64             `json:"sessionkey"`
	Paramaters map[string]string `json:"params"`
}

func NewJsonRequest(message string) *JsonRequest {
	r := &JsonRequest{}

	err := json.Unmarshal([]byte(message), r)
	if err != nil {
		return nil
	}
	return r
}
