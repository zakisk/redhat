package models

type Response struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}
