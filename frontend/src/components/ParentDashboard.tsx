import { useState, useEffect, useRef } from 'preact/hooks';
import { parentApi, authApi, childrenApi } from '../lib/api';

// Dynamically import Chart.js to avoid loading it on other pages
let Chart: any = null;
let chartModulesLoaded = false;

async function loadChartJS() {
  if (chartModulesLoaded) return;
  const mod = await import('chart.js');
  mod.Chart.register(
    mod.LineController,
    mod.BarController,
    mod.RadarController,
    mod.LineElement,
    mod.BarElement,
    mod.PointElement,
    mod.RadialLinearScale,
    mod.CategoryScale,
    mod.LinearScale,
    mod.Filler,
    mod.Tooltip,
    mod.Legend,
  );
  Chart = mod.Chart;
  chartModulesLoaded = true;
}

// ---- Types ----
interface Child {
  id: number;
  name: string;
  age_group: string;
  avatar_id: number;
  xp_total: number;
  current_level: number;
  streak_days: number;
  daily_limit_min?: number | null;
}

interface CategoryProgress {
  category_slug: string;
  category_name: string;
  category_color: string;
  completed: number;
  total_available: number;
  score: number;
  max_possible: number;
}

interface RecentActivity {
  lesson_id: number;
  lesson_title: string;
  status: string;
  score: number;
  max_possible: number;
  completed_at: string;
}

interface ProgressSummary {
  total_lessons_completed: number;
  total_score: number;
  total_max_possible: number;
  total_time_spent_sec: number;
  category_progress: CategoryProgress[] | null;
  recent_activity: RecentActivity[] | null;
}

interface WeeklyActivity {
  date: string;
  lessons_count: number;
  time_spent_sec: number;
}

interface CategoryScore {
  category_slug: string;
  category_name: string;
  category_color: string;
  avg_score_pct: number;
}

interface DailyTime {
  date: string;
  time_spent_sec: number;
}

interface Badge {
  id: number;
  slug: string;
  name: string;
  description: string;
  icon: string;
  earned: boolean;
  earned_at?: string;
}

interface ChildDashboard {
  child: Child;
  progress: ProgressSummary | null;
  weekly_activity: WeeklyActivity[] | null;
  category_scores: CategoryScore[] | null;
  daily_time: DailyTime[] | null;
  badges: Badge[] | null;
}

// ---- Constants ----
const AVATAR_EMOJIS = ['🦁', '🐼', '🦊', '🐸', '🦄', '🐶', '🐱', '🐰'];

const CATEGORY_EMOJIS: Record<string, string> = {
  math: '🧮', science: '🔬', coding: '💻',
  language: '📚', logic: '🧩', art: '🎨', toddlers: '✨',
};

const BADGE_EMOJIS: Record<string, string> = {
  'badge-first': '🌟', 'badge-three': '📖', 'badge-five': '🎯',
  'badge-math': '🧮', 'badge-code': '💻', 'badge-science': '🔬',
  'badge-language': '📚', 'badge-logic': '🧩', 'badge-art': '🎨',
  'badge-toddler': '✨', 'badge-streak3': '🔥', 'badge-streak7': '⭐',
  'badge-streak14': '💪', 'badge-streak30': '🏅', 'badge-perfect': '💯',
  'badge-ten': '🏆', 'badge-twenty': '👑', 'badge-xp100': '⚡',
  'badge-xp500': '🌟', 'badge-level2': '📈', 'badge-level5': '🚀',
};

const AGE_GROUP_LABELS: Record<string, string> = {
  toddlers: 'Balita', explorers: 'Penjelajah',
  builders: 'Pembangun', challengers: 'Penantang',
};

function formatTime(sec: number): string {
  if (sec < 60) return `${sec} detik`;
  const min = Math.floor(sec / 60);
  if (min < 60) return `${min} menit`;
  const hr = Math.floor(min / 60);
  const remMin = min % 60;
  return `${hr}j ${remMin > 0 ? `${remMin}m` : ''}`;
}

function formatShortDate(dateStr: string): string {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
}

// ---- Component ----
export default function ParentDashboard() {
  const [pinVerified, setPinVerified] = useState(false);
  const [pinInput, setPinInput] = useState('');
  const [pinError, setPinError] = useState('');
  const [loading, setLoading] = useState(true);
  const [dashboards, setDashboards] = useState<ChildDashboard[]>([]);
  const [selectedIdx, setSelectedIdx] = useState(0);
  const [savingLimit, setSavingLimit] = useState(false);
  const [limitMsg, setLimitMsg] = useState('');
  const [difficultiesMap, setDifficultiesMap] = useState<Record<number, any[]>>({});

  const loadDifficulty = async (childId: number) => {
    const res = await childrenApi.getDifficulty(childId);
    if (res.data?.difficulties) {
      setDifficultiesMap(prev => ({ ...prev, [childId]: res.data.difficulties }));
    }
  };

  const handleAcceptDifficultyParent = async (childId: number, categoryId: number) => {
    await childrenApi.acceptDifficulty(childId, categoryId);
    loadDifficulty(childId);
  };

  const handleDismissDifficultyParent = async (childId: number, categoryId: number) => {
    await childrenApi.dismissDifficulty(childId, categoryId);
    loadDifficulty(childId);
  };

  const handleSaveDailyLimit = async (childId: number, limitMin: number | null) => {
    setSavingLimit(true);
    setLimitMsg('');
    const res = await childrenApi.updateDailyLimit(childId, limitMin);
    setSavingLimit(false);
    if (res.error) {
      setLimitMsg(res.error);
    } else {
      setLimitMsg('Batas waktu harian berhasil disimpan! ✨');
      setDashboards(prev =>
        prev.map(d => {
          if (d.child.id === childId) {
            return {
              ...d,
              child: {
                ...d.child,
                daily_limit_min: limitMin,
              },
            };
          }
          return d;
        })
      );
      setTimeout(() => setLimitMsg(''), 3000);
    }
  };

  // Chart refs
  const activityChartRef = useRef<HTMLCanvasElement>(null);
  const radarChartRef = useRef<HTMLCanvasElement>(null);
  const timeChartRef = useRef<HTMLCanvasElement>(null);
  const chartInstances = useRef<any[]>([]);

  // Check if PIN is set
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    const res = await authApi.me();
    if (res.error) {
      window.location.href = '/masuk';
      return;
    }
    // If no PIN is set, skip PIN gate
    if (!res.data?.account?.parent_pin || res.data.account.parent_pin === '') {
      setPinVerified(true);
      loadDashboard();
    }
  };

  const verifyPin = async () => {
    if (pinInput.length !== 4) {
      setPinError('PIN harus 4 digit');
      return;
    }
    const res = await authApi.verifyPin(pinInput);
    if (res.error) {
      setPinError('PIN salah. Coba lagi.');
      setPinInput('');
      return;
    }
    setPinVerified(true);
    setPinError('');
    loadDashboard();
  };

  const loadDashboard = async () => {
    setLoading(true);
    const res = await parentApi.getDashboard();
    if (res.data?.children) {
      setDashboards(res.data.children);
    }
    setLoading(false);
  };

  // Render charts & fetch difficulty when data or selected child changes
  useEffect(() => {
    if (!pinVerified || loading || dashboards.length === 0) return;

    const data = dashboards[selectedIdx];
    if (data?.child?.id) {
      loadDifficulty(data.child.id);
    }

    const renderCharts = async () => {
      await loadChartJS();
      if (!Chart) return;

      // Destroy existing chart instances
      chartInstances.current.forEach(c => c?.destroy());
      chartInstances.current = [];

      const data = dashboards[selectedIdx];
      if (!data) return;

      // ---- Activity Line Chart ----
      if (activityChartRef.current) {
        const ctx = activityChartRef.current.getContext('2d');
        if (ctx) {
          const weekly = data.weekly_activity || [];
          const labels = weekly.map(w => formatShortDate(w.date));
          const values = weekly.map(w => w.lessons_count);

          const gradient = ctx.createLinearGradient(0, 0, 0, 260);
          gradient.addColorStop(0, 'rgba(56, 189, 248, 0.3)');
          gradient.addColorStop(1, 'rgba(56, 189, 248, 0.02)');

          chartInstances.current.push(new Chart(ctx, {
            type: 'line',
            data: {
              labels,
              datasets: [{
                label: 'Pelajaran Selesai',
                data: values,
                borderColor: '#38bdf8',
                backgroundColor: gradient,
                fill: true,
                tension: 0.4,
                pointBackgroundColor: '#38bdf8',
                pointBorderColor: '#1e293b',
                pointBorderWidth: 2,
                pointRadius: 4,
                pointHoverRadius: 6,
              }],
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              plugins: {
                legend: { display: false },
                tooltip: {
                  backgroundColor: 'rgba(15, 23, 42, 0.9)',
                  titleColor: '#e2e8f0',
                  bodyColor: '#94a3b8',
                  borderColor: 'rgba(56, 189, 248, 0.3)',
                  borderWidth: 1,
                  cornerRadius: 8,
                  padding: 10,
                },
              },
              scales: {
                x: {
                  grid: { color: 'rgba(255,255,255,0.04)' },
                  ticks: { color: '#64748b', font: { size: 11 } },
                },
                y: {
                  beginAtZero: true,
                  grid: { color: 'rgba(255,255,255,0.04)' },
                  ticks: {
                    color: '#64748b',
                    font: { size: 11 },
                    stepSize: 1,
                  },
                },
              },
            },
          }));
        }
      }

      // ---- Category Radar Chart ----
      if (radarChartRef.current) {
        const ctx = radarChartRef.current.getContext('2d');
        if (ctx) {
          const scores = data.category_scores || [];
          const labels = scores.map(s => s.category_name);
          const values = scores.map(s => s.avg_score_pct);
          const colors = scores.map(s => s.category_color);

          chartInstances.current.push(new Chart(ctx, {
            type: 'radar',
            data: {
              labels,
              datasets: [{
                label: 'Skor Rata-rata (%)',
                data: values,
                backgroundColor: 'rgba(167, 139, 250, 0.2)',
                borderColor: '#a78bfa',
                borderWidth: 2,
                pointBackgroundColor: colors,
                pointBorderColor: '#1e293b',
                pointBorderWidth: 2,
                pointRadius: 5,
                pointHoverRadius: 7,
              }],
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              plugins: {
                legend: { display: false },
                tooltip: {
                  backgroundColor: 'rgba(15, 23, 42, 0.9)',
                  titleColor: '#e2e8f0',
                  bodyColor: '#94a3b8',
                  borderColor: 'rgba(167, 139, 250, 0.3)',
                  borderWidth: 1,
                  cornerRadius: 8,
                  callbacks: {
                    label: (ctx: any) => `${ctx.parsed.r.toFixed(1)}%`,
                  },
                },
              },
              scales: {
                r: {
                  beginAtZero: true,
                  max: 100,
                  ticks: {
                    color: '#64748b',
                    backdropColor: 'transparent',
                    font: { size: 10 },
                    stepSize: 25,
                  },
                  grid: { color: 'rgba(255,255,255,0.06)' },
                  pointLabels: {
                    color: '#94a3b8',
                    font: { size: 11, weight: '600' as any },
                  },
                  angleLines: { color: 'rgba(255,255,255,0.06)' },
                },
              },
            },
          }));
        }
      }

      // ---- Time Spent Bar Chart ----
      if (timeChartRef.current) {
        const ctx = timeChartRef.current.getContext('2d');
        if (ctx) {
          const daily = data.daily_time || [];
          const labels = daily.map(d => formatShortDate(d.date));
          const values = daily.map(d => Math.round(d.time_spent_sec / 60));

          const gradient = ctx.createLinearGradient(0, 0, 0, 260);
          gradient.addColorStop(0, 'rgba(16, 185, 129, 0.7)');
          gradient.addColorStop(1, 'rgba(16, 185, 129, 0.15)');

          chartInstances.current.push(new Chart(ctx, {
            type: 'bar',
            data: {
              labels,
              datasets: [{
                label: 'Menit Belajar',
                data: values,
                backgroundColor: gradient,
                borderColor: '#10b981',
                borderWidth: 1,
                borderRadius: 6,
                borderSkipped: false,
              }],
            },
            options: {
              responsive: true,
              maintainAspectRatio: false,
              plugins: {
                legend: { display: false },
                tooltip: {
                  backgroundColor: 'rgba(15, 23, 42, 0.9)',
                  titleColor: '#e2e8f0',
                  bodyColor: '#94a3b8',
                  borderColor: 'rgba(16, 185, 129, 0.3)',
                  borderWidth: 1,
                  cornerRadius: 8,
                  callbacks: {
                    label: (ctx: any) => `${ctx.parsed.y} menit`,
                  },
                },
              },
              scales: {
                x: {
                  grid: { display: false },
                  ticks: { color: '#64748b', font: { size: 11 } },
                },
                y: {
                  beginAtZero: true,
                  grid: { color: 'rgba(255,255,255,0.04)' },
                  ticks: {
                    color: '#64748b',
                    font: { size: 11 },
                    callback: (v: any) => `${v}m`,
                  },
                },
              },
            },
          }));
        }
      }
    };

    // Small delay to ensure DOM is ready
    const timeout = setTimeout(renderCharts, 100);
    return () => {
      clearTimeout(timeout);
      chartInstances.current.forEach(c => c?.destroy());
    };
  }, [pinVerified, loading, selectedIdx, dashboards]);

  // ---- PIN Gate ----
  if (!pinVerified) {
    return (
      <div class="parent-dashboard">
        <nav class="pd-nav">
          <div class="pd-nav-inner">
            <a href="/beranda" class="dash-logo" style="text-decoration:none;">
              <span class="logo-icon">🎓</span>
              <span class="logo-text">Ez<span class="text-gradient">Edu</span></span>
            </a>
            <div class="pd-nav-right">
              <a href="/beranda" class="btn-ghost">← Beranda</a>
            </div>
          </div>
        </nav>

        <div style="display:flex;flex-direction:column;align-items:center;justify-content:center;min-height:70vh;padding:var(--space-xl);">
          <div style="background:rgba(30,41,59,0.7);border:1px solid rgba(255,255,255,0.1);border-radius:var(--radius-xl);padding:var(--space-2xl);max-width:380px;width:100%;text-align:center;">
            <div style="font-size:3rem;margin-bottom:var(--space-md);">🔒</div>
            <h2 style="font-size:1.3rem;font-weight:700;color:var(--color-text-primary);margin-bottom:var(--space-xs);">
              Dashboard Orangtua
            </h2>
            <p style="color:var(--color-text-muted);font-size:0.9rem;margin-bottom:var(--space-lg);">
              Masukkan PIN orangtua untuk melihat laporan kemajuan anak.
            </p>

            <div style="display:flex;gap:var(--space-xs);justify-content:center;margin-bottom:var(--space-md);">
              {[0, 1, 2, 3].map(i => (
                <input
                  key={i}
                  type="password"
                  maxLength={1}
                  value={pinInput[i] || ''}
                  onInput={(e: any) => {
                    const val = e.target.value.replace(/\D/g, '');
                    const newPin = pinInput.split('');
                    newPin[i] = val;
                    const joined = newPin.join('').slice(0, 4);
                    setPinInput(joined);
                    if (val && i < 3) {
                      const next = e.target.parentElement.children[i + 1] as HTMLInputElement;
                      next?.focus();
                    }
                  }}
                  onKeyDown={(e: any) => {
                    if (e.key === 'Backspace' && !pinInput[i] && i > 0) {
                      const prev = e.target.parentElement.children[i - 1] as HTMLInputElement;
                      prev?.focus();
                    }
                    if (e.key === 'Enter' && pinInput.length === 4) {
                      verifyPin();
                    }
                  }}
                  style="width:52px;height:58px;text-align:center;font-size:1.5rem;font-weight:700;background:rgba(15,23,42,0.6);border:1px solid rgba(255,255,255,0.15);border-radius:var(--radius-md);color:var(--color-text-primary);outline:none;"
                />
              ))}
            </div>

            {pinError && (
              <p style="color:#f87171;font-size:0.85rem;margin-bottom:var(--space-sm);">{pinError}</p>
            )}

            <button
              type="button"
              onClick={verifyPin}
              class="btn-primary"
              style="width:100%;padding:var(--space-sm) var(--space-lg);font-size:1rem;border-radius:var(--radius-md);"
            >
              Buka Dashboard 📊
            </button>
          </div>
        </div>
      </div>
    );
  }

  // ---- Loading ----
  if (loading) {
    return (
      <div class="pd-loading">
        <div class="loading-spinner" />
        <p>Memuat dashboard orangtua...</p>
      </div>
    );
  }

  // ---- No Children ----
  if (dashboards.length === 0) {
    return (
      <div class="parent-dashboard">
        <nav class="pd-nav">
          <div class="pd-nav-inner">
            <a href="/beranda" class="dash-logo" style="text-decoration:none;">
              <span class="logo-icon">🎓</span>
              <span class="logo-text">Ez<span class="text-gradient">Edu</span></span>
            </a>
          </div>
        </nav>
        <div class="pd-empty" style="padding-top:15vh;">
          <p style="font-size:2rem;margin-bottom:var(--space-md);">👶</p>
          <h2 style="color:var(--color-text-primary);margin-bottom:var(--space-sm);">Belum Ada Profil Anak</h2>
          <p>Buat profil anak terlebih dahulu untuk melihat dashboard kemajuan.</p>
          <a href="/profil" class="btn-primary" style="margin-top:var(--space-lg);display:inline-block;">
            Buat Profil Anak
          </a>
        </div>
      </div>
    );
  }

  // ---- Selected child data ----
  const current = dashboards[selectedIdx];
  const child = current.child;
  const progress = current.progress;
  const badges = current.badges || [];
  const earnedBadges = badges.filter((b: Badge) => b.earned);

  const completionPct = progress?.total_max_possible && progress.total_max_possible > 0
    ? Math.round((progress.total_score / progress.total_max_possible) * 100)
    : 0;

  return (
    <div class="parent-dashboard">
      {/* Nav */}
      <nav class="pd-nav">
        <div class="pd-nav-inner">
          <a href="/beranda" class="dash-logo" style="text-decoration:none;">
            <span class="logo-icon">🎓</span>
            <span class="logo-text">Ez<span class="text-gradient">Edu</span></span>
          </a>
          <div class="pd-nav-right">
            <a href="/beranda" class="btn-ghost" id="pd-back-beranda">← Beranda</a>
          </div>
        </div>
      </nav>

      {/* Hero */}
      <div class="pd-hero pd-slide-up">
        <h1>📊 Dashboard Orangtua</h1>
        <p>Pantau kemajuan belajar anak Anda secara visual dan detail.</p>
      </div>

      {/* Child Tabs */}
      <div class="pd-child-tabs pd-fade-in">
        {dashboards.map((d, i) => (
          <button
            key={d.child.id}
            type="button"
            class={`pd-child-tab ${i === selectedIdx ? 'active' : ''}`}
            onClick={() => setSelectedIdx(i)}
          >
            <span class="pd-child-tab-avatar">
              {AVATAR_EMOJIS[d.child.avatar_id - 1] || '🦁'}
            </span>
            <span>{d.child.name}</span>
            <span style="font-size:0.75rem;opacity:0.7;">
              ({AGE_GROUP_LABELS[d.child.age_group] || d.child.age_group})
            </span>
          </button>
        ))}
      </div>

      {/* Dashboard Content */}
      <div class="pd-content pd-fade-in" key={selectedIdx}>
        {/* Stat Cards */}
        <div class="pd-stats-grid">
          <div class="pd-stat-card">
            <span class="pd-stat-icon">⭐</span>
            <span class="pd-stat-value">{child.xp_total}</span>
            <span class="pd-stat-label">Total XP</span>
          </div>
          <div class="pd-stat-card">
            <span class="pd-stat-icon">🏆</span>
            <span class="pd-stat-value">Lv.{child.current_level}</span>
            <span class="pd-stat-label">Level</span>
          </div>
          <div class="pd-stat-card">
            <span class="pd-stat-icon">🔥</span>
            <span class="pd-stat-value">{child.streak_days}</span>
            <span class="pd-stat-label">Hari Streak</span>
          </div>
          <div class="pd-stat-card">
            <span class="pd-stat-icon">📚</span>
            <span class="pd-stat-value">{progress?.total_lessons_completed || 0}</span>
            <span class="pd-stat-label">Pelajaran Selesai</span>
          </div>
          <div class="pd-stat-card">
            <span class="pd-stat-icon">💯</span>
            <span class="pd-stat-value">{completionPct}%</span>
            <span class="pd-stat-label">Akurasi</span>
          </div>
          <div class="pd-stat-card">
            <span class="pd-stat-icon">⏱️</span>
            <span class="pd-stat-value">{formatTime(progress?.total_time_spent_sec || 0)}</span>
            <span class="pd-stat-label">Waktu Belajar</span>
          </div>
        </div>

        {/* Daily Time Limit Settings */}
        <div class="pd-time-limit-card">
          <div class="pd-tl-header">
            <div>
              <h3>⏰ Batas Waktu Harian ({child.name})</h3>
              <p>Atur durasi maksimal belajar per hari. Timer akan muncul saat anak belajar.</p>
            </div>
            <span class={`pd-tl-badge ${child.daily_limit_min ? 'active' : ''}`}>
              {child.daily_limit_min ? `${child.daily_limit_min} Min / Hari` : 'Tanpa Batas'}
            </span>
          </div>

          <div class="pd-tl-presets">
            {[15, 30, 45, 60, 90, 120].map(mins => (
              <button
                key={mins}
                type="button"
                class={`pd-tl-preset-btn ${child.daily_limit_min === mins ? 'active' : ''}`}
                onClick={() => handleSaveDailyLimit(child.id, mins)}
                disabled={savingLimit}
              >
                {mins} Menit
              </button>
            ))}
            <button
              type="button"
              class={`pd-tl-preset-btn ${!child.daily_limit_min ? 'active' : ''}`}
              onClick={() => handleSaveDailyLimit(child.id, null)}
              disabled={savingLimit}
            >
              Tanpa Batas
            </button>
          </div>

          {limitMsg && <p class="pd-tl-msg">{limitMsg}</p>}
        </div>

        {/* Adaptive Difficulty Section */}
        <div class="pd-time-limit-card">
          <div class="pd-tl-header">
            <div>
              <h3>🧠 Tingkat Kesulitan Adaptif ({child.name})</h3>
              <p>Sistem otomatis mengevaluasi hasil kuis dan merekomendasikan penyesuaian level jika anak kesulitan atau sangat mahir.</p>
            </div>
          </div>

          <div class="pd-diff-grid">
            {(difficultiesMap[child.id] || []).map((diff: any) => (
              <div class={`pd-diff-card ${diff.recommendation ? 'has-rec' : ''}`} key={diff.category_id}>
                <div class="pd-diff-header">
                  <span class="pd-diff-cat">{diff.category_name || diff.category_slug}</span>
                  <span class="pd-diff-level">Level {diff.current_level}</span>
                </div>

                {diff.recommendation ? (
                  <div class="pd-diff-rec-box">
                    <span class="pd-diff-rec-text">
                      {diff.recommendation === 'easier'
                        ? `📉 Disarankan turun ke Level ${diff.recommended_level}`
                        : `🚀 Disarankan naik ke Level ${diff.recommended_level}`}
                    </span>
                    <div class="pd-diff-rec-btns">
                      <button
                        type="button"
                        class="pd-tl-preset-btn active"
                        onClick={() => handleAcceptDifficultyParent(child.id, diff.category_id)}
                      >
                        Terapkan
                      </button>
                      <button
                        type="button"
                        class="pd-tl-preset-btn"
                        onClick={() => handleDismissDifficultyParent(child.id, diff.category_id)}
                      >
                        Abaikan
                      </button>
                    </div>
                  </div>
                ) : (
                  <div class="pd-diff-status-ok">
                    <span>Stabil ✅</span>
                    <span class="pd-diff-sub">
                      {diff.consecutive_high > 0 && `${diff.consecutive_high}/3 skor tinggi`}
                      {diff.consecutive_low > 0 && `${diff.consecutive_low}/3 butuh bantuan`}
                      {diff.consecutive_high === 0 && diff.consecutive_low === 0 && 'Performa seimbang'}
                    </span>
                  </div>
                )}
              </div>
            ))}
            {(!difficultiesMap[child.id] || difficultiesMap[child.id].length === 0) && (
              <p class="pd-empty" style="grid-column: 1 / -1;">
                {child.age_group === 'toddlers'
                  ? 'Kelompok usia Toddlers fokus pada eksplorasi bebas tanpa penyesuaian kesulitan.'
                  : 'Belum ada data evaluasi kesulitan. Mulai kerjakan kuis untuk mengaktifkan mesin adaptif!'}
              </p>
            )}
          </div>
        </div>

        {/* Charts */}
        <div class="pd-charts-grid">
          {/* Activity Line Chart */}
          <div class="pd-chart-card pd-chart-full">
            <h3>📈 Aktivitas Belajar (4 Minggu Terakhir)</h3>
            {(current.weekly_activity && current.weekly_activity.length > 0) ? (
              <div class="pd-chart-wrap">
                <canvas ref={activityChartRef} />
              </div>
            ) : (
              <div class="pd-chart-empty">
                Belum ada data aktivitas. Mulai belajar untuk melihat grafik!
              </div>
            )}
          </div>

          {/* Category Radar Chart */}
          <div class="pd-chart-card">
            <h3>🎯 Skor per Kategori</h3>
            {(current.category_scores && current.category_scores.length > 0) ? (
              <div class="pd-chart-wrap">
                <canvas ref={radarChartRef} />
              </div>
            ) : (
              <div class="pd-chart-empty">
                Belum ada data skor kategori.
              </div>
            )}
          </div>

          {/* Time Spent Bar Chart */}
          <div class="pd-chart-card">
            <h3>⏱️ Waktu Belajar Harian (14 Hari)</h3>
            {(current.daily_time && current.daily_time.length > 0) ? (
              <div class="pd-chart-wrap">
                <canvas ref={timeChartRef} />
              </div>
            ) : (
              <div class="pd-chart-empty">
                Belum ada data waktu belajar.
              </div>
            )}
          </div>
        </div>

        {/* Category Progress Bars */}
        <h2 class="pd-section-title">📊 Kemajuan per Kategori</h2>
        <div class="pd-cat-list">
          {(progress?.category_progress || []).map(cat => {
            const pct = cat.total_available > 0
              ? Math.round((cat.completed / cat.total_available) * 100)
              : 0;
            return (
              <div class="pd-cat-row" key={cat.category_slug}>
                <div class="pd-cat-header">
                  <div class="pd-cat-left">
                    <span class="pd-cat-emoji">{CATEGORY_EMOJIS[cat.category_slug] || '📖'}</span>
                    <div>
                      <span class="pd-cat-name">{cat.category_name}</span>
                      <span class="pd-cat-sub">
                        {cat.completed}/{cat.total_available} selesai
                        {cat.max_possible > 0 && ` • Skor: ${cat.score}/${cat.max_possible}`}
                      </span>
                    </div>
                  </div>
                  <span class="pd-cat-pct">{pct}%</span>
                </div>
                <div class="pd-cat-bar-bg">
                  <div
                    class="pd-cat-bar-fill"
                    style={`width:${pct}%;background:linear-gradient(90deg,${cat.category_color},#38bdf8);`}
                  />
                </div>
              </div>
            );
          })}
          {(!progress?.category_progress || progress.category_progress.length === 0) && (
            <p class="pd-empty">Belum ada data kemajuan.</p>
          )}
        </div>

        {/* Badge Showcase */}
        <h2 class="pd-section-title">🏅 Lencana ({earnedBadges.length}/{badges.length})</h2>
        <div class="pd-badge-grid">
          {badges.map((badge: Badge) => (
            <div
              class={`pd-badge-card ${badge.earned ? 'pd-badge-earned' : 'pd-badge-locked'}`}
              key={badge.slug}
            >
              <span class="pd-badge-icon">
                {badge.earned ? (BADGE_EMOJIS[badge.icon] || '🏅') : '🔒'}
              </span>
              <div class="pd-badge-name">{badge.name}</div>
              <div class="pd-badge-desc">{badge.description}</div>
              {badge.earned && badge.earned_at && (
                <div style="font-size:0.68rem;color:var(--color-text-muted);margin-top:4px;">
                  {new Date(badge.earned_at).toLocaleDateString('id-ID')}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Recent Activity */}
        <h2 class="pd-section-title">📋 Aktivitas Terbaru</h2>
        <div class="pd-activity-list">
          {(progress?.recent_activity || []).map((act, i) => (
            <a
              href={`/pelajaran/${act.lesson_id}`}
              class="pd-activity-item"
              key={i}
            >
              <div class="pd-activity-left">
                <span class="pd-activity-status">
                  {act.status === 'completed' ? '✅' : '🔄'}
                </span>
                <div>
                  <span class="pd-activity-title">{act.lesson_title}</span>
                  {act.completed_at && (
                    <span class="pd-activity-date">{formatShortDate(act.completed_at)}</span>
                  )}
                </div>
              </div>
              <span class="pd-activity-score">{act.score}/{act.max_possible}</span>
            </a>
          ))}
          {(!progress?.recent_activity || progress.recent_activity.length === 0) && (
            <p class="pd-empty">Belum ada aktivitas terbaru.</p>
          )}
        </div>
      </div>
    </div>
  );
}
