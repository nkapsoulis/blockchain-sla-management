package models

type Response struct {
	Message string `json:"message,omitempty"`
	Err     string `json:"error,omitempty"`
	Name    string `json:"name,omitempty"`
}
