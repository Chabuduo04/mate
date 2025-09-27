import { useState, useRef, useCallback } from 'react';

export const useAudio = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const chunksRef = useRef<Blob[]>([]);

  const startRecording = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const mediaRecorder = new MediaRecorder(stream);
      mediaRecorderRef.current = mediaRecorder;
      chunksRef.current = [];

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          chunksRef.current.push(event.data);
        }
      };

      mediaRecorder.onstop = () => {
        const audioBlob = new Blob(chunksRef.current, { type: 'audio/wav' });
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

  const stopRecording = useCallback(() => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop();
      setIsRecording(false);
    }
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
      const audioBlob = new Blob(chunksRef.current, { type: 'audio/wav' });
      return new File([audioBlob], 'recording.wav', { type: 'audio/wav' });
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
