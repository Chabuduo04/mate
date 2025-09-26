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

  const playAudio = useCallback((audioBase64?: string) => {
    if (!audioBase64) return;
    // Use the standard mp3 MIME type for widest compatibility
    const src = `data:audio/mpeg;base64,${audioBase64}`;
    const audio = new Audio(src);
    audio.preload = 'auto';
    audio.onplay = () => setIsPlaying(true);
    audio.onended = () => setIsPlaying(false);
    audio.onerror = (e) => {
      setIsPlaying(false);
      // Fallback attempt with wav header if needed
      const fallback = new Audio(`data:audio/wav;base64,${audioBase64}`);
      fallback.onplay = () => setIsPlaying(true);
      fallback.onended = () => setIsPlaying(false);
      fallback.onerror = () => setIsPlaying(false);
      fallback.play().catch(() => {});
    };
    const playPromise = audio.play();
    if (playPromise && typeof playPromise.catch === 'function') {
      playPromise.catch(() => {
        setIsPlaying(false);
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
