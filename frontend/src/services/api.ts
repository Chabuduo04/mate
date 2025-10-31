import axios, { AxiosResponse } from 'axios';
import { Role, ChatRequest, ChatResponse, VoiceChatResponse, ASRResponse, TTSRequest, TTSResponse, Message } from '../types';

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
    // 后端在 /role/list 返回 { data: { roles: [...] } }
    const response = await api.get<ApiResponse<{ roles: Role[] }>>('/role/list');
    return unwrap(response).roles;
  },
};

export const chatService = {
  async sendMessage(request: ChatRequest): Promise<ChatResponse> {
    // Chat 路由在后端被注册为 /chat/llm（需要 JWT）
    const response = await api.post<ApiResponse<ChatResponse>>('/chat/llm', request);
    return unwrap(response);
  },
  async getChatList(params: { user_id: string; role_id: string }): Promise<Message[]> {
    // /chat/list 返回 { data: { records: Message[] } }
    const response = await api.post<ApiResponse<{ records: Message[] }>>('/chat/list', {
      user_id: params.user_id,
      role_id: params.role_id
    });
    return unwrap(response).records;
  },
};

export const asrService = {
  async transcribe(audioFile: File): Promise<ASRResponse> {
    const formData = new FormData();
    formData.append('audio', audioFile);

    // ASR 接口在后端文档/handlers 中是 /api/asr
    const response = await api.post<ApiResponse<ASRResponse>>('/api/asr', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return unwrap(response);
  },
};

export const ttsService = {
  async synthesize(request: TTSRequest): Promise<TTSResponse> {
    // TTS 也通常位于 /api/tts（受保护）
    const response = await api.post<ApiResponse<TTSResponse>>('/api/tts', request);
    return unwrap(response);
  },
  async listVoices(): Promise<any> {
    // 列表接口在根路径下 /voice/list
    const response = await api.get<ApiResponse<any>>('/voice/list');
    return unwrap(response);
  },
};

export const uploadService = {
  async uploadAudio(audioFile: File): Promise<{ upload: boolean }> {
    const formData = new FormData();
    formData.append('audio', audioFile);

    // 上传走 /api/upload
    const response = await api.post<ApiResponse<{ upload: boolean }>>('/api/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return unwrap(response);
  },
};

export const authService = {
  async register(username: string, password: string) {
    const resp = await api.post<ApiResponse<any>>('/user/register', { username, password });
    return unwrap(resp);
  },
  async login(username: string, password: string) {
    const resp = await api.post<ApiResponse<{
      token: string;
      user_id: string;
      username: string;
    }>>('/user/login', { username, password });
    const data = unwrap(resp);
    return {
      token: data.token,
      user_id: data.user_id,
      username: data.username
    };
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

    // voice-chat 在后端被注册为 /chat/voice-chat
    const response = await api.post<ApiResponse<VoiceChatResponse>>('/chat/voice-chat', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return unwrap(response);
  },
};

// Generic API response wrapper used by backend: { code?: number, message?: string, data: T }
type ApiResponse<T> = {
  code?: number;
  message?: string;
  data: T;
};

function unwrap<T>(resp: AxiosResponse<ApiResponse<T>>): T {
  // defensive: if backend returns directly, fall back to resp.data
  if (resp && resp.data && Object.prototype.hasOwnProperty.call(resp.data, 'data')) {
    return resp.data.data as T;
  }
  // @ts-ignore
  return resp.data as unknown as T;
}
