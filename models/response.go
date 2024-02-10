package models

type Response struct {
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}

type ListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Count   int    `json:"count"`
	Files   []*File `json:"files"`
}

type File struct {
	Name       string
	Mode       string
	ModifiedAt string
	Size       string
}
