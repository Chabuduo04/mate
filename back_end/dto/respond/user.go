package respond

type LoginRespond struct {
	Token    string `json:"token"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type RegisterRespond struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}
