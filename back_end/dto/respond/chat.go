package respond

type ChatTextRespond struct {
	ReplyText   string `json:"reply_text"`
	AudioURL    string `json:"audio_url,omitempty"`
	AudioBase64 string `json:"audio_base64,omitempty"` // if returning base64 payload
}

type ChatRecordRespond struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id"`
	Message   string `json:"message"`
	Response  string `json:"response"`
	CreatedAt int64  `json:"created_at"`
}
