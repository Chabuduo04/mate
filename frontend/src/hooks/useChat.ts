import { useState, useCallback } from 'react';
import { Role, Message, ChatState } from '../types';
import { chatService, asrService, ttsService, voiceChatService } from '../services/api';

export const useChat = () => {
  const [state, setState] = useState<ChatState>({
    currentRole: null,
    messages: [],
    isLoading: false,
    isRecording: false,
    isPlaying: false,
    selectedVoice: undefined,
  });

  const selectRole = useCallback((role: Role) => {
    setState(prev => ({
      ...prev,
      currentRole: role,
      messages: [],
    }));
  }, []);

  const sendMessage = useCallback(async (content: string) => {
    if (!state.currentRole || !content.trim()) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      content: content.trim(),
      isUser: true,
      timestamp: new Date(),
    };

    setState(prev => ({
      ...prev,
      messages: [...prev.messages, userMessage],
      isLoading: true,
    }));

    try {
      const response = await chatService.sendMessage({
        role_id: state.currentRole.id,
        message: content.trim(),
        user_id: 'user123', // 可以改为动态用户ID
        voice: state.selectedVoice,
      });

      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: response.reply_text,
        isUser: false,
        timestamp: new Date(),
        audioUrl: response.audio_base64,
      };

      setState(prev => ({
        ...prev,
        messages: [...prev.messages, aiMessage],
        isLoading: false,
      }));
    } catch (error) {
      console.error('Error sending message:', error);
      setState(prev => ({
        ...prev,
        isLoading: false,
      }));
    }
  }, [state.currentRole, state.selectedVoice]);

  const transcribeAudio = useCallback(async (audioFile: File) => {
    try {
      const response = await asrService.transcribe(audioFile);
      return response.text;
    } catch (error) {
      console.error('Error transcribing audio:', error);
      throw error;
    }
  }, []);

  const synthesizeSpeech = useCallback(async (text: string) => {
    try {
      const response = await ttsService.synthesize({ text, voice: state.selectedVoice });
      return response.audio_base64;
    } catch (error) {
      console.error('Error synthesizing speech:', error);
      throw error;
    }
  }, [state.selectedVoice]);

  const sendVoiceMessage = useCallback(async (audioFile: File) => {
    if (!state.currentRole) return;

    setState(prev => ({
      ...prev,
      isLoading: true,
    }));

    try {
      const response = await voiceChatService.sendVoiceMessage(
        audioFile, 
        state.currentRole.id, 
        'user123', // 可以改为动态用户ID
        state.selectedVoice
      );

      // 添加用户消息（转录的文字）
      const userMessage: Message = {
        id: Date.now().toString(),
        content: response.transcribed_text,
        isUser: true,
        timestamp: new Date(),
      };

      // 添加AI回复消息
      const aiMessage: Message = {
        id: (Date.now() + 1).toString(),
        content: response.reply_text,
        isUser: false,
        timestamp: new Date(),
        audioUrl: response.audio_base64,
      };

      setState(prev => ({
        ...prev,
        messages: [...prev.messages, userMessage, aiMessage],
        isLoading: false,
      }));
    } catch (error) {
      console.error('Error sending voice message:', error);
      setState(prev => ({
        ...prev,
        isLoading: false,
      }));
    }
  }, [state.currentRole, state.selectedVoice]);

  const clearMessages = useCallback(() => {
    setState(prev => ({
      ...prev,
      messages: [],
    }));
  }, []);

  return {
    ...state,
    selectRole,
    sendMessage,
    sendVoiceMessage,
    transcribeAudio,
    synthesizeSpeech,
    clearMessages,
    setSelectedVoice: (voice?: string) => setState(prev => ({ ...prev, selectedVoice: voice })),
  };
};
