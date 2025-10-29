package respond

type LoginRespond struct {
	Token    string `json:"token"`
	Id       string `json:"id"`
	Username string `json:"username"`
}

type RegisterRespond struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
