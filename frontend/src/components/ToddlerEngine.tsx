import { useState, useEffect } from 'preact/hooks';
import { speakIndonesian } from '../lib/speech';
import SensorySparkle from './SensorySparkle';

interface Activity {
  id: number;
  lesson_id: number;
  type: string;
  question_json: string;
  max_score: number;
}

interface Props {
  activity: Activity;
  onComplete: (score: number) => void;
}

export default function ToddlerEngine({ activity, onComplete }: Props) {
  const [questionData, setQuestionData] = useState<any>(null);
  const [selectedOption, setSelectedOption] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<{ isCorrect: boolean; text: string } | null>(null);
  const [isLocked, setIsLocked] = useState(false);
  const [lockClickCount, setLockClickCount] = useState(0);

  useEffect(() => {
    try {
      const q = JSON.parse(activity.question_json);
      setQuestionData(q);
      setSelectedOption(null);
      setFeedback(null);
      // Auto speak prompt when activity loads
      if (q?.prompt) {
        speakIndonesian(q.prompt);
      }
    } catch (e) {
      console.error('Failed to parse toddler activity JSON', e);
    }
  }, [activity]);

  const handleOptionTap = (opt: string) => {
    if (feedback?.isCorrect) return;

    setSelectedOption(opt);
    speakIndonesian(opt);

    const isCorrect = opt === questionData.answer;
    if (isCorrect) {
      setFeedback({
        isCorrect: true,
        text: questionData.explanation || 'Pintar Sekali! 🎉',
      });
      setTimeout(() => {
        onComplete(activity.max_score);
      }, 2000);
    } else {
      setFeedback({
        isCorrect: false,
        text: 'Ayo coba lagi! 💪',
      });
      setTimeout(() => {
        setFeedback(null);
      }, 1500);
    }
  };

  const handleToggleLock = () => {
    if (lockClickCount >= 1) {
      setIsLocked(!isLocked);
      setLockClickCount(0);
    } else {
      setLockClickCount(1);
      setTimeout(() => setLockClickCount(0), 1000);
    }
  };

  if (!questionData) return null;

  return (
    <div class="toddler-engine-container theme-toddler animate-fade-in">
      <SensorySparkle active={true} />

      {/* Toddler Safety Guard Lock Button */}
      <div class="toddler-guard-bar">
        <button
          type="button"
          class={`toddler-lock-btn ${isLocked ? 'locked' : ''}`}
          onClick={handleToggleLock}
          title="Ketuk 2x untuk Kunci Modul Balita"
        >
          {isLocked ? '🔒 Modul Balita Terkunci (Ketuk 2x)' : '🔓 Ketuk 2x untuk Kunci Layar Balita'}
        </button>
      </div>

      {/* Main Toddler Activity Card */}
      <div class="toddler-card card mt-md">
        <div class="toddler-prompt-header">
          <h1 class="toddler-prompt-title">{questionData.prompt}</h1>
          <button
            type="button"
            class="btn-speak-prompt"
            onClick={() => speakIndonesian(questionData.prompt)}
            title="Dengarkan Suara"
          >
            🔊 Dengarkan
          </button>
        </div>

        {/* Options Grid */}
        <div class="toddler-options-grid mt-xl">
          {questionData.options?.map((opt: string, idx: number) => {
            const isSelected = selectedOption === opt;
            const isCorrect = feedback?.isCorrect && isSelected;
            const isWrong = feedback && !feedback.isCorrect && isSelected;

            return (
              <button
                key={idx}
                type="button"
                class={`toddler-option-card ${isSelected ? 'selected' : ''} ${isCorrect ? 'correct animate-bounce' : ''} ${isWrong ? 'wrong shake-card' : ''}`}
                disabled={isLocked || (feedback?.isCorrect ?? false)}
                onClick={() => handleOptionTap(opt)}
              >
                <span class="toddler-option-text">{opt}</span>
              </button>
            );
          })}
        </div>

        {/* Instant Toddler Feedback Banner */}
        {feedback && (
          <div class={`toddler-feedback-banner mt-xl ${feedback.isCorrect ? 'success' : 'retry'} animate-slide-up`}>
            <span class="feedback-icon">{feedback.isCorrect ? '🎉' : '💡'}</span>
            <span class="feedback-text">{feedback.text}</span>
          </div>
        )}
      </div>
    </div>
  );
}
