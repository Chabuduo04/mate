import React from 'react';
import { Message } from '../types';
import { cn } from '../utils/cn';
import { Play, Volume2 } from 'lucide-react';

interface MessageBubbleProps {
  message: Message;
  onPlayAudio?: (audioBase64: string) => void;
}

export const MessageBubble: React.FC<MessageBubbleProps> = ({
  message,
  onPlayAudio,
}) => {
  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div
      className={cn(
        "flex mb-4 animate-fade-in",
        message.isUser ? "justify-end" : "justify-start"
      )}
    >
      <div
        className={cn(
          "message-bubble",
          message.isUser ? "message-user" : "message-ai"
        )}
      >
        <p className="text-sm leading-relaxed">{message.content}</p>
        <div className="flex items-center justify-between mt-2">
          <span className="text-xs opacity-70">
            {formatTime(message.timestamp)}
          </span>
          {message.audioUrl && onPlayAudio && (
            <button
              onClick={() => onPlayAudio(message.audioUrl!)}
              className="ml-2 p-1 hover:bg-white/20 rounded-full transition-colors"
              title="播放语音"
            >
              <Volume2 size={16} />
            </button>
          )}
        </div>
      </div>
    </div>
  );
};
