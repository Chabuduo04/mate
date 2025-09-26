import React, { useState, useEffect } from 'react';
import { Role } from './types';
import { roleService, ttsService } from './services/api';
import { useChat } from './hooks/useChat';
import { useAudio } from './hooks/useAudio';
import { RoleSelector } from './components/RoleSelector';
import { ChatInterface } from './components/ChatInterface';
import { LoadingSpinner } from './components/LoadingSpinner';
import { cn } from './utils/cn';

function App() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoadingRoles, setIsLoadingRoles] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const {
    currentRole,
    messages,
    isLoading,
    selectRole,
    sendMessage,
    sendVoiceMessage,
    transcribeAudio,
    synthesizeSpeech,
    clearMessages,
    setSelectedVoice,
    selectedVoice,
  } = useChat();

  const {
    isRecording,
    isPlaying,
    startRecording,
    stopRecording,
    playAudio,
    getRecordedAudio,
  } = useAudio();

  // 加载角色列表
  useEffect(() => {
    const loadRoles = async () => {
      try {
        setIsLoadingRoles(true);
        setError(null);
        const rolesData = await roleService.getRoles();
        setRoles(rolesData);
      } catch (err) {
        console.error('Error loading roles:', err);
        setError('加载角色失败，请刷新页面重试');
      } finally {
        setIsLoadingRoles(false);
      }
    };

    loadRoles();
  }, []);

  // 加载音色列表
  const [voices, setVoices] = useState<any[]>([]);
  useEffect(() => {
    const loadVoices = async () => {
      try {
        const list = await ttsService.listVoices();
        // 兼容可能返回数组对象或简单数组
        if (Array.isArray(list)) {
          setVoices(list);
        } else if (list && Array.isArray(list.data)) {
          setVoices(list.data);
        }
      } catch (e) {
        console.warn('加载音色失败:', e);
      }
    };
    loadVoices();
  }, []);

  // 处理语音录制
  const handleStartRecording = async () => {
    try {
      await startRecording();
    } catch (error) {
      console.error('Error starting recording:', error);
      alert('无法访问麦克风，请检查权限设置');
    }
  };

  const handleStopRecording = async () => {
    try {
      await stopRecording();
      
      // 获取录制的音频并直接发送语音消息
      const audioFile = getRecordedAudio();
      if (audioFile) {
        await sendVoiceMessage(audioFile);
      }
    } catch (error) {
      console.error('Error processing recording:', error);
      alert('语音处理失败，请重试');
    }
  };

  // 处理音频播放
  const handlePlayAudio = async (audioBase64: string) => {
    try {
      playAudio(audioBase64);
    } catch (error) {
      console.error('Error playing audio:', error);
    }
  };

  // 自动播放最新的AI语音回复
  useEffect(() => {
    if (messages.length === 0) return;
    const last = messages[messages.length - 1];
    if (!last.isUser && last.audioUrl) {
      try {
        playAudio(last.audioUrl);
      } catch (e) {
        console.warn('Auto play failed:', e);
      }
    }
  }, [messages, playAudio]);

  // 处理消息发送
  const handleSendMessage = async (message: string) => {
    try {
      await sendMessage(message);
    } catch (error) {
      console.error('Error sending message:', error);
      alert('发送消息失败，请重试');
    }
  };

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-500 text-6xl mb-4">⚠️</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">加载失败</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="btn-primary"
          >
            刷新页面
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">
            Mate
          </h1>
          <p className="text-gray-600">
            AI角色扮演语音聊天平台
          </p>
        </div>

        <div className="max-w-6xl mx-auto">
          {!currentRole ? (
            <div className="card p-8">
              {isLoadingRoles ? (
                <div className="text-center">
                  <LoadingSpinner size="lg" className="mx-auto mb-4" />
                  <p className="text-gray-600">加载角色中...</p>
                </div>
              ) : (
                <RoleSelector
                  roles={roles}
                  selectedRole={currentRole}
                  onSelectRole={selectRole}
                  isLoading={isLoadingRoles}
                />
              )}
            </div>
          ) : (
            <div className="card h-[600px] flex flex-col">
              {/* Voice selector */}
              <div className="p-4 border-b border-gray-200 flex items-center gap-3">
                <label className="text-sm text-gray-600">音色</label>
                <select
                  className="border rounded px-2 py-1 text-sm"
                  value={selectedVoice || ''}
                  onChange={(e) => setSelectedVoice(e.target.value || undefined)}
                >
                  <option value="">默认</option>
                  {voices.map((v: any, idx: number) => {
                    // 只传 voice_type，label 用 voice_name
                    const value = v.voice_type;
                    const label = v.voice_name;
                    return (
                      <option key={idx} value={value}>
                        {label}
                      </option>
                    );
                  })}
                </select>
              </div>
              <ChatInterface
                currentRole={currentRole}
                messages={messages}
                isLoading={isLoading}
                isRecording={isRecording}
                onSendMessage={handleSendMessage}
                onStartRecording={handleStartRecording}
                onStopRecording={handleStopRecording}
                onPlayAudio={handlePlayAudio}
                onClearMessages={clearMessages}
              />
            </div>
          )}
        </div>

        {/* 状态指示器 */}
        {(isRecording || isPlaying) && (
          <div className="fixed bottom-4 right-4 z-50">
            <div className="bg-white rounded-lg shadow-lg p-4 flex items-center gap-3">
              <div className={cn(
                "w-3 h-3 rounded-full animate-pulse",
                isRecording ? "bg-red-500" : "bg-green-500"
              )} />
              <span className="text-sm font-medium text-gray-700">
                {isRecording ? "正在录音..." : "正在播放..."}
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
