package dto

type PulseDto struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	ToId     string `json:"toId"`
}
