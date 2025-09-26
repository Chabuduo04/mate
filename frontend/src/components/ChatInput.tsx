import React, { useState, useRef } from 'react';
import { cn } from '../utils/cn';
import { Send, Mic, MicOff, Loader2 } from 'lucide-react';

interface ChatInputProps {
  onSendMessage: (message: string) => void;
  onStartRecording: () => void;
  onStopRecording: () => void;
  isRecording: boolean;
  isLoading: boolean;
  disabled?: boolean;
}

export const ChatInput: React.FC<ChatInputProps> = ({
  onSendMessage,
  onStartRecording,
  onStopRecording,
  isRecording,
  isLoading,
  disabled = false,
}) => {
  const [message, setMessage] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (message.trim() && !isLoading && !disabled) {
      onSendMessage(message);
      setMessage('');
      if (textareaRef.current) {
        textareaRef.current.style.height = 'auto';
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleTextareaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
    const textarea = e.target;
    textarea.style.height = 'auto';
    textarea.style.height = `${Math.min(textarea.scrollHeight, 120)}px`;
  };

  return (
    <div className="border-t border-gray-200 bg-white p-4">
      <form onSubmit={handleSubmit} className="flex items-end gap-2">
        <div className="flex-1 relative">
          <textarea
            ref={textareaRef}
            value={message}
            onChange={handleTextareaChange}
            onKeyDown={handleKeyDown}
            placeholder="输入消息..."
            disabled={disabled || isLoading}
            className={cn(
              "input-field resize-none min-h-[44px] max-h-[120px] pr-12",
              disabled && "opacity-50 cursor-not-allowed"
            )}
            rows={1}
          />
          <button
            type="button"
            onClick={isRecording ? onStopRecording : onStartRecording}
            disabled={disabled || isLoading}
            className={cn(
              "absolute right-2 bottom-2 p-2 rounded-full transition-colors",
              isRecording
                ? "bg-red-500 hover:bg-red-600 text-white"
                : "bg-gray-100 hover:bg-gray-200 text-gray-600",
              disabled && "opacity-50 cursor-not-allowed"
            )}
            title={isRecording ? "停止录音" : "开始录音"}
          >
            {isRecording ? <MicOff size={16} /> : <Mic size={16} />}
          </button>
        </div>
        <button
          type="submit"
          disabled={!message.trim() || isLoading || disabled}
          className={cn(
            "btn-primary flex items-center gap-2 px-4 py-2 min-w-[80px]",
            (!message.trim() || isLoading || disabled) && "opacity-50 cursor-not-allowed"
          )}
        >
          {isLoading ? (
            <Loader2 size={16} className="animate-spin" />
          ) : (
            <Send size={16} />
          )}
          {isLoading ? "发送中" : "发送"}
        </button>
      </form>
    </div>
  );
};
