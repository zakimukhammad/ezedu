import { useState, useEffect } from 'preact/hooks';
import { childrenApi, authApi } from '../lib/api';
import ParentPinModal from './ParentPinModal';

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
  const [hasParentPin, setHasParentPin] = useState(false);

  // Modals state
  const [showForm, setShowForm] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [showPinModal, setShowPinModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [pinMode, setPinMode] = useState<'verify' | 'setup'>('verify');
  const [pendingAction, setPendingAction] = useState<'add' | 'edit' | 'delete' | null>(null);
  const [selectedChildId, setSelectedChildId] = useState<number | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  // Form state
  const [name, setName] = useState('');
  const [birthYear, setBirthYear] = useState(2018);
  const [avatarId, setAvatarId] = useState(1);
  const [error, setError] = useState('');
  const [formLoading, setFormLoading] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const showNotification = (msg: string) => {
    setToast(msg);
    setTimeout(() => {
      setToast(null);
    }, 3500);
  };

  const loadData = async () => {
    setLoading(true);
    const [childRes, meRes] = await Promise.all([
      childrenApi.list(),
      authApi.me(),
    ]);
    setLoading(false);

    if (childRes.data?.children) {
      setChildren(childRes.data.children);
    }
    if (meRes.data) {
      setHasParentPin(!!meRes.data.has_parent_pin);
    }
  };

  const handleAddClick = () => {
    setEditMode(false);
    setName('');
    setBirthYear(2018);
    setAvatarId(1);
    if (hasParentPin) {
      setPinMode('verify');
      setPendingAction('add');
      setShowPinModal(true);
    } else {
      setShowForm(true);
    }
  };

  const handleEditClick = (e: Event, child: Child) => {
    e.stopPropagation();
    setSelectedChildId(child.id);
    setName(child.name);
    setBirthYear(child.birth_year);
    setAvatarId(child.avatar_id);
    setEditMode(true);

    if (hasParentPin) {
      setPinMode('verify');
      setPendingAction('edit');
      setShowPinModal(true);
    } else {
      setShowForm(true);
    }
  };

  const handleDeleteClick = (e: Event, childId: number) => {
    e.stopPropagation();
    setSelectedChildId(childId);
    if (hasParentPin) {
      setPinMode('verify');
      setPendingAction('delete');
      setShowPinModal(true);
    } else {
      setShowDeleteModal(true);
    }
  };

  const executeDelete = async () => {
    if (!selectedChildId) return;
    setShowDeleteModal(false);
    await childrenApi.delete(selectedChildId);
    showNotification('Profil anak berhasil dihapus');
    setSelectedChildId(null);
    loadData();
  };

  const handlePinSuccess = () => {
    setShowPinModal(false);
    if (pendingAction === 'add') {
      setShowForm(true);
    } else if (pendingAction === 'edit') {
      setShowForm(true);
    } else if (pendingAction === 'delete' && selectedChildId) {
      setShowDeleteModal(true);
    } else if (pinMode === 'setup') {
      setHasParentPin(true);
      showNotification('PIN Orang Tua berhasil disimpan! 🔐');
    }
    setPendingAction(null);
  };

  const handleOpenPinSetup = () => {
    if (hasParentPin) {
      setPinMode('verify');
      setPendingAction(null);
      setShowPinModal(true);
    } else {
      setPinMode('setup');
      setShowPinModal(true);
    }
  };

  const handleSave = async (e: Event) => {
    e.preventDefault();
    setError('');

    if (!name.trim()) {
      setError('Nama anak wajib diisi');
      return;
    }

    setFormLoading(true);
    let apiError: string | undefined;

    if (editMode && selectedChildId) {
      const res = await childrenApi.update(selectedChildId, { name, birth_year: birthYear, avatar_id: avatarId });
      apiError = res.error;
    } else {
      const res = await childrenApi.create(name, birthYear, avatarId);
      apiError = res.error;
    }
    setFormLoading(false);

    if (apiError) {
      setError(apiError);
      return;
    }

    showNotification(editMode ? 'Profil anak berhasil diperbarui! ✨' : 'Profil anak baru berhasil dibuat! 🎉');
    setName('');
    setShowForm(false);
    loadData();
  };

  const selectChild = (child: Child) => {
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
      {/* Header & Parent PIN Control Bar */}
      <div class="profile-header-area">
        <h2 class="profile-title">Siapa yang mau belajar hari ini? 🎉</h2>
        <div class="parent-pin-bar">
          <button
            class={`parent-pin-btn ${hasParentPin ? 'pin-active' : 'pin-inactive'}`}
            onClick={handleOpenPinSetup}
            id="parent-pin-btn"
            title="Kelola Kunci PIN Orang Tua"
          >
            <span class="pin-icon">{hasParentPin ? '🔒' : '⚙️'}</span>
            <span class="pin-text">{hasParentPin ? 'PIN Orang Tua: Aktif' : 'Atur PIN Orang Tua'}</span>
          </button>
        </div>
      </div>

      {/* Children Profiles Grid */}
      <div class="profile-grid stagger-children">
        {children.map((child) => (
          <div key={child.id} class="profile-card-wrapper">
            <div class="profile-card-actions">
              <button
                type="button"
                class="pca-btn pca-edit"
                onClick={(e) => handleEditClick(e, child)}
                title="Edit Profil"
              >
                ✏️
              </button>
              <button
                type="button"
                class="pca-btn pca-delete"
                onClick={(e) => handleDeleteClick(e, child.id)}
                title="Hapus Profil"
              >
                🗑️
              </button>
            </div>

            <button
              type="button"
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
                <span class="profile-streak">🔥 {child.streak_days} hari streak</span>
              )}
            </button>
          </div>
        ))}

        {children.length < 4 && (
          <div class="profile-card-wrapper">
            <button
              type="button"
              class="profile-card profile-add"
              onClick={handleAddClick}
              id="add-profile-btn"
            >
              <div class="profile-avatar add-avatar">+</div>
              <span class="profile-name">Tambah Profil</span>
            </button>
          </div>
        )}
      </div>

      {/* Custom Toast Notification */}
      {toast && (
        <div class="custom-toast animate-slide-up">
          <span class="toast-icon">✨</span>
          <span class="toast-message">{toast}</span>
        </div>
      )}

      {/* Parent PIN Modal */}
      <ParentPinModal
        isOpen={showPinModal}
        mode={pinMode}
        onSuccess={() => {
          if (hasParentPin && pendingAction === null && pinMode === 'verify') {
            setPinMode('setup');
          } else {
            handlePinSuccess();
          }
        }}
        onCancel={() => {
          setShowPinModal(false);
          setPendingAction(null);
        }}
      />

      {/* Delete Confirmation Modal */}
      {showDeleteModal && (
        <div class="modal-overlay" onClick={() => setShowDeleteModal(false)}>
          <div class="modal-content animate-scale-in text-center" onClick={(e) => e.stopPropagation()}>
            <div class="delete-warning-icon">⚠️</div>
            <h3>Hapus Profil Anak?</h3>
            <p class="text-muted mt-sm" style="font-size: 0.9rem; color: var(--color-text-muted);">
              Apakah Anda yakin ingin menghapus profil ini? Semua data progres dan skor kuis akan terhapus secara permanen.
            </p>
            <div class="modal-actions mt-xl">
              <button type="button" class="btn btn-secondary" onClick={() => setShowDeleteModal(false)}>
                Batal
              </button>
              <button type="button" class="btn btn-danger" onClick={executeDelete} id="confirm-delete-btn">
                Ya, Hapus
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Add / Edit Child Form Modal */}
      {showForm && (
        <div class="modal-overlay" onClick={() => setShowForm(false)}>
          <div class="modal-content animate-scale-in" onClick={(e) => e.stopPropagation()}>
            <h3>{editMode ? 'Edit Profil Anak ✏️' : 'Tambah Profil Anak 👶'}</h3>
            <form onSubmit={handleSave}>
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
                  {formLoading ? 'Menyimpan...' : (editMode ? 'Simpan Perubahan' : 'Simpan Profil')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
