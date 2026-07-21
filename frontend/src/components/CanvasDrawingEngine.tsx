import { useState, useRef, useEffect } from 'preact/hooks';
import { sounds } from '../lib/sound';

interface Props {
  activity: any;
  onComplete: (score: number) => void;
}

const COLORS = [
  { name: 'Merah', hex: '#ef4444' },
  { name: 'Kuning', hex: '#f59e0b' },
  { name: 'Hijau', hex: '#10b981' },
  { name: 'Biru', hex: '#3b82f6' },
  { name: 'Ungu', hex: '#8b5cf6' },
  { name: 'Pink', hex: '#ec4899' },
  { name: 'Hitam', hex: '#1e293b' },
];

const BRUSH_SIZES = [
  { name: 'Tipis', size: 4 },
  { name: 'Sedang', size: 10 },
  { name: 'Tebal', size: 20 },
];

const STAMPS = ['⭐️', '🎨', '🌈', '🐱', '🌸', '🚀', '❤️', '☀️'];

export default function CanvasDrawingEngine({ activity, onComplete }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [color, setColor] = useState('#ef4444');
  const [brushSize, setBrushSize] = useState(10);
  const [mode, setMode] = useState<'draw' | 'erase' | 'stamp'>('draw');
  const [selectedStamp, setSelectedStamp] = useState('⭐️');
  const [isDrawing, setIsDrawing] = useState(false);
  const [completed, setCompleted] = useState(false);

  let question: any = {};
  try {
    question = JSON.parse(activity.question_json || '{}');
  } catch (e) {}

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas dimensions
    canvas.width = canvas.parentElement?.clientWidth || 600;
    canvas.height = 360;

    // Fill white background
    ctx.fillStyle = '#ffffff';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
  }, []);

  const getCoordinates = (e: MouseEvent | TouchEvent) => {
    const canvas = canvasRef.current;
    if (!canvas) return { x: 0, y: 0 };
    const rect = canvas.getBoundingClientRect();
    let clientX = 0;
    let clientY = 0;

    if ('touches' in e && e.touches.length > 0) {
      clientX = e.touches[0].clientX;
      clientY = e.touches[0].clientY;
    } else if ('clientX' in e) {
      clientX = (e as MouseEvent).clientX;
      clientY = (e as MouseEvent).clientY;
    }

    return {
      x: clientX - rect.left,
      y: clientY - rect.top,
    };
  };

  const startDrawing = (e: any) => {
    const { x, y } = getCoordinates(e);
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    if (mode === 'stamp') {
      ctx.font = '36px sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillText(selectedStamp, x, y);
      sounds.pop();
      return;
    }

    setIsDrawing(true);
    ctx.beginPath();
    ctx.moveTo(x, y);
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
    ctx.strokeStyle = mode === 'erase' ? '#ffffff' : color;
    ctx.lineWidth = mode === 'erase' ? brushSize * 2 : brushSize;
  };

  const draw = (e: any) => {
    if (!isDrawing || mode === 'stamp') return;
    const { x, y } = getCoordinates(e);
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.lineTo(x, y);
    ctx.stroke();
  };

  const stopDrawing = () => {
    setIsDrawing(false);
  };

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.fillStyle = '#ffffff';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
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
    <div class="canvas-drawing-engine container">
      <div class="drawing-card glass-card p-lg animate-fade-in">
        {question.hint && <p class="text-muted text-sm text-center mb-md">💡 {question.hint}</p>}

        {/* Toolbar Controls */}
        <div class="drawing-toolbar mb-md">
          {/* Colors */}
          <div class="toolbar-section">
            <span class="toolbar-label">Warna:</span>
            <div class="color-picker-grid">
              {COLORS.map((c) => (
                <button
                  key={c.hex}
                  type="button"
                  class={`color-btn ${color === c.hex && mode === 'draw' ? 'active' : ''}`}
                  style={{ backgroundColor: c.hex }}
                  onClick={() => {
                    setColor(c.hex);
                    setMode('draw');
                    sounds.pop();
                  }}
                  title={c.name}
                />
              ))}
            </div>
          </div>

          {/* Brush Sizes */}
          <div class="toolbar-section">
            <span class="toolbar-label">Ukuran Kuas:</span>
            <div class="brush-picker">
              {BRUSH_SIZES.map((b) => (
                <button
                  key={b.size}
                  type="button"
                  class={`brush-btn ${brushSize === b.size ? 'active' : ''}`}
                  onClick={() => {
                    setBrushSize(b.size);
                    sounds.pop();
                  }}
                >
                  <span class="brush-dot" style={{ width: `${b.size}px`, height: `${b.size}px` }}></span>
                  {b.name}
                </button>
              ))}
            </div>
          </div>

          {/* Eraser & Clear */}
          <div class="toolbar-section">
            <button
              type="button"
              class={`btn ${mode === 'erase' ? 'btn-primary' : 'btn-secondary'} btn-sm`}
              onClick={() => {
                setMode('erase');
                sounds.pop();
              }}
            >
              🧹 Penghapus
            </button>
            <button type="button" class="btn btn-secondary btn-sm" onClick={clearCanvas}>
              🗑️ Hapus Semua
            </button>
          </div>
        </div>

        {/* Stamps Section */}
        <div class="stamps-bar mb-md">
          <span class="toolbar-label">Stempel Stiker:</span>
          <div class="stamps-grid">
            {STAMPS.map((s) => (
              <button
                key={s}
                type="button"
                class={`stamp-btn ${mode === 'stamp' && selectedStamp === s ? 'active' : ''}`}
                onClick={() => {
                  setSelectedStamp(s);
                  setMode('stamp');
                  sounds.pop();
                }}
              >
                {s}
              </button>
            ))}
          </div>
        </div>

        {/* Interactive Canvas */}
        <div class="canvas-container">
          <canvas
            ref={canvasRef}
            class="drawing-canvas"
            onMouseDown={startDrawing}
            onMouseMove={draw}
            onMouseUp={stopDrawing}
            onMouseLeave={stopDrawing}
            onTouchStart={startDrawing}
            onTouchMove={draw}
            onTouchEnd={stopDrawing}
          />
        </div>

        {/* Submit Drawing Button */}
        <div class="drawing-actions text-center mt-lg">
          <button
            type="button"
            class="btn btn-primary btn-lg w-full"
            onClick={handleFinish}
            disabled={completed}
          >
            {completed ? 'Karya Tersimpan! 🎉' : 'Simpan & Selesai Melukis ✨'}
          </button>
        </div>
      </div>
    </div>
  );
}
