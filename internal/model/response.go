package model

type ResponseMessage struct {
	Message string `json:"message"`
}

type TokenResponseMessage struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
