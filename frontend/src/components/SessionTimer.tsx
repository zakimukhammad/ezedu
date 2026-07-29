import { useState, useEffect } from 'preact/hooks';
import { childrenApi } from '../lib/api';
import BreakTimeModal from './BreakTimeModal';

interface Child {
  id: number;
  name: string;
  age_group: string;
}

export default function SessionTimer() {
  const [child, setChild] = useState<Child | null>(null);
  const [hasLimit, setHasLimit] = useState(false);
  const [dailyLimitMin, setDailyLimitMin] = useState(0);
  const [timeUsedSec, setTimeUsedSec] = useState(0);
  const [remainingSec, setRemainingSec] = useState<number | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [isMinimized, setIsMinimized] = useState(false);

  useEffect(() => {
    const stored = sessionStorage.getItem('ezedu_child');
    if (stored) {
      try {
        const c = JSON.parse(stored);
        setChild(c);
        fetchTimeLimit(c.id);
      } catch (e) {
        // ignore
      }
    }
  }, []);

  const fetchTimeLimit = async (childId: number) => {
    const res = await childrenApi.getRemainingTime(childId);
    if (res.data?.has_limit) {
      setHasLimit(true);
      setDailyLimitMin(res.data.daily_limit_min || 0);
      setTimeUsedSec(res.data.time_used_sec || 0);
      const rem = res.data.remaining_sec ?? 0;
      setRemainingSec(rem);

      if (rem <= 0) {
        setShowModal(true);
      }
    } else {
      setHasLimit(false);
    }
  };

  useEffect(() => {
    if (!hasLimit || remainingSec === null || remainingSec <= 0) return;

    const interval = setInterval(() => {
      setRemainingSec(prev => {
        if (prev === null || prev <= 1) {
          clearInterval(interval);
          setShowModal(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [hasLimit, remainingSec !== null && remainingSec > 0]);

  const handleExtend15Min = () => {
    // Add 15 minutes (900s)
    setRemainingSec(prev => (prev || 0) + 900);
    setShowModal(false);
  };

  if (!child || !hasLimit || remainingSec === null) {
    return null;
  }

  const formatRemaining = (sec: number) => {
    if (sec <= 0) return '00:00';
    const mins = Math.floor(sec / 60);
    const secs = sec % 60;
    if (mins >= 60) {
      const hrs = Math.floor(mins / 60);
      const rMins = mins % 60;
      return `${hrs}j ${rMins}m`;
    }
    return `${mins}m ${secs < 10 ? '0' : ''}${secs}s`;
  };

  const minutesLearned = Math.round(timeUsedSec / 60);

  const isWarning = remainingSec > 0 && remainingSec <= 300; // <= 5 min
  const isDanger = remainingSec > 0 && remainingSec <= 60;   // <= 1 min

  return (
    <>
      <div
        class={`session-timer-pill ${isDanger ? 'session-timer--danger' : isWarning ? 'session-timer--warning' : ''} ${isMinimized ? 'session-timer--minimized' : ''}`}
      >
        <button
          type="button"
          class="session-timer-toggle"
          onClick={() => setIsMinimized(!isMinimized)}
          title={isMinimized ? 'Tampilkan Timer' : 'Sembunyikan Timer'}
        >
          ⏰
        </button>

        {!isMinimized && (
          <div class="session-timer-content">
            <span class="session-timer-label">Sisa Waktu:</span>
            <span class="session-timer-time">{formatRemaining(remainingSec)}</span>
          </div>
        )}
      </div>

      {showModal && (
        <BreakTimeModal
          minutesLearned={minutesLearned}
          onContinue={handleExtend15Min}
          onGoHome={() => (window.location.href = '/beranda')}
        />
      )}
    </>
  );
}
