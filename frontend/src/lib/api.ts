// EzEdu — API Client
// Handles all communication with the Go backend

const API_BASE = '/api';

interface ApiResponse<T = any> {
  data?: T;
  error?: string;
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    const res = await fetch(`${API_BASE}${endpoint}`, {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    const data = await res.json();

    if (!res.ok) {
      return { error: data.error || 'Terjadi kesalahan. Silakan coba lagi.' };
    }

    return { data };
  } catch (err) {
    return { error: 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.' };
  }
}

// Auth API
export const authApi = {
  signup: (email: string, password: string, parentName: string) =>
    request('/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ email, password, parent_name: parentName }),
    }),

  login: (email: string, password: string) =>
    request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  logout: () =>
    request('/auth/logout', { method: 'POST' }),

  me: () =>
    request('/auth/me'),

  updatePin: (pin: string) =>
    request('/auth/pin', {
      method: 'PUT',
      body: JSON.stringify({ pin }),
    }),

  verifyPin: (pin: string) =>
    request('/auth/pin/verify', {
      method: 'POST',
      body: JSON.stringify({ pin }),
    }),
};

// Children API
export const childrenApi = {
  list: () =>
    request('/children'),

  create: (name: string, birthYear: number, avatarId: number = 1) =>
    request('/children', {
      method: 'POST',
      body: JSON.stringify({ name, birth_year: birthYear, avatar_id: avatarId }),
    }),

  update: (id: number, data: { name?: string; birth_year?: number; avatar_id?: number }) =>
    request(`/children/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    request(`/children/${id}`, { method: 'DELETE' }),

  updateDailyLimit: (id: number, limitMin: number | null) =>
    request(`/children/${id}/daily-limit`, {
      method: 'PUT',
      body: JSON.stringify({ daily_limit_min: limitMin }),
    }),

  getRemainingTime: (id: number) =>
    request(`/children/${id}/remaining-time`),

  getDifficulty: (id: number) =>
    request(`/children/${id}/difficulty`),

  acceptDifficulty: (childId: number, categoryId: number) =>
    request(`/children/${childId}/difficulty/${categoryId}/accept`, { method: 'POST' }),

  dismissDifficulty: (childId: number, categoryId: number) =>
    request(`/children/${childId}/difficulty/${categoryId}/dismiss`, { method: 'POST' }),

  toggleLeaderboardOptIn: (id: number, optIn: boolean) =>
    request(`/children/${id}/leaderboard-opt-in`, {
      method: 'PUT',
      body: JSON.stringify({ opt_in: optIn }),
    }),
};

// Leaderboard API
export const leaderboardApi = {
  getWeekly: (childId: number) =>
    request(`/leaderboard?child_id=${childId}`),
};

// Categories API
export const categoriesApi = {
  list: () =>
    request('/categories'),
};

// Lessons API
export const lessonsApi = {
  listByCategory: (categorySlug: string, ageGroup: string = 'builders', childId?: number) => {
    let url = `/categories/${categorySlug}/lessons?age_group=${ageGroup}`;
    if (childId) url += `&child_id=${childId}`;
    return request(url);
  },

  getById: (id: number) =>
    request(`/lessons/${id}`),

  complete: (lessonId: number, childId: number, finalScore: number, maxScore: number, timeSpentSec: number) =>
    request(`/lessons/${lessonId}/complete`, {
      method: 'POST',
      body: JSON.stringify({
        child_id: childId,
        final_score: finalScore,
        max_score: maxScore,
        time_spent_sec: timeSpentSec,
      }),
    }),
};

// Activities API
export const activitiesApi = {
  submit: (activityId: number, childId: number, answer: any, attemptNo: number = 1) =>
    request(`/activities/${activityId}/submit`, {
      method: 'POST',
      body: JSON.stringify({
        child_id: childId,
        answer,
        attempt_no: attemptNo,
      }),
    }),
};

// Progress API
export const progressApi = {
  getChildProgress: (childId: number) =>
    request(`/children/${childId}/progress`),
};

// Badges API
export const badgesApi = {
  getChildBadges: (childId: number) =>
    request(`/children/${childId}/badges`),
};

// Daily Challenge API
export const dailyChallengeApi = {
  getToday: (ageGroup: string, childId?: number) => {
    let url = `/daily-challenge?age_group=${ageGroup}`;
    if (childId) url += `&child_id=${childId}`;
    return request(url);
  },

  submit: (childId: number, challengeId: number, answer: any, score: number) =>
    request('/daily-challenge/submit', {
      method: 'POST',
      body: JSON.stringify({
        child_id: childId,
        challenge_id: challengeId,
        answer,
        score,
      }),
    }),

  getStreak: (childId: number) =>
    request(`/daily-challenge/streak?child_id=${childId}`),
};

// Parent Dashboard API
export const parentApi = {
  getDashboard: () =>
    request('/parent/dashboard'),
};

