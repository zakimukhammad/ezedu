import { useState, useEffect } from 'preact/hooks';
import { lessonsApi, activitiesApi } from '../lib/api';
import { sounds } from '../lib/sound';
import BlockCodingEngine from './BlockCodingEngine';

interface Activity {
  id: number;
  lesson_id: number;
  type: string;
  sort_order: number;
  question_json: string;
  max_score: number;
}

interface Lesson {
  id: number;
  category_id: number;
  title: string;
  description: string;
  content_json: string;
  estimated_minutes: number;
  xp_reward: number;
}

interface Props {
  lessonId: number;
}

interface Child {
  id: number;
  name: string;
  age_group: string;
}

export default function QuizEngine({ lessonId }: Props) {
  const [lesson, setLesson] = useState<Lesson | null>(null);
  const [activities, setActivities] = useState<Activity[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [child, setChild] = useState<Child | null>(null);
  const [loading, setLoading] = useState(true);

  // Activity state
  const [selectedChoice, setSelectedChoice] = useState<string>('');
  const [fillAnswer, setFillAnswer] = useState<string>('');
  const [dragItems, setDragItems] = useState<string[]>([]);
  const [blockAnswer, setBlockAnswer] = useState<string[]>([]);
  const [submitted, setSubmitted] = useState(false);
  const [feedback, setFeedback] = useState<{ isCorrect: boolean; text: string; hint?: string; explanation?: string } | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [showHint, setShowHint] = useState(false);
  const [attemptCount, setAttemptCount] = useState(1);

  // Lesson summary state
  const [totalScore, setTotalScore] = useState(0);
  const [maxPossibleScore, setMaxPossibleScore] = useState(0);
  const [startTime] = useState(Date.now());
  const [completed, setCompleted] = useState(false);
  const [xpEarned, setXpEarned] = useState(0);

  useEffect(() => {
    const stored = sessionStorage.getItem('ezedu_child');
    if (stored) {
      setChild(JSON.parse(stored));
    }
    loadLessonData();
  }, [lessonId]);

  const loadLessonData = async () => {
    setLoading(true);
    const { data } = await lessonsApi.getById(lessonId);
    setLoading(false);

    if (data?.lesson) {
      setLesson(data.lesson);
      setActivities(data.activities || []);

      let totalMax = 0;
      (data.activities || []).forEach((a: Activity) => {
        totalMax += a.max_score;
      });
      setMaxPossibleScore(totalMax);

      if (data.activities?.length > 0) {
        initActivityState(data.activities[0]);
      }
    }
  };

  const initActivityState = (act: Activity) => {
    setSelectedChoice('');
    setFillAnswer('');
    setBlockAnswer([]);
    setSubmitted(false);
    setFeedback(null);
    setShowHint(false);
    setAttemptCount(1);

    try {
      const q = JSON.parse(act.question_json);
      if (act.type === 'drag_drop' || act.type === 'sequencing') {
        setDragItems(q.items ? [...q.items] : []);
      }
    } catch (e) {
      console.error('Failed to parse question JSON', e);
    }
  };

  const currentActivity = activities[currentIndex];
  let currentQuestion: any = null;
  if (currentActivity) {
    try {
      currentQuestion = JSON.parse(currentActivity.question_json);
    } catch (e) {}
  }

  // Handle Drag & Drop item reordering (move item up/down)
  const moveItem = (fromIndex: number, toIndex: number) => {
    if (toIndex < 0 || toIndex >= dragItems.length) return;
    const newItems = [...dragItems];
    const [moved] = newItems.splice(fromIndex, 1);
    newItems.splice(toIndex, 0, moved);
    setDragItems(newItems);
  };

  const handleSubmit = async () => {
    if (!currentActivity || !child) return;

    let answer: any = null;
    if (currentActivity.type === 'multiple_choice') {
      if (!selectedChoice) return;
      answer = selectedChoice;
    } else if (currentActivity.type === 'fill_blank') {
      if (!fillAnswer.trim()) return;
      answer = fillAnswer.trim();
    } else if (currentActivity.type === 'drag_drop' || currentActivity.type === 'sequencing') {
      answer = dragItems;
    } else if (currentActivity.type === 'block_code') {
      if (blockAnswer.length === 0) return;
      answer = blockAnswer;
    }

    setSubmitting(true);
    const { data } = await activitiesApi.submit(
      currentActivity.id,
      child.id,
      answer,
      attemptCount
    );
    setSubmitting(false);

    if (data) {
      setSubmitted(true);
      setFeedback({
        isCorrect: data.is_correct,
        text: data.feedback,
        hint: data.hint,
        explanation: data.explanation,
      });

      if (data.is_correct) {
        sounds.playCorrect();
        setTotalScore((prev) => prev + data.score);
      } else {
        sounds.playWrong();
        setAttemptCount((prev) => prev + 1);
      }
    }
  };

  const handleNextActivity = async () => {
    if (currentIndex + 1 < activities.length) {
      const nextIdx = currentIndex + 1;
      setCurrentIndex(nextIdx);
      initActivityState(activities[nextIdx]);
    } else {
      finishLesson();
    }
  };

  const finishLesson = async () => {
    if (!lesson || !child) return;

    const timeSpentSec = Math.round((Date.now() - startTime) / 1000);
    const { data } = await lessonsApi.complete(
      lesson.id,
      child.id,
      totalScore,
      maxPossibleScore,
      timeSpentSec
    );

    setXpEarned(data?.xp_earned || lesson.xp_reward);
    sounds.playFanfare();
    setCompleted(true);
  };

  if (loading) {
    return (
      <div class="quiz-loading flex-center">
        <div class="loading-spinner"></div>
        <p>Menyiapkan ruang belajar...</p>
      </div>
    );
  }

  if (!lesson || activities.length === 0) {
    return (
      <div class="quiz-error text-center">
        <h2>Belum ada soal untuk pelajaran ini 🚀</h2>
        <a href="/beranda" class="btn btn-primary mt-lg">Kembali ke Beranda</a>
      </div>
    );
  }

  // Summary Celebration Screen
  if (completed) {
    return (
      <div class="quiz-summary-card card animate-scale-in text-center">
        <div class="celebration-emoji">🏆</div>
        <h2>Pelajaran Selesai!</h2>
        <p class="summary-subtitle text-muted mt-xs">Kamu sudah menyelesaikan <strong>{lesson.title}</strong></p>

        <div class="summary-stats-grid mt-xl">
          <div class="summary-stat-box">
            <span class="stat-number text-gradient">+{xpEarned} XP</span>
            <span class="stat-label">Hadiah XP</span>
          </div>
          <div class="summary-stat-box">
            <span class="stat-number">{totalScore} / {maxPossibleScore}</span>
            <span class="stat-label">Total Skor</span>
          </div>
        </div>

        <div class="mt-2xl flex-center gap-md">
          <a href="/beranda" class="btn btn-primary btn-lg w-full" id="finish-to-dashboard">
            Lanjut Belajar 🚀
          </a>
        </div>
      </div>
    );
  }

  const progressPercent = Math.round(((currentIndex + 1) / activities.length) * 100);

  return (
    <div class="quiz-engine-view container">
      {/* Top Header & Progress */}
      <div class="quiz-top-bar">
        <a href="/beranda" class="close-quiz-btn">✕ Keluar</a>
        <div class="quiz-progress-container">
          <div class="quiz-progress-bar" style={`width: ${progressPercent}%`}></div>
        </div>
        <span class="quiz-progress-text">{currentIndex + 1} dari {activities.length} Soal</span>
      </div>

      {/* Main Question Card */}
      <div class={`quiz-card card mt-lg animate-fade-in ${feedback && !feedback.isCorrect ? 'shake-card' : ''}`}>
        <div class="quiz-type-badge">
          {currentActivity.type === 'multiple_choice'
            ? '💡 Pilihan Ganda'
            : currentActivity.type === 'fill_blank'
            ? '✏️ Ketik Jawaban'
            : currentActivity.type === 'block_code'
            ? '🤖 Blok Koding'
            : '🧩 Seret & Urutkan'}
        </div>

        <h2 class="quiz-prompt mt-md">{currentQuestion?.prompt}</h2>

        {/* Multiple Choice Component */}
        {currentActivity.type === 'multiple_choice' && (
          <div class="options-grid mt-xl">
            {currentQuestion?.options?.map((opt: string, idx: number) => {
              const isSelected = selectedChoice === opt;
              return (
                <button
                  key={idx}
                  class={`option-card ${isSelected ? 'option-selected' : ''}`}
                  disabled={submitted && feedback?.isCorrect}
                  onClick={() => !submitted && setSelectedChoice(opt)}
                >
                  <span class="option-label">{String.fromCharCode(65 + idx)}</span>
                  <span class="option-text">{opt}</span>
                </button>
              );
            })}
          </div>
        )}

        {/* Fill in the Blank Component */}
        {currentActivity.type === 'fill_blank' && (
          <div class="fill-blank-container mt-xl">
            <input
              type="text"
              class="form-input fill-input"
              placeholder="Ketik jawabanmu di sini..."
              value={fillAnswer}
              disabled={submitted && feedback?.isCorrect}
              onInput={(e) => setFillAnswer((e.target as HTMLInputElement).value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSubmit()}
              autoFocus
            />
          </div>
        )}

        {/* Visual Block Coding Component */}
        {currentActivity.type === 'block_code' && (
          <BlockCodingEngine
            availableBlocks={currentQuestion?.available_blocks || []}
            onChange={setBlockAnswer}
            disabled={submitted && feedback?.isCorrect}
          />
        )}

        {/* Drag & Drop / Sequencing Component */}
        {(currentActivity.type === 'drag_drop' || currentActivity.type === 'sequencing') && (
          <div class="reorder-container mt-xl">
            <p class="reorder-instruction text-muted">
              Gunakan tombol panah ⬆️ ⬇️ untuk mengurutkan item dari atas ke bawah:
            </p>
            <div class="reorder-list mt-md">
              {dragItems.map((item, idx) => (
                <div key={idx} class="reorder-item card">
                  <span class="item-index">{idx + 1}</span>
                  <span class="item-text">{item}</span>
                  <div class="reorder-controls">
                    <button
                      type="button"
                      class="btn-arrow"
                      disabled={idx === 0 || (submitted && feedback?.isCorrect)}
                      onClick={() => moveItem(idx, idx - 1)}
                    >
                      ⬆️
                    </button>
                    <button
                      type="button"
                      class="btn-arrow"
                      disabled={idx === dragItems.length - 1 || (submitted && feedback?.isCorrect)}
                      onClick={() => moveItem(idx, idx + 1)}
                    >
                      ⬇️
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Instant Feedback Panel */}
        {feedback && (
          <div class={`feedback-panel mt-xl ${feedback.isCorrect ? 'feedback-success' : 'feedback-wrong'} animate-slide-up`}>
            <div class="feedback-header">
              <span class="feedback-icon">{feedback.isCorrect ? '🎉' : '💡'}</span>
              <strong class="feedback-title">{feedback.text}</strong>
            </div>

            {feedback.explanation && (
              <p class="feedback-explanation mt-xs">{feedback.explanation}</p>
            )}

            {!feedback.isCorrect && feedback.hint && (
              <div class="hint-section mt-sm">
                {!showHint ? (
                  <button type="button" class="btn-hint" onClick={() => setShowHint(true)}>
                    🔍 Lihat Petunjuk (Hint)
                  </button>
                ) : (
                  <p class="hint-text"><strong>Petunjuk:</strong> {feedback.hint}</p>
                )}
              </div>
            )}
          </div>
        )}

        {/* Submit & Navigation Buttons */}
        <div class="quiz-footer mt-xl">
          {!submitted || (feedback && !feedback.isCorrect) ? (
            <button
              class="btn btn-primary btn-lg w-full"
              disabled={
                submitting ||
                (currentActivity.type === 'multiple_choice' && !selectedChoice) ||
                (currentActivity.type === 'fill_blank' && !fillAnswer.trim()) ||
                (currentActivity.type === 'block_code' && blockAnswer.length === 0)
              }
              onClick={handleSubmit}
              id="quiz-submit-btn"
            >
              {submitting ? 'Memeriksa...' : 'Jawab Now 🚀'}
            </button>
          ) : (
            <button
              class="btn btn-primary btn-lg w-full animate-bounce"
              onClick={handleNextActivity}
              id="quiz-next-btn"
            >
              {currentIndex + 1 < activities.length ? 'Soal Berikutnya ➔' : 'Selesaikan Pelajaran 🎉'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
