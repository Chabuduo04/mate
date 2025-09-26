import React, { useEffect, useRef } from 'react';
import { Message, Role } from '../types';
import { MessageBubble } from './MessageBubble';
import { ChatInput } from './ChatInput';
import { cn } from '../utils/cn';
import { Bot, User, RotateCcw } from 'lucide-react';

interface ChatInterfaceProps {
  currentRole: Role | null;
  messages: Message[];
  isLoading: boolean;
  isRecording: boolean;
  onSendMessage: (message: string) => void;
  onStartRecording: () => void;
  onStopRecording: () => void;
  onPlayAudio: (audioBase64: string) => void;
  onClearMessages: () => void;
}

export const ChatInterface: React.FC<ChatInterfaceProps> = ({
  currentRole,
  messages,
  isLoading,
  isRecording,
  onSendMessage,
  onStartRecording,
  onStopRecording,
  onPlayAudio,
  onClearMessages,
}) => {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  if (!currentRole) {
    return (
      <div className="flex-1 flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Bot size={64} className="mx-auto text-gray-400 mb-4" />
          <h3 className="text-xl font-semibold text-gray-600 mb-2">
            请先选择一个AI角色
          </h3>
          <p className="text-gray-500">
            选择上方的角色开始对话
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-br from-primary-400 to-primary-600 rounded-full flex items-center justify-center text-white font-bold">
              {currentRole.name.charAt(0)}
            </div>
            <div>
              <h2 className="text-lg font-semibold text-gray-800">
                {currentRole.name}
              </h2>
              <p className="text-sm text-gray-500">
                {currentRole.skills.join(' • ')}
              </p>
            </div>
          </div>
          <button
            onClick={onClearMessages}
            className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
            title="清空对话"
          >
            <RotateCcw size={20} />
          </button>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
        {messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <Bot size={48} className="mx-auto text-gray-400 mb-3" />
              <h3 className="text-lg font-medium text-gray-600 mb-2">
                开始与 {currentRole.name} 对话
              </h3>
              <p className="text-gray-500 text-sm">
                输入消息或点击麦克风开始语音对话
              </p>
            </div>
          </div>
        ) : (
          messages.map((message) => (
            <MessageBubble
              key={message.id}
              message={message}
              onPlayAudio={onPlayAudio}
            />
          ))
        )}
        
        {isLoading && (
          <div className="flex justify-start">
            <div className="message-bubble message-ai">
              <div className="flex items-center gap-2">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce"></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                </div>
                <span className="text-sm text-gray-500">正在思考...</span>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <ChatInput
        onSendMessage={onSendMessage}
        onStartRecording={onStartRecording}
        onStopRecording={onStopRecording}
        isRecording={isRecording}
        isLoading={isLoading}
      />
    </div>
  );
};
