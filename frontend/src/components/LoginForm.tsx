import React, { useState } from 'react';
import { authService } from '../services/api';

export const LoginForm: React.FC<{onLogin: () => void}> = ({ onLogin }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      const res: any = await authService.login(username, password);
      if (res && res.token) {
        localStorage.setItem('jwt_token', res.token);
        localStorage.setItem('user_id', res.user_id);
        localStorage.setItem('username', res.username);
        onLogin();
      } else {
        setError('登录失败');
      }
    } catch (err: any) {
      setError(err?.response?.data?.error || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={submit} className="space-y-3">
      {error && <div className="text-red-500">{error}</div>}
      <input className="input" placeholder="用户名" value={username} onChange={(e) => setUsername(e.target.value)} />
      <input className="input" type="password" placeholder="密码" value={password} onChange={(e) => setPassword(e.target.value)} />
      <button className="btn-primary" type="submit" disabled={loading}>{loading ? '登录中...' : '登录'}</button>
    </form>
  );
};
