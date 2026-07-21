import { useState, useEffect } from 'preact/hooks';
import { categoriesApi, authApi, lessonsApi, dailyChallengeApi } from '../lib/api';

interface Category {
  id: number;
  slug: string;
  name: string;
  description: string;
  icon: string;
  color: string;
}

interface Child {
  id: number;
  name: string;
  age_group: string;
  avatar_id: number;
  xp_total: number;
  current_level: number;
  streak_days: number;
}

const CATEGORY_EMOJIS: Record<string, string> = {
  toddlers: '🎈',
  math: '🧮',
  science: '🔬',
  coding: '💻',
  language: '📚',
  logic: '🧩',
  art: '🎨',
};

const AVATAR_EMOJIS = ['🦁', '🐼', '🦊', '🐸', '🦄', '🐶', '🐱', '🐰'];

const AGE_GROUP_THEMES: Record<string, string> = {
  toddlers: 'theme-toddler',
  explorers: 'theme-explorer',
  builders: 'theme-builder',
  challengers: 'theme-challenger',
};

const AGE_GROUP_GREETING: Record<string, string> = {
  toddlers: 'Selamat datang, Adik Cilik! 🎈',
  explorers: 'Ayo kita mulai belajar! 🌟',
  builders: 'Siap belajar hari ini? 🚀',
  challengers: 'Tantangan hari ini menantimu! ⚡',
};

interface NextLesson {
  id: number;
  title: string;
  category_name: string;
  description: string;
}

export default function Dashboard() {
  const [child, setChild] = useState<Child | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [nextLesson, setNextLesson] = useState<NextLesson | null>(null);
  const [dailyChallenge, setDailyChallenge] = useState<any>(null);
  const [dailyCompleted, setDailyCompleted] = useState(false);
  const [dailyStreak, setDailyStreak] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Get selected child from sessionStorage
    const stored = sessionStorage.getItem('ezedu_child');
    if (!stored) {
      window.location.href = '/profil';
      return;
    }
    const childData = JSON.parse(stored);
    setChild(childData);
    loadData(childData);
  }, []);

  const loadData = async (childData: Child) => {
    const { data: catData } = await categoriesApi.list();
    if (catData?.categories) {
      const filtered = catData.categories.filter((cat: Category) => {
        if (childData.age_group === 'toddlers') {
          return cat.slug === 'toddlers';
        }
        return cat.slug !== 'toddlers';
      });
      setCategories(filtered);
    }

    // Load Daily Challenge
    try {
      const { data: dailyData } = await dailyChallengeApi.getToday(childData.age_group, childData.id);
      if (dailyData?.challenge) {
        setDailyChallenge(dailyData.challenge);
        setDailyCompleted(!!dailyData.completed);
      }
      const { data: streakData } = await dailyChallengeApi.getStreak(childData.id);
      if (streakData?.streak !== undefined) {
        setDailyStreak(streakData.streak);
      }
    } catch (e) {}

    // Try finding next uncompleted lesson in the appropriate age group category
    try {
      const targetSlug = childData.age_group === 'toddlers' ? 'toddlers' : 'math';
      const { data: targetData } = await lessonsApi.listByCategory(targetSlug, childData.age_group, childData.id);
      if (targetData?.lessons) {
        const progMap = targetData.progress || {};
        const uncompleted = targetData.lessons.find((l: any) => !progMap[l.id] || progMap[l.id].status !== 'completed');
        if (uncompleted) {
          setNextLesson({
            id: uncompleted.id,
            title: uncompleted.title,
            category_name: `${targetData.category?.name || 'Mengenal Dunia'} ${CATEGORY_EMOJIS[targetSlug] || '🎈'}`,
            description: uncompleted.description,
          });
        }
      }
    } catch(e) {}

    setLoading(false);
  };

  const handleLogout = async () => {
    await authApi.logout();
    sessionStorage.removeItem('ezedu_child');
    window.location.href = '/';
  };

  const switchProfile = () => {
    sessionStorage.removeItem('ezedu_child');
    window.location.href = '/profil';
  };

  if (!child || loading) {
    return (
      <div class="dashboard-loading flex-center">
        <div class="loading-spinner"></div>
        <p>Memuat beranda...</p>
      </div>
    );
  }

  const themeClass = AGE_GROUP_THEMES[child.age_group] || '';
  const greeting = AGE_GROUP_GREETING[child.age_group] || 'Selamat datang!';

  return (
    <div class={`dashboard ${themeClass}`}>
      {/* Top Nav */}
      <nav class="dash-nav">
        <div class="dash-nav-inner">
          <span class="dash-logo">🎓 EzEdu</span>
          <div class="dash-nav-right">
            <a href="/kemajuan" class="btn-ghost" id="view-progress">📊 Kemajuan</a>
            <button class="btn-ghost" onClick={switchProfile} id="switch-profile">
              Ganti Profil
            </button>
            <button class="btn-ghost" onClick={handleLogout} id="logout-btn">
              Keluar
            </button>
          </div>
        </div>
      </nav>

      {/* Welcome Section */}
      <section class="dash-welcome animate-slide-up">
        <div class="welcome-left">
          <div class="welcome-avatar">
            {AVATAR_EMOJIS[child.avatar_id - 1] || '🦁'}
          </div>
          <div>
            <h1 class="welcome-name">Halo, {child.name}! 👋</h1>
            <p class="welcome-greeting">{greeting}</p>
          </div>
        </div>
        <div class="welcome-stats">
          <div class="welcome-stat">
            <span class="welcome-stat-value">⭐ {child.xp_total}</span>
            <span class="welcome-stat-label">Total XP</span>
          </div>
          <div class="welcome-stat">
            <span class="welcome-stat-value">🏆 Lv.{child.current_level}</span>
            <span class="welcome-stat-label">Level</span>
          </div>
          {child.streak_days > 0 && (
            <div class="welcome-stat">
              <span class="welcome-stat-value">🔥 {child.streak_days}</span>
              <span class="welcome-stat-label">Hari Streak</span>
            </div>
          )}
        </div>
      </section>

      {/* Daily Challenge Card */}
      {dailyChallenge && (
        <section class="dash-daily animate-slide-up" style="max-width: 1200px; margin: var(--space-xl) auto 0; padding: 0 var(--space-lg);">
          <div class="daily-challenge-card">
            <div class="daily-challenge-header">
              <span class="daily-challenge-icon">🎯</span>
              <span class="daily-challenge-title">Tantangan Harian</span>
              {dailyStreak > 0 && (
                <span class="daily-challenge-streak">🔥 {dailyStreak} Hari Streak</span>
              )}
            </div>
            <div class="daily-challenge-body">
              <p class="daily-challenge-preview">
                {JSON.parse(dailyChallenge.question_json || '{}').prompt || 'Tantangan harian siap dimainkan!'}
              </p>
              {dailyCompleted ? (
                <div class="daily-challenge-completed">
                  <span>✅ Kamu sudah menyelesaikan tantangan hari ini! Sampai jumpa besok! 🎉</span>
                </div>
              ) : (
                <a
                  href={`/pelajaran/39`}
                  class="btn btn-primary mt-sm"
                  style="display: inline-block;"
                  id="start-daily-challenge-btn"
                >
                  Mainkan Tantangan Harian 🚀
                </a>
              )}
            </div>
          </div>
        </section>
      )}

      {/* Continue Learning Banner */}
      {nextLesson && (
        <section class="dash-continue animate-slide-up" style="max-width: 1200px; margin: var(--space-xl) auto 0; padding: 0 var(--space-lg);">
          <div class="continue-card" style="background: linear-gradient(135deg, rgba(99, 102, 241, 0.12), rgba(16, 185, 129, 0.12)); border: 2px solid var(--color-primary); border-radius: var(--radius-xl); padding: var(--space-xl); display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: var(--space-md);">
            <div>
              <span class="badge badge-primary mb-xs">🚀 Lanjutkan Belajar</span>
              <h2 style="font-size: 1.3rem; margin-top: var(--space-xs); color: var(--color-text);">{nextLesson.title}</h2>
              <p style="color: var(--color-text-muted); font-size: 0.9rem; margin-top: var(--space-xs);">{nextLesson.category_name} • {nextLesson.description}</p>
            </div>
            <a href={`/pelajaran/${nextLesson.id}`} class="btn btn-primary btn-lg" id="continue-lesson-btn">
              Mulai Belajar Sekarang →
            </a>
          </div>
        </section>
      )}

      {/* Categories Grid */}
      <section class="dash-content">
        <h2 class="dash-section-title">Pilih Kategori Belajar</h2>
        <div class="category-grid stagger-children">
          {categories.map((cat) => (
            <a
              key={cat.id}
              href={`/belajar/${cat.slug}`}
              class="dash-category-card"
              style={`--cat-color: ${cat.color}`}
              id={`cat-${cat.slug}`}
            >
              <div class="dash-cat-icon">
                {CATEGORY_EMOJIS[cat.slug] || '📖'}
              </div>
              <h3 class="dash-cat-name">{cat.name}</h3>
              <p class="dash-cat-desc">{cat.description}</p>
              <span class="dash-cat-cta">
                Mulai Belajar →
              </span>
            </a>
          ))}
        </div>
      </section>

      {/* Daily Tip */}
      <section class="dash-tip animate-fade-in">
        <div class="tip-card">
          <span class="tip-icon">💡</span>
          <div>
            <strong>Tips Hari Ini:</strong>
            <p>Belajar 10 menit setiap hari lebih efektif daripada belajar 1 jam seminggu sekali!</p>
          </div>
        </div>
      </section>
    </div>
  );
}
