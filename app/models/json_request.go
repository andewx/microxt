package models

import (
	"encoding/json"
	"fmt"
)

type JsonRequest struct {
	Controller string            `json:"controller"`
	RouteKey   string            `json:"routekey"`
	SessionKey string            `json:"sessionkey"`
	Paramaters map[string]string `json:"params"`
}

func NewJsonRequest(message string) *JsonRequest {
	r := &JsonRequest{}

	err := json.Unmarshal([]byte(message), r)
	if err != nil {
		fmt.Printf("Failed to decode json request: %s\n", err)
		return nil
	}
	return r
}
