package domain

type Task struct {
	Id              string `json:"id"`
	ImageBase       string `json:"image"`
	Status          string `json:"status"`
	FilterName      string `json:"filterName"`
	FilterParametes any    `json:"filterParameters"`
	Result          string `json:"result"`
}

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Session struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
}
