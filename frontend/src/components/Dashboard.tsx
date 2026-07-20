import { useState, useEffect } from 'preact/hooks';
import { categoriesApi, authApi } from '../lib/api';

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
  math: '🧮',
  science: '🔬',
  coding: '💻',
  language: '📚',
  logic: '🧩',
  art: '🎨',
};

const AVATAR_EMOJIS = ['🦁', '🐼', '🦊', '🐸', '🦄', '🐶', '🐱', '🐰'];

const AGE_GROUP_THEMES: Record<string, string> = {
  explorers: 'theme-explorers',
  builders: 'theme-builders',
  challengers: 'theme-challengers',
};

const AGE_GROUP_GREETING: Record<string, string> = {
  explorers: 'Ayo kita mulai belajar! 🌟',
  builders: 'Siap belajar hari ini? 🚀',
  challengers: 'Tantangan hari ini menantimu! ⚡',
};

export default function Dashboard() {
  const [child, setChild] = useState<Child | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Get selected child from sessionStorage
    const stored = sessionStorage.getItem('ezedu_child');
    if (!stored) {
      window.location.href = '/profil';
      return;
    }
    setChild(JSON.parse(stored));
    loadCategories();
  }, []);

  const loadCategories = async () => {
    const { data } = await categoriesApi.list();
    if (data?.categories) {
      setCategories(data.categories);
    }
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
