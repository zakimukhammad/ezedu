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

interface Problem {
  a: number;
  b: number;
  op: string;
  answer: number;
  display: string;
}

export default function MathRacerEngine({ activity, onComplete }: Props) {
  const question = JSON.parse(activity.question_json || '{}');
  const timeLimit = question.time_limit || 60;
  const maxNum = question.max_number || 50;
  const ops: string[] = question.operations || ['add', 'subtract'];

  const [started, setStarted] = useState(false);
  const [finished, setFinished] = useState(false);
  const [timeLeft, setTimeLeft] = useState(timeLimit);
  const [score, setScore] = useState(0);
  const [streak, setStreak] = useState(0);
  const [bestStreak, setBestStreak] = useState(0);
  const [problem, setProblem] = useState<Problem | null>(null);
  const [options, setOptions] = useState<number[]>([]);
  const [feedback, setFeedback] = useState<'correct' | 'wrong' | null>(null);
  const [questionsAnswered, setQuestionsAnswered] = useState(0);
  const timerRef = useRef<any>(null);

  const generateProblem = (): { prob: Problem; opts: number[] } => {
    const op = ops[Math.floor(Math.random() * ops.length)];
    let a: number, b: number, answer: number, display: string;

    if (op === 'add') {
      a = Math.floor(Math.random() * maxNum) + 1;
      b = Math.floor(Math.random() * maxNum) + 1;
      answer = a + b;
      display = `${a} + ${b} = ?`;
    } else if (op === 'subtract') {
      a = Math.floor(Math.random() * maxNum) + 10;
      b = Math.floor(Math.random() * Math.min(a, maxNum)) + 1;
      answer = a - b;
      display = `${a} − ${b} = ?`;
    } else if (op === 'multiply') {
      a = Math.floor(Math.random() * 12) + 1;
      b = Math.floor(Math.random() * 12) + 1;
      answer = a * b;
      display = `${a} × ${b} = ?`;
    } else {
      b = Math.floor(Math.random() * 12) + 1;
      answer = Math.floor(Math.random() * 12) + 1;
      a = b * answer;
      display = `${a} ÷ ${b} = ?`;
    }

    // Generate 3 wrong options + 1 correct
    const wrongSet = new Set<number>();
    while (wrongSet.size < 3) {
      const offset = Math.floor(Math.random() * 10) - 5;
      const wrong = answer + (offset === 0 ? 1 : offset);
      if (wrong !== answer && wrong >= 0) {
        wrongSet.add(wrong);
      }
    }

    const allOpts = [...wrongSet, answer];
    // Shuffle
    for (let i = allOpts.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [allOpts[i], allOpts[j]] = [allOpts[j], allOpts[i]];
    }

    return {
      prob: { a, b, op, answer, display },
      opts: allOpts,
    };
  };

  const startGame = () => {
    setStarted(true);
    setScore(0);
    setStreak(0);
    setBestStreak(0);
    setQuestionsAnswered(0);
    setTimeLeft(timeLimit);
    nextProblem();

    timerRef.current = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timerRef.current);
          setFinished(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  };

  const nextProblem = () => {
    const { prob, opts } = generateProblem();
    setProblem(prob);
    setOptions(opts);
    setFeedback(null);
  };

  const handleAnswer = (chosen: number) => {
    if (finished || feedback) return;

    setQuestionsAnswered((q) => q + 1);

    if (chosen === problem?.answer) {
      const streakBonus = streak >= 3 ? 5 : 0;
      const pointsEarned = 10 + streakBonus;
      setScore((prev) => prev + pointsEarned);
      setStreak((prev) => {
        const next = prev + 1;
        setBestStreak((best) => Math.max(best, next));
        return next;
      });
      setFeedback('correct');
      sounds.playCorrect();
    } else {
      setStreak(0);
      setFeedback('wrong');
      sounds.playWrong();
    }

    setTimeout(() => {
      if (!finished) {
        nextProblem();
      }
    }, 300);
  };

  useEffect(() => {
    if (finished) {
      sounds.playFanfare();
      setTimeout(() => {
        onComplete(Math.min(score, activity.max_score));
      }, 800);
    }
  }, [finished]);

  useEffect(() => {
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, []);

  const progressPercent = Math.round((timeLeft / timeLimit) * 100);
  const timerColor = timeLeft <= 10 ? 'var(--color-danger)' : timeLeft <= 20 ? 'var(--color-warning)' : 'var(--color-success)';

  if (!started) {
    return (
      <div class="math-racer-engine">
        <div class="racer-intro glass-card p-xl text-center animate-fade-in">
          <div class="racer-intro-emoji">⏱️🧮</div>
          <h2 class="mt-md">Math Racer</h2>
          <p class="text-muted mt-sm">Pecahkan soal matematika sebanyak mungkin dalam <strong>{timeLimit} detik</strong>!</p>
          <div class="racer-rules mt-lg">
            <div class="rule-item">🎯 Setiap jawaban benar = <strong>10 poin</strong></div>
            <div class="rule-item">🔥 Streak 3+ berturut-turut = <strong>+5 bonus</strong></div>
            <div class="rule-item">⚡ Semakin cepat, semakin banyak poin!</div>
          </div>
          <button class="btn btn-primary btn-lg w-full mt-xl" onClick={startGame}>
            Mulai Race! 🚀
          </button>
        </div>
      </div>
    );
  }

  if (finished) {
    return (
      <div class="math-racer-engine">
        <div class="racer-results glass-card p-xl text-center animate-scale-in">
          <div class="racer-results-emoji">🏆</div>
          <h2 class="mt-md">Race Selesai!</h2>
          <div class="racer-stats mt-lg">
            <div class="stat-item">
              <span class="stat-value">{score}</span>
              <span class="stat-label">Poin</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{questionsAnswered}</span>
              <span class="stat-label">Soal</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">🔥 {bestStreak}</span>
              <span class="stat-label">Best Streak</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div class="math-racer-engine">
      {/* Timer Bar */}
      <div class="racer-timer-container">
        <div
          class="racer-timer-bar"
          style={`width: ${progressPercent}%; background: ${timerColor}; transition: width 1s linear;`}
        />
        <span class="racer-timer-text">{timeLeft}s</span>
      </div>

      {/* Score & Streak */}
      <div class="racer-hud mt-md">
        <div class="hud-item">
          <span class="hud-label">Poin</span>
          <span class="hud-value">{score}</span>
        </div>
        <div class="hud-item">
          <span class="hud-label">Streak</span>
          <span class="hud-value">{streak >= 3 ? '🔥' : ''} {streak}</span>
        </div>
        <div class="hud-item">
          <span class="hud-label">Soal</span>
          <span class="hud-value">#{questionsAnswered + 1}</span>
        </div>
      </div>

      {/* Problem Display */}
      {problem && (
        <div class={`racer-problem-card glass-card mt-lg p-xl text-center animate-fade-in ${feedback === 'correct' ? 'flash-green' : feedback === 'wrong' ? 'flash-red' : ''}`}>
          <h1 class="racer-problem-text">{problem.display}</h1>

          <div class="racer-options mt-xl">
            {options.map((opt, idx) => (
              <button
                key={`${questionsAnswered}-${idx}`}
                type="button"
                class="racer-option-btn"
                onClick={() => handleAnswer(opt)}
              >
                {opt}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
