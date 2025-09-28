import { useState, useRef, useCallback } from 'react';

export const useAudio = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  // 记录实际使用的 mimeType
  const mimeTypeRef = useRef<string>('');

  const startRecording = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      // 自动适配 mimeType
      let mimeType = '';
      if (window.MediaRecorder && MediaRecorder.isTypeSupported) {
        if (MediaRecorder.isTypeSupported('audio/webm')) {
          mimeType = 'audio/webm';
        } else if (MediaRecorder.isTypeSupported('audio/ogg')) {
          mimeType = 'audio/ogg';
        } else if (MediaRecorder.isTypeSupported('audio/wav')) {
          mimeType = 'audio/wav';
        }
      }
      mimeTypeRef.current = mimeType;
      const mediaRecorder = mimeType
        ? new MediaRecorder(stream, { mimeType })
        : new MediaRecorder(stream);
      mediaRecorderRef.current = mediaRecorder;
      chunksRef.current = [];

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          chunksRef.current.push(event.data);
        }
      };

      mediaRecorder.onstop = () => {
        const audioBlob = new Blob(chunksRef.current, { type: mimeTypeRef.current || 'audio/webm' });
        const audioUrl = URL.createObjectURL(audioBlob);
        if (audioRef.current) {
          audioRef.current.src = audioUrl;
        }
        stream.getTracks().forEach(track => track.stop());
      };

      mediaRecorder.start();
      setIsRecording(true);
    } catch (error) {
      console.error('Error starting recording:', error);
    }
  }, []);

  // 停止录音，返回 Promise，onstop 完成后 resolve
  const stopRecording = useCallback((): Promise<void> => {
    return new Promise((resolve) => {
      if (mediaRecorderRef.current && isRecording) {
        const recorder = mediaRecorderRef.current;
        const handleStop = () => {
          recorder.removeEventListener('stop', handleStop);
          setIsRecording(false);
          resolve();
        };
        recorder.addEventListener('stop', handleStop);
        recorder.stop();
      } else {
        resolve();
      }
    });
  }, [isRecording]);

  // 当前播放的音频对象
  const currentAudioRef = useRef<HTMLAudioElement | null>(null);

  // 播放音频，确保同一时间只播放一段
  const playAudio = useCallback((audioBase64?: string) => {
    if (!audioBase64) return;
    let src = audioBase64;
    if (!/^data:audio\//.test(audioBase64)) {
      src = `data:audio/mp3;base64,${audioBase64}`;
    }
    // 停止之前的音频
    if (currentAudioRef.current) {
      currentAudioRef.current.pause();
      currentAudioRef.current.currentTime = 0;
      currentAudioRef.current = null;
    }
    const audio = new Audio(src);
    currentAudioRef.current = audio;
    audio.preload = 'auto';
    audio.onplay = () => setIsPlaying(true);
    audio.onended = () => {
      setIsPlaying(false);
      currentAudioRef.current = null;
    };
    audio.onerror = () => {
      setIsPlaying(false);
      currentAudioRef.current = null;
      // Fallback attempt with wav header if needed
      if (!/^data:audio\/wav/.test(src)) {
        const fallback = new Audio(`data:audio/wav;base64,${audioBase64}`);
        currentAudioRef.current = fallback;
        fallback.onplay = () => setIsPlaying(true);
        fallback.onended = () => {
          setIsPlaying(false);
          currentAudioRef.current = null;
        };
        fallback.onerror = () => {
          setIsPlaying(false);
          currentAudioRef.current = null;
        };
        fallback.play().catch(() => {});
      }
    };
    const playPromise = audio.play();
    if (playPromise && typeof playPromise.catch === 'function') {
      playPromise.catch(() => {
        setIsPlaying(false);
        currentAudioRef.current = null;
      });
    }
  }, []);

  const getRecordedAudio = useCallback((): File | null => {
    if (chunksRef.current.length > 0) {
      const mimeType = mimeTypeRef.current || 'audio/webm';
      const ext = mimeType.includes('ogg') ? 'ogg' : mimeType.includes('wav') ? 'wav' : 'webm';
      const audioBlob = new Blob(chunksRef.current, { type: mimeType });
      return new File([audioBlob], `recording.${ext}`, { type: mimeType });
    }
    return null;
  }, []);

  return {
    isRecording,
    isPlaying,
    startRecording,
    stopRecording,
    playAudio,
    getRecordedAudio,
  };
};
