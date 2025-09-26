package models

type Role struct {
    ID        string   `json:"id"`
    Name      string   `json:"name"`
    Skills    []string `json:"skills"`
    Prompt    string   `json:"prompt"`
    VoiceType string   `json:"voice_type,omitempty"`
}
