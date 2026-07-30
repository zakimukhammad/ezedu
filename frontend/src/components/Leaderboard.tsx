import { useState, useEffect } from 'preact/hooks';
import { leaderboardApi, childrenApi } from '../lib/api';
import ParentPinModal from './ParentPinModal';

interface LeaderboardEntry {
  rank: number;
  display_name: string;
  avatar_id: number;
  weekly_xp: number;
  is_me: boolean;
}

interface Child {
  id: number;
  name: string;
  age_group: string;
  avatar_id: number;
  leaderboard_opt_in?: boolean;
}

const AVATAR_EMOJIS: Record<number, string> = {
  1: '🦊',
  2: '🐱',
  3: '🐶',
  4: '🦁',
  5: '🐼',
  6: '🦄',
  7: '🚀',
  8: '🤖',
};

export default function Leaderboard() {
  const [child, setChild] = useState<Child | null>(null);
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [optIn, setOptIn] = useState<boolean>(false);
  const [childRank, setChildRank] = useState<number>(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Pin verification state
  const [showPinModal, setShowPinModal] = useState(false);
  const [pendingOptIn, setPendingOptIn] = useState<boolean | null>(null);
  const [toggling, setToggling] = useState(false);

  useEffect(() => {
    const stored = sessionStorage.getItem('ezedu_child');
    if (stored) {
      const c = JSON.parse(stored);
      setChild(c);
      loadLeaderboard(c.id);
    } else {
      setLoading(false);
    }
  }, []);

  const loadLeaderboard = async (childId: number) => {
    setLoading(true);
    setError(null);
    const { data, error: err } = await leaderboardApi.getWeekly(childId);
    setLoading(false);

    if (err) {
      setError(err);
      return;
    }

    if (data) {
      setEntries(data.entries || []);
      setOptIn(data.opt_in);
      setChildRank(data.child_rank);
    }
  };

  const handleRequestToggleOptIn = (targetOptIn: boolean) => {
    setPendingOptIn(targetOptIn);
    setShowPinModal(true);
  };

  const handlePinSuccess = async () => {
    setShowPinModal(false);
    if (pendingOptIn === null || !child) return;

    setToggling(true);
    const targetVal = pendingOptIn;
    setPendingOptIn(null);

    const { data, error: err } = await childrenApi.toggleLeaderboardOptIn(child.id, targetVal);
    setToggling(false);

    if (err) {
      alert(err);
      return;
    }

    if (data) {
      setOptIn(data.opt_in);
      // Update session storage child
      const updated = { ...child, leaderboard_opt_in: data.opt_in };
      setChild(updated);
      sessionStorage.setItem('ezedu_child', JSON.stringify(updated));
      loadLeaderboard(child.id);
    }
  };

  if (loading) {
    return (
      <div class="card p-xl text-center">
        <div class="spinner" style="margin: 0 auto 1rem;"></div>
        <p class="text-muted">Memuat papan peringkat...</p>
      </div>
    );
  }

  if (!child) {
    return (
      <div class="card p-xl text-center">
        <p class="text-muted">Pilih profil anak terlebih dahulu untuk melihat papan peringkat.</p>
        <a href="/profil" class="btn btn-primary mt-md">Pilih Profil</a>
      </div>
    );
  }

  if (child.age_group !== 'challengers') {
    return (
      <div class="card p-xl text-center">
        <div style="font-size: 3rem; margin-bottom: 0.5rem;">🏆</div>
        <h3>Papan Peringkat Kategori Challengers</h3>
        <p class="text-muted mt-sm">
          Papan peringkat mingguan khusus dirancang untuk kelompok usia <strong>Challengers (10–12 Tahun)</strong>.
        </p>
        <a href="/beranda" class="btn btn-secondary mt-lg">Kembali ke Beranda</a>
      </div>
    );
  }

  return (
    <div class="leaderboard-container">
      {/* Header Banner */}
      <div class="card leaderboard-header p-lg mb-lg">
        <div style="display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 1rem;">
          <div>
            <h2 style="margin: 0; display: flex; align-items: center; gap: 0.5rem;">
              <span>🏆</span> Papan Peringkat Mingguan
            </h2>
            <p class="text-muted mt-xs mb-0">
              Tunjukkan kemampuanmu secara anonim & raih posisi puncak tiap minggu!
            </p>
          </div>

          <div style="display: flex; align-items: center; gap: 0.75rem;">
            {optIn && childRank > 0 && (
              <div class="rank-badge">
                Peringkatmu: <strong>#{childRank}</strong>
              </div>
            )}
            <button
              type="button"
              onClick={() => handleRequestToggleOptIn(!optIn)}
              class={optIn ? 'btn-ghost text-muted' : 'btn btn-primary'}
              disabled={toggling}
              style="font-size: 0.85rem;"
            >
              {optIn ? '⚙️ Kelola Opt-in' : '✨ Bergabung Sekarang'}
            </button>
          </div>
        </div>
      </div>

      {/* Opt-In Prompt if Not Opted In */}
      {!optIn && (
        <div class="card p-xl text-center mb-lg opt-in-card" style="border: 2px dashed var(--color-primary-light, #6366f1);">
          <div style="font-size: 2.5rem; margin-bottom: 0.5rem;">🔒</div>
          <h3>Privasi & Keamanan Papan Peringkat</h3>
          <p class="text-muted mt-xs max-w-lg mx-auto" style="max-width: 500px; margin-left: auto; margin-right: auto;">
            Papan peringkat EzEdu menggunakan <strong>nama hewan anonim</strong> (seperti <em>Singa Pemberani</em>) untuk melindungi identitas anak.
            Izin orang tua (PIN) diperlukan untuk mengaktifkan fitur ini.
          </p>
          <button
            type="button"
            onClick={() => handleRequestToggleOptIn(true)}
            class="btn btn-primary btn-lg mt-md"
          >
            Aktifkan Papan Peringkat (PIN Orang Tua) 🚀
          </button>
        </div>
      )}

      {/* Leaderboard Table */}
      {optIn && (
        <div class="card p-0 overflow-hidden">
          {entries.length === 0 ? (
            <div class="p-xl text-center text-muted">
              <p>Belum ada aktivitas di papan peringkat minggu ini. Selesaikan pelajaran untuk mendapatkan XP pertama!</p>
            </div>
          ) : (
            <div class="leaderboard-table-wrapper">
              <table class="leaderboard-table">
                <thead>
                  <tr>
                    <th style="width: 70px; text-align: center;">Posisi</th>
                    <th>Nama Anonim</th>
                    <th style="text-align: right;">Weekly XP</th>
                  </tr>
                </thead>
                <tbody>
                  {entries.map((entry) => {
                    const isTop1 = entry.rank === 1;
                    const isTop2 = entry.rank === 2;
                    const isTop3 = entry.rank === 3;

                    return (
                      <tr key={entry.rank} class={`leaderboard-row ${entry.is_me ? 'row-is-me' : ''}`}>
                        <td style="text-align: center; font-weight: bold;">
                          {isTop1 ? '🥇' : isTop2 ? '🥈' : isTop3 ? '🥉' : `#${entry.rank}`}
                        </td>
                        <td>
                          <div style="display: flex; align-items: center; gap: 0.75rem;">
                            <span class="avatar-icon">{AVATAR_EMOJIS[entry.avatar_id] || '🦊'}</span>
                            <span class="player-name">{entry.display_name}</span>
                            {entry.is_me && <span class="badge-me">Kamu</span>}
                          </div>
                        </td>
                        <td style="text-align: right; font-weight: bold; color: var(--color-warning, #f59e0b);">
                          ⚡ {entry.weekly_xp} XP
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {/* Parent PIN Modal */}
      {showPinModal && (
        <ParentPinModal
          onSuccess={handlePinSuccess}
          onCancel={() => setShowPinModal(false)}
        />
      )}

      <style>{`
        .leaderboard-header {
          background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(168, 85, 247, 0.15));
          border-left: 4px solid var(--color-primary, #6366f1);
        }
        .rank-badge {
          background: rgba(245, 158, 11, 0.2);
          color: #f59e0b;
          border: 1px solid rgba(245, 158, 11, 0.4);
          padding: 0.35rem 0.75rem;
          border-radius: var(--radius-full, 9999px);
          font-size: 0.9rem;
        }
        .leaderboard-table-wrapper {
          overflow-x: auto;
        }
        .leaderboard-table {
          width: 100%;
          border-collapse: collapse;
          text-align: left;
        }
        .leaderboard-table th {
          background: rgba(255, 255, 255, 0.04);
          padding: 0.85rem 1.25rem;
          font-size: 0.85rem;
          text-transform: uppercase;
          letter-spacing: 0.05em;
          color: var(--color-text-muted, #94a3b8);
          border-bottom: 1px solid rgba(255, 255, 255, 0.08);
        }
        .leaderboard-table td {
          padding: 0.85rem 1.25rem;
          border-bottom: 1px solid rgba(255, 255, 255, 0.05);
        }
        .leaderboard-row:hover {
          background: rgba(255, 255, 255, 0.03);
        }
        .row-is-me {
          background: rgba(99, 102, 241, 0.12) !important;
          border-left: 3px solid #6366f1;
        }
        .avatar-icon {
          font-size: 1.5rem;
        }
        .player-name {
          font-weight: 600;
        }
        .badge-me {
          background: #6366f1;
          color: white;
          font-size: 0.7rem;
          font-weight: bold;
          padding: 0.15rem 0.45rem;
          border-radius: 4px;
        }
      `}</style>
    </div>
  );
}
