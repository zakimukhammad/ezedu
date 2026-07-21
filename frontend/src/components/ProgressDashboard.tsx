import { useState, useEffect } from 'preact/hooks';
import { progressApi, badgesApi } from '../lib/api';

interface Child {
  id: number;
  name: string;
  age_group: string;
  avatar_id: number;
  xp_total: number;
  current_level: number;
  streak_days: number;
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

interface Badge {
  id: number;
  slug: string;
  name: string;
  description: string;
  icon: string;
  earned: boolean;
  earned_at?: string;
}

const AVATAR_EMOJIS = ['🦁', '🐼', '🦊', '🐸', '🦄', '🐶', '🐱', '🐰'];

const CATEGORY_EMOJIS: Record<string, string> = {
  math: '🧮',
  science: '🔬',
  coding: '💻',
  language: '📚',
  logic: '🧩',
  art: '🎨',
};

const BADGE_EMOJIS: Record<string, string> = {
  'badge-first': '🌟',
  'badge-math': '🧮',
  'badge-code': '💻',
  'badge-streak3': '🔥',
  'badge-streak7': '⭐',
  'badge-perfect': '💯',
  'badge-ten': '🏆',
};

function formatTime(sec: number): string {
  if (sec < 60) return `${sec} detik`;
  const min = Math.floor(sec / 60);
  if (min < 60) return `${min} menit`;
  const hr = Math.floor(min / 60);
  const remMin = min % 60;
  return `${hr} jam ${remMin > 0 ? `${remMin} menit` : ''}`;
}

export default function ProgressDashboard() {
  const [child, setChild] = useState<Child | null>(null);
  const [progress, setProgress] = useState<ProgressSummary | null>(null);
  const [badges, setBadges] = useState<Badge[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'overview' | 'badges'>('overview');

  useEffect(() => {
    const stored = sessionStorage.getItem('ezedu_child');
    if (!stored) {
      window.location.href = '/profil';
      return;
    }
    const childData = JSON.parse(stored);
    setChild(childData);
    loadData(childData.id);
  }, []);

  const loadData = async (childId: number) => {
    setLoading(true);
    const [progRes, badgeRes] = await Promise.all([
      progressApi.getChildProgress(childId),
      badgesApi.getChildBadges(childId),
    ]);

    if (progRes.data?.progress) {
      setProgress(progRes.data.progress);
    }
    if (progRes.data?.child) {
      setChild(progRes.data.child);
      // Update sessionStorage with fresh data
      sessionStorage.setItem('ezedu_child', JSON.stringify(progRes.data.child));
    }
    if (badgeRes.data?.badges) {
      setBadges(badgeRes.data.badges);
    }
    setLoading(false);
  };

  if (!child || loading) {
    return (
      <div class="progress-loading flex-center">
        <div class="loading-spinner"></div>
        <p>Memuat data kemajuan...</p>
      </div>
    );
  }

  const earnedCount = badges.filter(b => b.earned).length;
  const completionPct = progress?.total_max_possible && progress.total_max_possible > 0
    ? Math.round((progress.total_score / progress.total_max_possible) * 100)
    : 0;

  return (
    <div class="progress-page">
      {/* Top Nav */}
      <nav class="dash-nav">
        <div class="dash-nav-inner">
          <a href="/beranda" class="dash-logo" style="text-decoration:none;">
            <span class="logo-icon">🎓</span>
            <span class="logo-text">Ez<span class="text-gradient">Edu</span></span>
          </a>
          <div class="dash-nav-right">
            <a href="/beranda" class="btn-ghost" id="back-to-dashboard">← Beranda</a>
          </div>
        </div>
      </nav>

      {/* Hero */}
      <section class="progress-hero animate-slide-up">
        <div class="progress-hero-inner">
          <div class="progress-hero-avatar">
            {AVATAR_EMOJIS[child.avatar_id - 1] || '🦁'}
          </div>
          <div class="progress-hero-info">
            <h1 class="progress-hero-name">Kemajuan {child.name}</h1>
            <p class="progress-hero-sub">Terus semangat belajar! 🚀</p>
          </div>
        </div>

        {/* Stat cards */}
        <div class="progress-stat-grid">
          <div class="progress-stat-card stat-xp">
            <span class="psc-icon">⭐</span>
            <span class="psc-value">{child.xp_total}</span>
            <span class="psc-label">Total XP</span>
          </div>
          <div class="progress-stat-card stat-level">
            <span class="psc-icon">🏆</span>
            <span class="psc-value">Lv.{child.current_level}</span>
            <span class="psc-label">Level</span>
          </div>
          <div class="progress-stat-card stat-streak">
            <span class="psc-icon">🔥</span>
            <span class="psc-value">{child.streak_days}</span>
            <span class="psc-label">Hari Streak</span>
          </div>
          <div class="progress-stat-card stat-lessons">
            <span class="psc-icon">📚</span>
            <span class="psc-value">{progress?.total_lessons_completed || 0}</span>
            <span class="psc-label">Pelajaran Selesai</span>
          </div>
          <div class="progress-stat-card stat-accuracy">
            <span class="psc-icon">💯</span>
            <span class="psc-value">{completionPct}%</span>
            <span class="psc-label">Akurasi</span>
          </div>
          <div class="progress-stat-card stat-time">
            <span class="psc-icon">⏱️</span>
            <span class="psc-value">{formatTime(progress?.total_time_spent_sec || 0)}</span>
            <span class="psc-label">Waktu Belajar</span>
          </div>
        </div>
      </section>

      {/* Tabs */}
      <div class="progress-tabs-container">
        <div class="progress-tabs">
          <button
            type="button"
            class={`progress-tab ${activeTab === 'overview' ? 'active' : ''}`}
            onClick={() => setActiveTab('overview')}
          >
            📊 Ringkasan
          </button>
          <button
            type="button"
            class={`progress-tab ${activeTab === 'badges' ? 'active' : ''}`}
            onClick={() => setActiveTab('badges')}
          >
            🏅 Lencana ({earnedCount}/{badges.length})
          </button>
        </div>
      </div>

      {/* Tab Content */}
      <section class="progress-content">
        {activeTab === 'overview' && (
          <div class="progress-overview animate-fade-in">
            {/* Category Progress */}
            <h2 class="progress-section-title">Kemajuan per Kategori</h2>
            <div class="category-progress-list">
              {(progress?.category_progress || []).map((cat) => {
                const pct = cat.total_available > 0
                  ? Math.round((cat.completed / cat.total_available) * 100)
                  : 0;
                return (
                  <div class="cat-progress-row" key={cat.category_slug}>
                    <div class="cat-progress-header">
                      <span class="cat-progress-emoji">{CATEGORY_EMOJIS[cat.category_slug] || '📖'}</span>
                      <div class="cat-progress-meta">
                        <span class="cat-progress-name">{cat.category_name}</span>
                        <span class="cat-progress-sub">{cat.completed} dari {cat.total_available} pelajaran selesai</span>
                      </div>
                      <span class="cat-progress-pct-badge">{pct}%</span>
                    </div>
                    <div class="cat-progress-bar-bg">
                      <div
                        class="cat-progress-bar-fill"
                        style={`width: ${pct}%; background: linear-gradient(90deg, ${cat.category_color}, #38bdf8);`}
                      />
                    </div>
                    {cat.max_possible > 0 && (
                      <div class="cat-progress-footer">
                        <span class="cat-progress-score">⭐ Skor: {cat.score}/{cat.max_possible} ({Math.round((cat.score / cat.max_possible) * 100)}%)</span>
                      </div>
                    )}
                  </div>
                );
              })}
              {(!progress?.category_progress || progress.category_progress.length === 0) && (
                <p class="progress-empty">Belum ada data kemajuan. Mulai belajar untuk melihat kemajuanmu!</p>
              )}
            </div>

            {/* Recent Activity */}
            <h2 class="progress-section-title" style="margin-top: var(--space-2xl);">Aktivitas Terbaru</h2>
            <div class="recent-activity-list">
              {(progress?.recent_activity || []).map((act, i) => (
                <a
                  href={`/pelajaran/${act.lesson_id}`}
                  class="recent-activity-item"
                  key={i}
                >
                  <div class="ra-left">
                    <span class={`ra-status ${act.status === 'completed' ? 'ra-done' : 'ra-progress'}`}>
                      {act.status === 'completed' ? '✅' : '🔄'}
                    </span>
                    <span class="ra-title">{act.lesson_title}</span>
                  </div>
                  <div class="ra-right">
                    <span class="ra-score">{act.score}/{act.max_possible}</span>
                  </div>
                </a>
              ))}
              {(!progress?.recent_activity || progress.recent_activity.length === 0) && (
                <p class="progress-empty">Belum ada aktivitas. Coba selesaikan pelajaran pertamamu!</p>
              )}
            </div>
          </div>
        )}

        {activeTab === 'badges' && (
          <div class="badge-gallery animate-fade-in">
            <p class="badge-gallery-subtitle">
              Kumpulkan lencana dengan belajar, menyelesaikan pelajaran, dan menjaga streak!
            </p>
            <div class="badge-grid">
              {badges.map((badge) => (
                <div
                  class={`badge-card ${badge.earned ? 'badge-earned' : 'badge-locked'}`}
                  key={badge.slug}
                >
                  <div class="badge-icon-circle">
                    {badge.earned
                      ? (BADGE_EMOJIS[badge.icon] || '🏅')
                      : '🔒'
                    }
                  </div>
                  <h3 class="badge-name">{badge.name}</h3>
                  <p class="badge-desc">{badge.description}</p>
                  {badge.earned && badge.earned_at && (
                    <span class="badge-date">
                      Diperoleh: {new Date(badge.earned_at).toLocaleDateString('id-ID')}
                    </span>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </section>
    </div>
  );
}
