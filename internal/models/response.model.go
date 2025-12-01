package models

type ResponseFailed struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResponseSucces struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Results any    `json:"results"`
}
