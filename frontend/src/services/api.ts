import axios from 'axios';
import { Role, ChatRequest, ChatResponse, VoiceChatResponse, ASRResponse, TTSRequest, TTSResponse } from '../types';

// 不再全局使用 '/api' 作为 baseURL，因为后端路由有的在根路径，有的在 /api 下
const API_BASE_URL = '';

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
    // 后端在 /role/list 返回 { roles: [...] }
    const response = await api.get<{ roles: Role[] }>('/role/list');
    return response.data.roles;
  },
};

export const chatService = {
  async sendMessage(request: ChatRequest): Promise<ChatResponse> {
    // Chat 路由在后端被注册为 /api/llm（需要 JWT）
    const response = await api.post<ChatResponse>('/api/llm', request);
    return response.data;
  },
};

export const asrService = {
  async transcribe(audioFile: File): Promise<ASRResponse> {
    const formData = new FormData();
    formData.append('audio', audioFile);

    // ASR 接口在后端文档/handlers 中是 /api/asr
    const response = await api.post<ASRResponse>('/api/asr', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

export const ttsService = {
  async synthesize(request: TTSRequest): Promise<TTSResponse> {
    // TTS 也通常位于 /api/tts（受保护）
    const response = await api.post<TTSResponse>('/api/tts', request);
    return response.data;
  },
  async listVoices(): Promise<any> {
    // 列表接口在根路径下 /voice/list
    const response = await api.get('/voice/list');
    return response.data;
  },
};

export const uploadService = {
  async uploadAudio(audioFile: File): Promise<{ upload: boolean }> {
    const formData = new FormData();
    formData.append('audio', audioFile);

    // 上传走 /api/upload
    const response = await api.post<{ upload: boolean }>('/api/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

export const authService = {
  async register(username: string, password: string) {
    const resp = await api.post('/user/register', { username, password });
    return resp.data;
  },
  async login(username: string, password: string) {
    const resp = await api.post('/user/login', { username, password });
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

    // voice-chat 在后端被注册为 /api/voice-chat
    const response = await api.post<VoiceChatResponse>('/api/voice-chat', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};
