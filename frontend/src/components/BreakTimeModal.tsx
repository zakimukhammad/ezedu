import { useState } from 'preact/hooks';
import { authApi } from '../lib/api';

interface Props {
  minutesLearned: number;
  onContinue: () => void; // Extends time by 15 min after PIN check
  onGoHome?: () => void;
}

export default function BreakTimeModal({ minutesLearned, onContinue, onGoHome }: Props) {
  const [showPinInput, setShowPinInput] = useState(false);
  const [pinInput, setPinInput] = useState('');
  const [pinError, setPinError] = useState('');
  const [verifying, setVerifying] = useState(false);

  const handleGoHome = () => {
    if (onGoHome) {
      onGoHome();
    } else {
      window.location.href = '/beranda';
    }
  };

  const handleVerifyPin = async () => {
    if (pinInput.length !== 4) {
      setPinError('PIN harus 4 digit');
      return;
    }
    setVerifying(true);
    setPinError('');

    const res = await authApi.verifyPin(pinInput);
    setVerifying(false);

    if (res.error || !res.data?.valid) {
      setPinError('PIN Orang Tua salah.');
      setPinInput('');
      return;
    }

    // PIN correct — grant 15 min extension
    onContinue();
  };

  return (
    <div class="break-modal-overlay">
      <div class="break-modal-card">
        <div class="break-modal-icon">😴 🌈 ⏰</div>
        <h2 class="break-modal-title">Waktunya Istirahat!</h2>
        <p class="break-modal-message">
          Kamu sudah belajar dengan hebat hari ini! Istirahat dulu ya, mata dan pikiran butuh kesegaran. 💪
        </p>

        <div class="break-modal-stats">
          <div class="break-stat-item">
            <span class="break-stat-icon">⏱️</span>
            <div>
              <span class="break-stat-val">{minutesLearned} menit</span>
              <span class="break-stat-lbl">Waktu belajar hari ini</span>
            </div>
          </div>
        </div>

        {!showPinInput ? (
          <div class="break-modal-actions">
            <button type="button" onClick={handleGoHome} class="btn-primary break-btn-home">
              🏠 Kembali ke Beranda
            </button>
            <button
              type="button"
              onClick={() => setShowPinInput(true)}
              class="btn-ghost break-btn-extend"
            >
              🔒 Perpanjang Waktu (Orangtua)
            </button>
          </div>
        ) : (
          <div class="break-pin-box">
            <p class="break-pin-label">Masukkan PIN Orang Tua (+15 Menit):</p>
            <div class="break-pin-inputs">
              {[0, 1, 2, 3].map(i => (
                <input
                  key={i}
                  type="password"
                  maxLength={1}
                  value={pinInput[i] || ''}
                  onInput={(e: any) => {
                    const val = e.target.value.replace(/\D/g, '');
                    const newPin = pinInput.split('');
                    newPin[i] = val;
                    const joined = newPin.join('').slice(0, 4);
                    setPinInput(joined);
                    if (val && i < 3) {
                      const next = e.target.parentElement.children[i + 1] as HTMLInputElement;
                      next?.focus();
                    }
                  }}
                  onKeyDown={(e: any) => {
                    if (e.key === 'Backspace' && !pinInput[i] && i > 0) {
                      const prev = e.target.parentElement.children[i - 1] as HTMLInputElement;
                      prev?.focus();
                    }
                    if (e.key === 'Enter' && pinInput.length === 4) {
                      handleVerifyPin();
                    }
                  }}
                  class="break-pin-input"
                />
              ))}
            </div>

            {pinError && <p class="break-pin-error">{pinError}</p>}

            <div class="break-pin-actions">
              <button
                type="button"
                onClick={handleVerifyPin}
                disabled={verifying || pinInput.length !== 4}
                class="btn-primary"
                style="padding: 0.5rem 1rem; font-size: 0.9rem;"
              >
                {verifying ? 'Memverifikasi...' : 'Tambah 15 Menit ⏱️'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowPinInput(false);
                  setPinInput('');
                  setPinError('');
                }}
                class="btn-ghost"
                style="padding: 0.5rem 1rem; font-size: 0.9rem;"
              >
                Batal
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
