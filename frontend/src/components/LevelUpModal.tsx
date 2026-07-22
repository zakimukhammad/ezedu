import { sounds } from '../lib/sound';

interface Props {
  newLevel: number;
  onClose: () => void;
}

export default function LevelUpModal({ newLevel, onClose }: Props) {
  return (
    <div class="modal-backdrop flex-center animate-fade-in" style="z-index: 10000;">
      <div class="modal-card glass-card p-2xl text-center animate-scale-in" style="max-width: 440px; background: linear-gradient(135deg, rgba(30, 27, 75, 0.95), rgba(15, 23, 42, 0.95)); border: 2px solid var(--color-secondary);">
        <div class="level-up-icon animate-bounce" style="font-size: 4rem;">🏆✨</div>
        <h1 class="text-secondary mt-md" style="font-size: 2rem; font-weight: 800;">LEVEL UP!</h1>
        <p class="text-muted mt-xs">Selamat! Kamu berhasil mencapai</p>
        <div class="level-badge mt-md p-md" style="background: rgba(245, 158, 11, 0.2); border-radius: var(--radius-lg); border: 1px solid var(--color-secondary);">
          <span style="font-size: 2.2rem; font-weight: 800; color: var(--color-secondary-light);">
            Level {newLevel} 🎉
          </span>
        </div>
        <p class="mt-md text-sm text-muted">Tetap semangat belajar untuk membuka petualangan baru!</p>

        <button
          type="button"
          class="btn btn-primary btn-lg w-full mt-xl"
          onClick={onClose}
        >
          Lanjutkan Petualangan! 🚀
        </button>
      </div>
    </div>
  );
}
