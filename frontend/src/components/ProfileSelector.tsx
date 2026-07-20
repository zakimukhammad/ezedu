import { useState, useEffect } from 'preact/hooks';
import { childrenApi } from '../lib/api';

interface Child {
  id: number;
  name: string;
  birth_year: number;
  age_group: string;
  avatar_id: number;
  xp_total: number;
  current_level: number;
  streak_days: number;
}

const AGE_GROUP_LABELS: Record<string, string> = {
  explorers: 'Penjelajah (4–6)',
  builders: 'Pembangun (7–9)',
  challengers: 'Penantang (10–12)',
};

const AVATAR_EMOJIS = ['🦁', '🐼', '🦊', '🐸', '🦄', '🐶', '🐱', '🐰'];

export default function ProfileSelector() {
  const [children, setChildren] = useState<Child[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [birthYear, setBirthYear] = useState(2018);
  const [avatarId, setAvatarId] = useState(1);
  const [error, setError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    loadChildren();
  }, []);

  const loadChildren = async () => {
    setLoading(true);
    const { data, error: apiError } = await childrenApi.list();
    setLoading(false);
    if (data?.children) {
      setChildren(data.children);
    }
  };

  const handleCreate = async (e: Event) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Nama anak wajib diisi');
      return;
    }

    setFormLoading(true);
    const { data, error: apiError } = await childrenApi.create(name, birthYear, avatarId);
    setFormLoading(false);

    if (apiError) {
      setError(apiError);
      return;
    }

    setName('');
    setShowForm(false);
    loadChildren();
  };

  const selectChild = (child: Child) => {
    // Store selected child in sessionStorage
    sessionStorage.setItem('ezedu_child', JSON.stringify(child));
    window.location.href = '/beranda';
  };

  if (loading) {
    return (
      <div class="profile-loading flex-center">
        <div class="loading-spinner"></div>
        <p>Memuat profil...</p>
      </div>
    );
  }

  return (
    <div class="profile-selector">
      <h2 class="profile-title">Siapa yang mau belajar hari ini? 🎉</h2>

      <div class="profile-grid stagger-children">
        {children.map((child) => (
          <button
            key={child.id}
            class="profile-card"
            onClick={() => selectChild(child)}
            id={`profile-${child.id}`}
          >
            <div class="profile-avatar">
              {AVATAR_EMOJIS[child.avatar_id - 1] || '🦁'}
            </div>
            <span class="profile-name">{child.name}</span>
            <span class="profile-group badge badge-primary">
              {AGE_GROUP_LABELS[child.age_group] || child.age_group}
            </span>
            <div class="profile-stats">
              <span class="profile-xp">⭐ {child.xp_total} XP</span>
              <span class="profile-level">Lv. {child.current_level}</span>
            </div>
            {child.streak_days > 0 && (
              <span class="profile-streak">🔥 {child.streak_days} hari berturut-turut</span>
            )}
          </button>
        ))}

        {children.length < 4 && (
          <button
            class="profile-card profile-add"
            onClick={() => setShowForm(true)}
            id="add-profile-btn"
          >
            <div class="profile-avatar add-avatar">+</div>
            <span class="profile-name">Tambah Profil</span>
          </button>
        )}
      </div>

      {/* Add Child Form Modal */}
      {showForm && (
        <div class="modal-overlay" onClick={() => setShowForm(false)}>
          <div class="modal-content animate-scale-in" onClick={(e) => e.stopPropagation()}>
            <h3>Tambah Profil Anak</h3>
            <form onSubmit={handleCreate}>
              <div class="form-group">
                <label class="form-label" for="child-name">Nama Anak</label>
                <input
                  id="child-name"
                  class="form-input"
                  type="text"
                  placeholder="Masukkan nama anak"
                  value={name}
                  onInput={(e) => setName((e.target as HTMLInputElement).value)}
                />
              </div>

              <div class="form-group">
                <label class="form-label" for="child-birth">Tahun Lahir</label>
                <select
                  id="child-birth"
                  class="form-input"
                  value={birthYear}
                  onChange={(e) => setBirthYear(Number((e.target as HTMLSelectElement).value))}
                >
                  {Array.from({ length: 13 }, (_, i) => 2014 + i).reverse().map(year => (
                    <option key={year} value={year}>{year}</option>
                  ))}
                </select>
              </div>

              <div class="form-group">
                <label class="form-label">Pilih Avatar</label>
                <div class="avatar-grid">
                  {AVATAR_EMOJIS.map((emoji, i) => (
                    <button
                      key={i}
                      type="button"
                      class={`avatar-option ${avatarId === i + 1 ? 'avatar-selected' : ''}`}
                      onClick={() => setAvatarId(i + 1)}
                    >
                      {emoji}
                    </button>
                  ))}
                </div>
              </div>

              {error && <p class="form-error" role="alert">⚠️ {error}</p>}

              <div class="modal-actions">
                <button type="button" class="btn btn-secondary" onClick={() => setShowForm(false)}>
                  Batal
                </button>
                <button type="submit" class="btn btn-primary" disabled={formLoading} id="create-child-submit">
                  {formLoading ? 'Menyimpan...' : 'Simpan Profil'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
