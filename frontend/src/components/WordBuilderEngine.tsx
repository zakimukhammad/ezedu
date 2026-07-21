import { useState, useEffect } from 'preact/hooks';
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

interface WordEntry {
  word: string;
  clue: string;
}

export default function WordBuilderEngine({ activity, onComplete }: Props) {
  const question = JSON.parse(activity.question_json || '{}');
  const words: WordEntry[] = question.words || [];

  const [currentWordIdx, setCurrentWordIdx] = useState(0);
  const [shuffledLetters, setShuffledLetters] = useState<string[]>([]);
  const [selectedLetters, setSelectedLetters] = useState<number[]>([]);
  const [feedback, setFeedback] = useState<'correct' | 'wrong' | null>(null);
  const [correctCount, setCorrectCount] = useState(0);
  const [finished, setFinished] = useState(false);
  const [shake, setShake] = useState(false);

  const currentWord = words[currentWordIdx];

  useEffect(() => {
    if (currentWord) {
      shuffleWord(currentWord.word);
    }
  }, [currentWordIdx]);

  const shuffleWord = (word: string) => {
    const letters = word.toUpperCase().split('');
    // Fisher-Yates shuffle
    for (let i = letters.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [letters[i], letters[j]] = [letters[j], letters[i]];
    }
    // Ensure it's not the same as the original
    if (letters.join('') === word.toUpperCase() && letters.length > 1) {
      [letters[0], letters[1]] = [letters[1], letters[0]];
    }
    setShuffledLetters(letters);
    setSelectedLetters([]);
    setFeedback(null);
  };

  const toggleLetter = (idx: number) => {
    if (feedback === 'correct') return;

    if (selectedLetters.includes(idx)) {
      setSelectedLetters(selectedLetters.filter((i) => i !== idx));
      sounds.pop();
    } else {
      setSelectedLetters([...selectedLetters, idx]);
      sounds.pop();
    }
  };

  const getCurrentAnswer = () => {
    return selectedLetters.map((i) => shuffledLetters[i]).join('');
  };

  const checkAnswer = () => {
    const answer = getCurrentAnswer();
    if (answer.toUpperCase() === currentWord.word.toUpperCase()) {
      setFeedback('correct');
      setCorrectCount((prev) => prev + 1);
      sounds.playCorrect();

      setTimeout(() => {
        if (currentWordIdx + 1 < words.length) {
          setCurrentWordIdx((prev) => prev + 1);
        } else {
          setFinished(true);
        }
      }, 800);
    } else {
      setFeedback('wrong');
      setShake(true);
      sounds.playWrong();
      setTimeout(() => {
        setShake(false);
        setFeedback(null);
      }, 600);
    }
  };

  const resetWord = () => {
    setSelectedLetters([]);
    setFeedback(null);
  };

  useEffect(() => {
    if (finished) {
      sounds.playFanfare();
      const scorePerWord = Math.floor(activity.max_score / Math.max(words.length, 1));
      setTimeout(() => {
        onComplete(correctCount * scorePerWord);
      }, 600);
    }
  }, [finished]);

  if (words.length === 0) {
    return <div class="text-center text-muted p-xl">Tidak ada kata tersedia.</div>;
  }

  if (finished) {
    return (
      <div class="word-builder-engine">
        <div class="word-results glass-card p-xl text-center animate-scale-in">
          <div class="word-results-emoji">📝✨</div>
          <h2 class="mt-md">Permainan Selesai!</h2>
          <div class="word-stats mt-lg">
            <div class="stat-item">
              <span class="stat-value">{correctCount}</span>
              <span class="stat-label">Benar</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{words.length}</span>
              <span class="stat-label">Total Kata</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div class="word-builder-engine">
      {/* Progress */}
      <div class="word-progress mb-md">
        <span class="text-muted">Kata {currentWordIdx + 1} dari {words.length}</span>
        <div class="word-progress-dots">
          {words.map((_, idx) => (
            <span
              key={idx}
              class={`progress-dot ${idx < currentWordIdx ? 'dot-done' : idx === currentWordIdx ? 'dot-active' : ''}`}
            />
          ))}
        </div>
      </div>

      {/* Clue Card */}
      <div class="word-clue-card glass-card p-lg text-center animate-fade-in">
        <p class="word-clue-text">{currentWord.clue}</p>
      </div>

      {/* Answer display */}
      <div class={`word-answer-display mt-lg ${shake ? 'shake-card' : ''}`}>
        {currentWord.word.split('').map((_, idx) => (
          <div key={idx} class={`answer-slot ${idx < selectedLetters.length ? 'slot-filled' : ''} ${feedback === 'correct' ? 'slot-correct' : feedback === 'wrong' ? 'slot-wrong' : ''}`}>
            {idx < selectedLetters.length ? shuffledLetters[selectedLetters[idx]] : ''}
          </div>
        ))}
      </div>

      {/* Shuffled letter tiles */}
      <div class="word-tiles mt-lg">
        {shuffledLetters.map((letter, idx) => (
          <button
            key={idx}
            type="button"
            class={`word-tile ${selectedLetters.includes(idx) ? 'tile-used' : ''}`}
            onClick={() => toggleLetter(idx)}
            disabled={feedback === 'correct'}
          >
            {letter}
          </button>
        ))}
      </div>

      {/* Actions */}
      <div class="word-actions mt-lg flex-center gap-md">
        <button type="button" class="btn btn-secondary" onClick={resetWord} disabled={feedback === 'correct'}>
          🔄 Ulang
        </button>
        <button
          type="button"
          class="btn btn-primary btn-lg"
          onClick={checkAnswer}
          disabled={selectedLetters.length !== currentWord.word.length || feedback === 'correct'}
        >
          Periksa ✓
        </button>
      </div>
    </div>
  );
}
