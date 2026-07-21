import { useState, useEffect } from 'preact/hooks';
import { lessonsApi } from '../lib/api';

interface Lesson {
  id: number;
  category_id: number;
  age_group: string;
  level: number;
  sort_order: number;
  title: string;
  description: string;
  estimated_minutes: number;
  xp_reward: number;
}

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
}

interface Props {
  categorySlug: string;
}

const CATEGORY_EMOJIS: Record<string, string> = {
  math: '🧮',
  science: '🔬',
  coding: '💻',
  language: '📚',
  logic: '🧩',
  art: '🎨',
};

export default function CategoryLessons({ categorySlug }: Props) {
  const [category, setCategory] = useState<Category | null>(null);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [progressMap, setProgressMap] = useState<Record<number, any>>({});
  const [child, setChild] = useState<Child | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const stored = sessionStorage.getItem('ezedu_child');
    let ageGroup = 'builders';
    let childId: number | undefined = undefined;

    if (stored) {
      const parsedChild = JSON.parse(stored);
      setChild(parsedChild);
      ageGroup = parsedChild.age_group || 'builders';
      childId = parsedChild.id;
    }

    loadLessons(categorySlug, ageGroup, childId);
  }, [categorySlug]);

  const loadLessons = async (slug: string, ageGroup: string, childId?: number) => {
    setLoading(true);
    const { data, error } = await lessonsApi.listByCategory(slug, ageGroup, childId);
    setLoading(false);

    if (data) {
      setCategory(data.category);
      setLessons(data.lessons || []);
      setProgressMap(data.progress || {});
    }
  };

  if (loading) {
    return (
      <div class="lessons-loading flex-center">
        <div class="loading-spinner"></div>
        <p>Memuat daftar pelajaran...</p>
      </div>
    );
  }

  if (child && ((child.age_group === 'toddlers' && categorySlug !== 'toddlers') || (child.age_group !== 'toddlers' && categorySlug === 'toddlers'))) {
    return (
      <div class="lessons-error text-center" style="padding: var(--space-3xl); max-width: 600px; margin: 0 auto;">
        <h2>Modul Ini Khusus Kelompok Usia Lain 🔒</h2>
        <p class="text-muted mt-md">Materi ini disesuaikan khusus untuk kelompok usia lain. Yuk pilih kategori pelajaran yang sesuai dengan {child.name}!</p>
        <a href="/beranda" class="btn btn-primary btn-lg mt-xl">Kembali ke Beranda 🚀</a>
      </div>
    );
  }

  if (!category) {
    return (
      <div class="lessons-error text-center">
        <h2>Kategori Tidak Ditemukan 🙁</h2>
        <a href="/beranda" class="btn btn-primary mt-lg">Kembali ke Beranda</a>
      </div>
    );
  }

  // Group lessons by level (Level 1 to Level 5)
  const levelMap: Record<number, Lesson[]> = {};
  for (let i = 1; i <= 5; i++) {
    levelMap[i] = [];
  }
  lessons.forEach((l) => {
    if (!levelMap[l.level]) levelMap[l.level] = [];
    levelMap[l.level].push(l);
  });

  return (
    <div class="category-lessons-view">
      {/* Category Header */}
      <header class="cat-header" style={`--cat-theme-color: ${category.color}`}>
        <div class="cat-header-inner">
          <a href="/beranda" class="back-link">← Kembali ke Beranda</a>
          <div class="cat-header-content mt-md">
            <div class="cat-header-icon">
              {CATEGORY_EMOJIS[category.slug] || '📖'}
            </div>
            <div>
              <h1 class="cat-header-title">{category.name}</h1>
              <p class="cat-header-desc">{category.description}</p>
              {child && (
                <span class="badge badge-primary mt-xs">
                  Modul: {child.age_group === 'explorers' ? 'Penjelajah (4–6 th)' : child.age_group === 'builders' ? 'Pembangun (7–9 th)' : 'Penantang (10–12 th)'}
                </span>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Levels & Lessons Timeline */}
      <main class="cat-levels-container container">
        {[1, 2, 3, 4, 5].map((levelNum) => {
          const levelLessons = levelMap[levelNum] || [];
          return (
            <div key={levelNum} class="level-section">
              <div class="level-badge-header">
                <span class="level-pill">Level {levelNum}</span>
                <div class="level-line"></div>
              </div>

              {levelLessons.length === 0 ? (
                <div class="empty-level-card card">
                  <p class="text-muted">Materi Level {levelNum} sedang dipersiapkan untuk grup usia ini 🚀</p>
                </div>
              ) : (
                <div class="lesson-cards-grid">
                  {levelLessons.map((lesson) => {
                    const prog = progressMap[lesson.id];
                    const isCompleted = prog?.status === 'completed';

                    return (
                      <div key={lesson.id} class={`lesson-card card ${isCompleted ? 'lesson-completed' : ''}`}>
                        <div class="lesson-card-header">
                          <span class="badge badge-primary">⏱️ {lesson.estimated_minutes} Min</span>
                          <span class="badge badge-warning">⭐ +{lesson.xp_reward} XP</span>
                        </div>
                        <h3 class="lesson-card-title mt-sm">{lesson.title}</h3>
                        <p class="lesson-card-desc">{lesson.description}</p>
                        
                        <div class="lesson-card-footer mt-md">
                          {isCompleted ? (
                            <div class="completed-info">
                              <span class="completed-badge">✅ Selesai ({prog.score}/{prog.max_possible || 20})</span>
                              <a href={`/pelajaran/${lesson.id}`} class="btn btn-secondary btn-sm">
                                Ulangi 🔄
                              </a>
                            </div>
                          ) : (
                            <a href={`/pelajaran/${lesson.id}`} class="btn btn-primary btn-sm w-full">
                              Mulai Pelajaran 🚀
                            </a>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          );
        })}
      </main>
    </div>
  );
}
