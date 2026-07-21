import { useState } from 'preact/hooks';
import { authApi } from '../lib/api';

interface Props {
  isOpen: boolean;
  mode: 'verify' | 'setup';
  title?: string;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function ParentPinModal({ isOpen, mode, title, onSuccess, onCancel }: Props) {
  const [pin, setPin] = useState(['', '', '', '']);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  if (!isOpen) return null;

  const handleDigitChange = (index: number, value: string) => {
    if (!/^\d*$/.test(value)) return;
    const newPin = [...pin];
    newPin[index] = value.slice(-1);
    setPin(newPin);
    setError('');

    // Auto-focus next input
    if (value && index < 3) {
      const nextInput = document.getElementById(`pin-input-${index + 1}`);
      nextInput?.focus();
    }
  };

  const handleKeyDown = (index: number, e: KeyboardEvent) => {
    if (e.key === 'Backspace' && !pin[index] && index > 0) {
      const prevInput = document.getElementById(`pin-input-${index - 1}`);
      prevInput?.focus();
    }
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    const pinStr = pin.join('');
    if (pinStr.length !== 4) {
      setError('PIN harus 4 digit angka');
      return;
    }

    setLoading(true);
    setError('');

    if (mode === 'setup') {
      const { error: apiError } = await authApi.updatePin(pinStr);
      setLoading(false);
      if (apiError) {
        setError(apiError);
        return;
      }
      onSuccess();
    } else {
      const { data, error: apiError } = await authApi.verifyPin(pinStr);
      setLoading(false);
      if (apiError || !data?.valid) {
        setError(apiError || 'PIN salah. Silakan coba lagi.');
        setPin(['', '', '', '']);
        document.getElementById('pin-input-0')?.focus();
        return;
      }
      onSuccess();
    }
  };

  const defaultTitle = mode === 'setup'
    ? 'Atur PIN Orang Tua 🔒'
    : 'Verifikasi PIN Orang Tua 🔒';

  return (
    <div class="modal-overlay" onClick={onCancel}>
      <div class="modal-content pin-modal-content animate-scale-in" onClick={(e) => e.stopPropagation()}>
        <h3>{title || defaultTitle}</h3>
        <p class="pin-modal-subtitle">
          {mode === 'setup'
            ? 'PIN ini digunakan untuk mengamankan profil dan pengaturan orang tua.'
            : 'Masukkan 4 digit PIN Anda untuk melanjutkan.'
          }
        </p>

        <form onSubmit={handleSubmit}>
          <div class="pin-inputs-row">
            {pin.map((digit, i) => (
              <input
                key={i}
                id={`pin-input-${i}`}
                type="password"
                maxLength={1}
                inputMode="numeric"
                class="pin-digit-box"
                value={digit}
                onInput={(e) => handleDigitChange(i, (e.target as HTMLInputElement).value)}
                onKeyDown={(e) => handleKeyDown(i, e)}
                autoFocus={i === 0}
              />
            ))}
          </div>

          {error && <p class="form-error mt-md text-center" role="alert">⚠️ {error}</p>}

          <div class="modal-actions mt-xl">
            <button type="button" class="btn btn-secondary" onClick={onCancel}>
              Batal
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              disabled={loading || pin.join('').length !== 4}
              id="pin-submit-btn"
            >
              {loading ? 'Memproses...' : (mode === 'setup' ? 'Simpan PIN' : 'Verifikasi')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
