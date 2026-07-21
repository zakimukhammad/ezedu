import { useState } from 'preact/hooks';
import { sounds } from '../lib/sound';

interface Props {
  activity: any;
  onComplete: (score: number) => void;
}

const PALETTE = [
  { name: 'Merah', hex: '#ef4444' },
  { name: 'Kuning', hex: '#f59e0b' },
  { name: 'Hijau', hex: '#10b981' },
  { name: 'Biru', hex: '#3b82f6' },
  { name: 'Ungu', hex: '#8b5cf6' },
  { name: 'Pink', hex: '#ec4899' },
  { name: 'Cokelat', hex: '#b45309' },
  { name: 'Hitam', hex: '#1e293b' },
  { name: 'Putih', hex: '#ffffff' },
];

export default function PixelArtEngine({ activity, onComplete }: Props) {
  const [gridSize, setGridSize] = useState<8 | 12>(8);
  const [pixels, setPixels] = useState<string[]>(() => Array(8 * 8).fill('#ffffff'));
  const [selectedColor, setSelectedColor] = useState('#ef4444');
  const [isEraser, setIsEraser] = useState(false);
  const [completed, setCompleted] = useState(false);

  let question: any = {};
  try {
    question = JSON.parse(activity.question_json || '{}');
  } catch (e) {}

  const handlePixelClick = (index: number) => {
    const newColor = isEraser ? '#ffffff' : selectedColor;
    const newPixels = [...pixels];
    newPixels[index] = newColor;
    setPixels(newPixels);
    sounds.pop();
  };

  const handleGridChange = (size: 8 | 12) => {
    setGridSize(size);
    setPixels(Array(size * size).fill('#ffffff'));
    sounds.pop();
  };

  const clearGrid = () => {
    setPixels(Array(gridSize * gridSize).fill('#ffffff'));
    sounds.pop();
  };

  const handleFinish = () => {
    setCompleted(true);
    sounds.playCorrect();
    setTimeout(() => {
      onComplete(activity.max_score || 10);
    }, 400);
  };

  return (
    <div class="pixel-art-engine container">
      <div class="pixel-card glass-card p-lg animate-fade-in text-center">
        {question.hint && <p class="text-muted text-sm text-center mb-md">💡 {question.hint}</p>}

        {/* Controls */}
        <div class="pixel-toolbar mb-md flex-center gap-md wrap">
          <div class="grid-toggle">
            <span class="toolbar-label">Ukuran Grid:</span>
            <button
              type="button"
              class={`btn btn-sm ${gridSize === 8 ? 'btn-primary' : 'btn-secondary'}`}
              onClick={() => handleGridChange(8)}
            >
              8x8
            </button>
            <button
              type="button"
              class={`btn btn-sm ${gridSize === 12 ? 'btn-primary' : 'btn-secondary'}`}
              onClick={() => handleGridChange(12)}
            >
              12x12
            </button>
          </div>

          <div class="tool-actions">
            <button
              type="button"
              class={`btn btn-sm ${isEraser ? 'btn-primary' : 'btn-secondary'}`}
              onClick={() => {
                setIsEraser(!isEraser);
                sounds.pop();
              }}
            >
              🧹 Penghapus
            </button>
            <button type="button" class="btn btn-secondary btn-sm" onClick={clearGrid}>
              🗑️ Riset Grid
            </button>
          </div>
        </div>

        {/* Palette */}
        <div class="pixel-palette mb-lg flex-center gap-xs wrap">
          {PALETTE.map((c) => (
            <button
              key={c.hex}
              type="button"
              class={`palette-swatch ${selectedColor === c.hex && !isEraser ? 'active' : ''}`}
              style={{ backgroundColor: c.hex }}
              onClick={() => {
                setSelectedColor(c.hex);
                setIsEraser(false);
                sounds.pop();
              }}
              title={c.name}
            />
          ))}
        </div>

        {/* Pixel Canvas Grid */}
        <div class="pixel-grid-container mb-xl">
          <div
            class="pixel-grid"
            style={{
              gridTemplateColumns: `repeat(${gridSize}, 1fr)`,
              maxWidth: gridSize === 8 ? '320px' : '380px',
            }}
          >
            {pixels.map((color, idx) => (
              <button
                key={idx}
                type="button"
                class="pixel-cell"
                style={{ backgroundColor: color }}
                onClick={() => handlePixelClick(idx)}
              />
            ))}
          </div>
        </div>

        {/* Action button */}
        <button
          type="button"
          class="btn btn-primary btn-lg w-full"
          onClick={handleFinish}
          disabled={completed}
        >
          {completed ? 'Seni Pixel Tersimpan! 👾' : 'Simpan & Selesai 🚀'}
        </button>
      </div>
    </div>
  );
}
