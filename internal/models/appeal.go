package models

type Appeal struct {
	UserId    string  `json:"userId"`
	Name      string  `json:"name"`
	Picture   string  `json:"picture"`
	Username  string  `json:"username"`
	Message   string  `json:"message"`
}

type AppealRequest struct {
	UserId   string  `json:"userId"`
	Message  string  `json:"message"`
}