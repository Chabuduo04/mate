import axios from 'axios';
import { Role, ChatRequest, ChatResponse, VoiceChatResponse, ASRResponse, TTSRequest, TTSResponse } from '../types';

const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
});

// Attach JWT if present
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwt_token');
  if (token && config && config.headers) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
});

export const roleService = {
  async getRoles(): Promise<Role[]> {
    const response = await api.get<Role[]>('/roles');
    return response.data;
  },
};

export const chatService = {
  async sendMessage(request: ChatRequest): Promise<ChatResponse> {
    const response = await api.post<ChatResponse>('/llm', request);
    return response.data;
  },
};

export const asrService = {
  async transcribe(audioFile: File): Promise<ASRResponse> {
    const formData = new FormData();
    formData.append('audio', audioFile);
    
    const response = await api.post<ASRResponse>('/asr', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

export const ttsService = {
  async synthesize(request: TTSRequest): Promise<TTSResponse> {
    const response = await api.post<TTSResponse>('/tts', request);
    return response.data;
  },
  async listVoices(): Promise<any> {
    const response = await api.get('/voice/list');
    return response.data;
  },
};

export const uploadService = {
  async uploadAudio(audioFile: File): Promise<{ upload: boolean }> {
    const formData = new FormData();
    formData.append('audio', audioFile);
    
    const response = await api.post<{ upload: boolean }>('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

export const authService = {
  async register(username: string, password: string) {
    const resp = await api.post('/register', { username, password });
    return resp.data;
  },
  async login(username: string, password: string) {
    const resp = await api.post('/login', { username, password });
    return resp.data; // should contain token, user_id, username
  },
  logout() {
    localStorage.removeItem('jwt_token');
    localStorage.removeItem('user_id');
    localStorage.removeItem('username');
  }
};

export const voiceChatService = {
  async sendVoiceMessage(audioFile: File, roleId: string, userId?: string, voice?: string): Promise<VoiceChatResponse> {
    const formData = new FormData();
    formData.append('audio', audioFile);
    formData.append('role_id', roleId);
    if (userId) {
      formData.append('user_id', userId);
    }
    if (voice) {
      formData.append('voice', voice);
    }
    
    const response = await api.post<VoiceChatResponse>('/voice-chat', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};
