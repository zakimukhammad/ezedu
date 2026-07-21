import { useState, useEffect, useRef } from 'preact/hooks';
import { sounds } from '../lib/sound';

interface Props {
  activity: {
    id: number;
    type: string;
    question_json: string;
    max_score: number;
  };
  onComplete: (score: number) => void;
}

interface MazeData {
  width: number;
  height: number;
  start: [number, number];
  goal: [number, number];
  walls: [number, number][];
}

type Direction = 'up' | 'down' | 'left' | 'right';

const DIR_EMOJI: Record<Direction, string> = {
  up: '⬆️',
  down: '⬇️',
  left: '⬅️',
  right: '➡️',
};

const DEFAULT_MAZE: MazeData = {
  width: 6,
  height: 6,
  start: [0, 0],
  goal: [5, 5],
  walls: [
    [0, 1], [1, 1], [2, 1], [3, 1], [4, 1],
    [1, 3], [2, 3], [3, 3], [4, 3], [5, 3],
    [0, 5], [1, 5], [2, 5], [3, 5], [4, 5]
  ]
};

export default function MazeLogicEngine({ activity, onComplete }: Props) {
  const question = JSON.parse(activity.question_json || '{}');
  let mazeData: MazeData = question.maze_data;
  if (
    !mazeData ||
    !mazeData.walls ||
    mazeData.walls.length === 0 ||
    mazeData.walls.some((w) => (w[0] === 1 && w[1] === 0) || (w[0] === 0 && w[1] === 1))
  ) {
    mazeData = DEFAULT_MAZE;
  }

  const [commands, setCommands] = useState<Direction[]>([]);
  const [playerPos, setPlayerPos] = useState<[number, number]>([...mazeData.start]);
  const [executing, setExecuting] = useState(false);
  const [execStep, setExecStep] = useState(-1);
  const [completed, setCompleted] = useState(false);
  const [failed, setFailed] = useState(false);
  const [attempts, setAttempts] = useState(0);
  const execRef = useRef<any>(null);

  const wallSet = new Set(mazeData.walls.map((w) => `${w[0]},${w[1]}`));

  const isWall = (x: number, y: number) => {
    return wallSet.has(`${x},${y}`) || x < 0 || y < 0 || x >= mazeData.width || y >= mazeData.height;
  };

  const addCommand = (dir: Direction) => {
    if (executing || completed) return;
    setCommands((prev) => [...prev, dir]);
    sounds.pop();
  };

  const removeLastCommand = () => {
    if (executing || completed) return;
    setCommands((prev) => prev.slice(0, -1));
    sounds.pop();
  };

  const clearCommands = () => {
    if (executing || completed) return;
    setCommands([]);
    sounds.pop();
  };

  const executeCommands = () => {
    if (commands.length === 0 || executing || completed) return;
    setExecuting(true);
    setFailed(false);
    setPlayerPos([...mazeData.start]);
    setExecStep(0);
    setAttempts((prev) => prev + 1);
  };

  useEffect(() => {
    if (!executing || execStep < 0) return;

    if (execStep >= commands.length) {
      // Done executing
      setExecuting(false);
      if (playerPos[0] === mazeData.goal[0] && playerPos[1] === mazeData.goal[1]) {
        setCompleted(true);
        sounds.playCorrect();
        setTimeout(() => {
          sounds.playFanfare();
          onComplete(activity.max_score);
        }, 800);
      } else {
        setFailed(true);
        sounds.playWrong();
      }
      return;
    }

    execRef.current = setTimeout(() => {
      const dir = commands[execStep];
      setPlayerPos((prev) => {
        let [x, y] = prev;
        if (dir === 'up') y--;
        else if (dir === 'down') y++;
        else if (dir === 'left') x--;
        else if (dir === 'right') x++;

        if (isWall(x, y)) {
          // Hit a wall — stop execution
          setExecuting(false);
          setFailed(true);
          sounds.playWrong();
          return prev;
        }

        // Check if reached goal
        if (x === mazeData.goal[0] && y === mazeData.goal[1]) {
          setCompleted(true);
          sounds.playCorrect();
          setTimeout(() => {
            sounds.playFanfare();
            onComplete(activity.max_score);
          }, 800);
        }

        return [x, y];
      });
      setExecStep((prev) => prev + 1);
    }, 400);

    return () => {
      if (execRef.current) clearTimeout(execRef.current);
    };
  }, [executing, execStep]);

  const resetGame = () => {
    setPlayerPos([...mazeData.start]);
    setCommands([]);
    setExecuting(false);
    setExecStep(-1);
    setFailed(false);
  };

  if (completed) {
    return (
      <div class="maze-engine">
        <div class="maze-results glass-card p-xl text-center animate-scale-in">
          <div class="maze-results-emoji">🎉🧩</div>
          <h2 class="mt-md">Labirin Terpecahkan!</h2>
          <p class="text-muted mt-sm">Kamu berhasil dalam {attempts} percobaan dengan {commands.length} langkah!</p>
        </div>
      </div>
    );
  }

  return (
    <div class="maze-engine">
      {/* Maze Grid */}
      <div class="maze-grid-container">
        <div
          class="maze-grid"
          style={`grid-template-columns: repeat(${mazeData.width}, 1fr); grid-template-rows: repeat(${mazeData.height}, 1fr);`}
        >
          {Array.from({ length: mazeData.height }, (_, y) =>
            Array.from({ length: mazeData.width }, (_, x) => {
              const isPlayer = playerPos[0] === x && playerPos[1] === y;
              const isGoal = mazeData.goal[0] === x && mazeData.goal[1] === y;
              const isStart = mazeData.start[0] === x && mazeData.start[1] === y;
              const isWallCell = wallSet.has(`${x},${y}`);

              return (
                <div
                  key={`${x}-${y}`}
                  class={`maze-cell ${isWallCell ? 'cell-wall' : 'cell-path'} ${isPlayer ? 'cell-player' : ''} ${isGoal ? 'cell-goal' : ''} ${isStart && !isPlayer ? 'cell-start' : ''}`}
                >
                  {isPlayer ? (
                    <span class="maze-char">🐱</span>
                  ) : isGoal ? (
                    <span class="maze-char">⭐</span>
                  ) : isWallCell ? (
                    <span class="maze-wall-icon">🧱</span>
                  ) : null}
                </div>
              );
            })
          )}
        </div>
      </div>

      {/* Command Queue */}
      <div class="maze-commands mt-md">
        <div class="command-label">Perintah ({commands.length}):</div>
        <div class="command-queue">
          {commands.length === 0 ? (
            <span class="text-muted text-sm">Tekan tombol arah untuk menambah perintah...</span>
          ) : (
            commands.map((cmd, idx) => (
              <span key={idx} class={`command-badge ${executing && idx === execStep ? 'badge-active' : executing && idx < execStep ? 'badge-done' : ''}`}>
                {DIR_EMOJI[cmd]}
              </span>
            ))
          )}
        </div>
      </div>

      {/* Feedback */}
      {failed && (
        <div class="maze-feedback mt-sm text-center animate-fade-in">
          <p class="text-danger">💥 Ups! {question.hint || 'Coba rute yang berbeda!'}</p>
        </div>
      )}

      {/* Direction Buttons */}
      <div class="maze-controls mt-md">
        <div class="direction-pad">
          <div class="dpad-row">
            <button type="button" class="dpad-btn" onClick={() => addCommand('up')} disabled={executing}>⬆️</button>
          </div>
          <div class="dpad-row">
            <button type="button" class="dpad-btn" onClick={() => addCommand('left')} disabled={executing}>⬅️</button>
            <button type="button" class="dpad-btn dpad-center" disabled>🐱</button>
            <button type="button" class="dpad-btn" onClick={() => addCommand('right')} disabled={executing}>➡️</button>
          </div>
          <div class="dpad-row">
            <button type="button" class="dpad-btn" onClick={() => addCommand('down')} disabled={executing}>⬇️</button>
          </div>
        </div>
      </div>

      {/* Action Buttons */}
      <div class="maze-actions mt-md flex-center gap-md">
        <button type="button" class="btn btn-secondary" onClick={removeLastCommand} disabled={executing || commands.length === 0}>
          ↩️ Hapus
        </button>
        <button type="button" class="btn btn-secondary" onClick={clearCommands} disabled={executing || commands.length === 0}>
          🗑️ Reset
        </button>
        <button type="button" class="btn btn-primary btn-lg" onClick={executeCommands} disabled={executing || commands.length === 0}>
          ▶️ Jalankan!
        </button>
      </div>

      {failed && (
        <button type="button" class="btn btn-primary w-full mt-md" onClick={resetGame}>
          🔄 Coba Lagi
        </button>
      )}
    </div>
  );
}
