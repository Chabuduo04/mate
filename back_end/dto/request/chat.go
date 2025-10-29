package request

type ChatRequest struct {
	RoleID  string `json:"role_id" binding:"required"`
	Message string `json:"message" binding:"required"`
	UserID  string `json:"user_id"`         // optional for session key
	Voice   string `json:"voice,omitempty"` // optional selected TTS voice
}
