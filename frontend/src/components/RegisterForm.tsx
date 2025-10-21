import React, { useState } from 'react';
import { authService } from '../services/api';

export const RegisterForm: React.FC<{onRegistered: () => void}> = ({ onRegistered }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res: any = await authService.register(username, password);
      if (res && res.username) {
        onRegistered();
      } else {
        setError('注册失败');
      }
    } catch (err: any) {
      setError(err?.response?.data?.error || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={submit} className="space-y-3">
      {error && <div className="text-red-500">{error}</div>}
      <input className="input" placeholder="用户名" value={username} onChange={(e) => setUsername(e.target.value)} />
      <input className="input" type="password" placeholder="密码（至少6位）" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button className="btn-primary" type="submit" disabled={loading}>{loading ? '注册中...' : '注册'}</button>
    </form>
  );
};
