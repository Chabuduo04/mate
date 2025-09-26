export interface Role {
  id: string;
  name: string;
  skills: string[];
  prompt: string;
}

export interface ChatRequest {
  role_id: string;
  message: string;
  user_id?: string;
  voice?: string;
}

export interface ChatResponse {
  reply_text: string;
  audio_url?: string;
  audio_base64?: string;
}

export interface VoiceChatResponse {
  transcribed_text: string;
  reply_text: string;
  audio_base64?: string;
}

export interface ASRResponse {
  text: string;
}

export interface TTSRequest {
  text: string;
  voice?: string;
}

export interface TTSResponse {
  audio_base64: string;
}

export interface Message {
  id: string;
  content: string;
  isUser: boolean;
  timestamp: Date;
  audioUrl?: string;
}

export interface ChatState {
  currentRole: Role | null;
  messages: Message[];
  isLoading: boolean;
  isRecording: boolean;
  isPlaying: boolean;
  selectedVoice?: string;
}
