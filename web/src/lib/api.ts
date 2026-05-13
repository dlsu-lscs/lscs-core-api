const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers
    },
    ...options
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ error: 'API request failed' }));
    throw new Error(error.error || 'API request failed');
  }

  return res.json();
}

export const api = {
  // Auth
  login: () => {
    window.location.href = `${API_BASE}/auth/google/login`;
  },

  logout: async () => {
    await fetchAPI('/auth/logout', { method: 'POST' });
  },

  getMe: () => fetchAPI<Member>('/auth/me'),

  // Members
  getMember: (id: number) => fetchAPI<Member>(`/auth/members/${id}`),
  updateMember: (id: number, data: Partial<Member>) =>
    fetchAPI<Member>(`/auth/members/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data)
    }),

  updateMe: (data: Partial<Member>) =>
    fetchAPI<Member>('/auth/me', {
      method: 'PUT',
      body: JSON.stringify(data)
    }),

  // Upload
  generateUploadUrl: (contentType: string) =>
    fetchAPI<{ upload_url: string; object_key: string }>('/upload/profile-image', {
      method: 'POST',
      body: JSON.stringify({ content_type: contentType })
    }),

  getImageUrl: (objectKey: string) =>
    fetchAPI<{ url: string; object_key: string; expires_in_seconds: number }>(
      `/upload/profile-image/url?key=${encodeURIComponent(objectKey)}`
    ),

  completeUpload: (objectKey: string) =>
    fetchAPI<{ image_url: string }>('/upload/profile-image/complete', {
      method: 'POST',
      body: JSON.stringify({ object_key: objectKey })
    }),

  deleteImage: () => fetchAPI('/upload/profile-image', { method: 'DELETE' }),

  // API Keys
  listApiKeys: () => fetchAPI<APIKey[]>('/api-keys'),

  revokeApiKey: (id: number) => fetchAPI<{ message: string }>(`/api-keys/${id}`, { method: 'DELETE' }),

  // API Key Request (uses Google OAuth, not session)
  requestKey: (data: RequestKeyRequest) =>
    fetchAPI<RequestKeyResponse>('/request-key', {
      method: 'POST',
      body: JSON.stringify(data)
    })
};

// Types
export interface APIKey {
  api_key_id: number;
  member_email: string;
  project?: string;
  allowed_origin?: string;
  is_dev: boolean;
  is_admin: boolean;
  created_at: string;
  expires_at?: string;
}

export interface RequestKeyRequest {
  project?: string;
  allowed_origin?: string;
  is_dev: boolean;
}

export interface RequestKeyResponse {
  email: string;
  api_key: string;
  expires_at?: string;
}

export interface Member {
  id: number;
  email: string;
  full_name: string;
  nickname?: string;
  image_url?: string;
  committee_id?: string;
  committee_name?: string;
  division_id?: string;
  division_name?: string;
  position_id?: string;
  position_name?: string;
  house_name?: string;
  contact_number?: string;
  college?: string;
  program?: string;
  interests?: string;
  discord?: string;
  fb_link?: string;
  telegram?: string;
}
