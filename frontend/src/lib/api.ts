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
};

// Categories API
export const categoriesApi = {
  list: () =>
    request('/categories'),
};
