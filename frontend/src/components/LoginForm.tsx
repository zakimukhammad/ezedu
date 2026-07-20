import { useState } from 'preact/hooks';
import { authApi } from '../lib/api';

export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');

    if (!email.trim()) {
      setError('Email wajib diisi');
      return;
    }
    if (!password.trim()) {
      setError('Kata sandi wajib diisi');
      return;
    }

    setLoading(true);
    const { data, error: apiError } = await authApi.login(email, password);
    setLoading(false);

    if (apiError) {
      setError(apiError);
      return;
    }

    // Redirect to child profile selection
    window.location.href = '/profil';
  };

  return (
    <form onSubmit={handleSubmit} class="auth-form">
      <div class="form-group">
        <label class="form-label" for="login-email">Email</label>
        <input
          id="login-email"
          class="form-input"
          type="email"
          placeholder="contoh@email.com"
          value={email}
          onInput={(e) => setEmail((e.target as HTMLInputElement).value)}
          autocomplete="email"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="login-password">Kata Sandi</label>
        <input
          id="login-password"
          class="form-input"
          type="password"
          placeholder="Masukkan kata sandi"
          value={password}
          onInput={(e) => setPassword((e.target as HTMLInputElement).value)}
          autocomplete="current-password"
        />
      </div>

      {error && <p class="form-error" role="alert">⚠️ {error}</p>}

      <button
        type="submit"
        class="btn btn-primary btn-lg w-full mt-md"
        disabled={loading}
        id="login-submit"
      >
        {loading ? 'Memproses...' : 'Masuk 🚀'}
      </button>

      <p class="auth-footer text-center mt-lg text-muted">
        Belum punya akun?{' '}
        <a href="/daftar" class="auth-link">Daftar gratis</a>
      </p>
    </form>
  );
}
