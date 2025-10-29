package respond

type ChatTextRespond struct {
	ReplyText   string `json:"reply_text"`
	AudioURL    string `json:"audio_url,omitempty"`
	AudioBase64 string `json:"audio_base64,omitempty"` // if returning base64 payload
}
