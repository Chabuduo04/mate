package api

type ChatRequest struct {
	RoleID  string `json:"role_id" binding:"required"`
	Message string `json:"message" binding:"required"`
	UserID  string `json:"user_id"`         // optional for session key
	Voice   string `json:"voice,omitempty"` // optional selected TTS voice
}

type ChatResponse struct {
	ReplyText   string `json:"reply_text"`
	AudioURL    string `json:"audio_url,omitempty"`
	AudioBase64 string `json:"audio_base64,omitempty"` // if returning base64 payload
}
