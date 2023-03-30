package models

type APIResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  string            `json:"errors"`
	Data    map[string]string `json:"data"`
}
