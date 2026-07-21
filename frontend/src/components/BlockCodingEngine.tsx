import { useState } from 'preact/hooks';

interface Props {
  availableBlocks: string[];
  onChange: (blocks: string[]) => void;
  disabled?: boolean;
}

export default function BlockCodingEngine({ availableBlocks, onChange, disabled }: Props) {
  const [assembledBlocks, setAssembledBlocks] = useState<string[]>([]);

  const addBlock = (block: string) => {
    if (disabled) return;
    const newBlocks = [...assembledBlocks, block];
    setAssembledBlocks(newBlocks);
    onChange(newBlocks);
  };

  const removeBlock = (index: number) => {
    if (disabled) return;
    const newBlocks = assembledBlocks.filter((_, i) => i !== index);
    setAssembledBlocks(newBlocks);
    onChange(newBlocks);
  };

  const moveBlock = (fromIndex: number, toIndex: number) => {
    if (disabled || toIndex < 0 || toIndex >= assembledBlocks.length) return;
    const newBlocks = [...assembledBlocks];
    const [moved] = newBlocks.splice(fromIndex, 1);
    newBlocks.splice(toIndex, 0, moved);
    setAssembledBlocks(newBlocks);
    onChange(newBlocks);
  };

  const resetBlocks = () => {
    if (disabled) return;
    setAssembledBlocks([]);
    onChange([]);
  };

  return (
    <div class="block-coding-container mt-xl">
      {/* Palette Section */}
      <div class="block-palette">
        <span class="palette-title text-muted">Pilih Blok Perintah (Klik untuk menambah):</span>
        <div class="palette-grid mt-sm">
          {availableBlocks.map((blk, idx) => (
            <button
              key={idx}
              type="button"
              class="block-btn palette-btn"
              disabled={disabled}
              onClick={() => addBlock(blk)}
            >
              + {blk}
            </button>
          ))}
        </div>
      </div>

      {/* Code Workspace Canvas */}
      <div class="block-canvas card mt-lg">
        <div class="canvas-header">
          <span class="canvas-title">
            🤖 Program Kode Robot ({assembledBlocks.length} blok)
          </span>
          {assembledBlocks.length > 0 && !disabled && (
            <button type="button" class="btn-clear-canvas" onClick={resetBlocks}>
              Hapus Semua 🗑️
            </button>
          )}
        </div>

        {assembledBlocks.length === 0 ? (
          <div class="empty-canvas text-center text-muted">
            <p>Belum ada blok kode. Klik blok di atas untuk mulai menyusun program robot!</p>
          </div>
        ) : (
          <div class="assembled-blocks-list mt-md">
            {assembledBlocks.map((blk, idx) => (
              <div key={idx} class="assembled-block-item">
                <span class="block-number">{idx + 1}</span>
                <span class="block-badge">{blk}</span>
                <div class="block-actions">
                  <button
                    type="button"
                    class="btn-block-action"
                    disabled={idx === 0 || disabled}
                    onClick={() => moveBlock(idx, idx - 1)}
                  >
                    ⬆️
                  </button>
                  <button
                    type="button"
                    class="btn-block-action"
                    disabled={idx === assembledBlocks.length - 1 || disabled}
                    onClick={() => moveBlock(idx, idx + 1)}
                  >
                    ⬇️
                  </button>
                  <button
                    type="button"
                    class="btn-block-delete"
                    disabled={disabled}
                    onClick={() => removeBlock(idx)}
                  >
                    ✕
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
