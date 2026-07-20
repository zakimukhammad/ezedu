import { useState } from 'preact/hooks';
import { authApi } from '../lib/api';

export default function SignupForm() {
  const [parentName, setParentName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');

    // Validation — all in Bahasa Indonesia
    if (!parentName.trim()) {
      setError('Nama orang tua wajib diisi');
      return;
    }
    if (!email.trim()) {
      setError('Email wajib diisi');
      return;
    }
    if (password.length < 6) {
      setError('Kata sandi minimal 6 karakter');
      return;
    }
    if (password !== confirmPassword) {
      setError('Konfirmasi kata sandi tidak cocok');
      return;
    }

    setLoading(true);
    const { data, error: apiError } = await authApi.signup(email, password, parentName);
    setLoading(false);

    if (apiError) {
      setError(apiError);
      return;
    }

    // Redirect to child profile creation
    window.location.href = '/profil';
  };

  return (
    <form onSubmit={handleSubmit} class="auth-form">
      <div class="form-group">
        <label class="form-label" for="signup-name">Nama Orang Tua</label>
        <input
          id="signup-name"
          class="form-input"
          type="text"
          placeholder="Masukkan nama Anda"
          value={parentName}
          onInput={(e) => setParentName((e.target as HTMLInputElement).value)}
          autocomplete="name"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="signup-email">Email</label>
        <input
          id="signup-email"
          class="form-input"
          type="email"
          placeholder="contoh@email.com"
          value={email}
          onInput={(e) => setEmail((e.target as HTMLInputElement).value)}
          autocomplete="email"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="signup-password">Kata Sandi</label>
        <input
          id="signup-password"
          class="form-input"
          type="password"
          placeholder="Minimal 6 karakter"
          value={password}
          onInput={(e) => setPassword((e.target as HTMLInputElement).value)}
          autocomplete="new-password"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="signup-confirm">Konfirmasi Kata Sandi</label>
        <input
          id="signup-confirm"
          class="form-input"
          type="password"
          placeholder="Ulangi kata sandi"
          value={confirmPassword}
          onInput={(e) => setConfirmPassword((e.target as HTMLInputElement).value)}
          autocomplete="new-password"
        />
      </div>

      {error && <p class="form-error" role="alert">⚠️ {error}</p>}

      <button
        type="submit"
        class="btn btn-primary btn-lg w-full mt-md"
        disabled={loading}
        id="signup-submit"
      >
        {loading ? 'Mendaftar...' : 'Buat Akun 🎉'}
      </button>

      <p class="auth-footer text-center mt-lg text-muted">
        Sudah punya akun?{' '}
        <a href="/masuk" class="auth-link">Masuk di sini</a>
      </p>
    </form>
  );
}
