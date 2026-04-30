const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export interface MeResponse {
  id: string
  external_user_id: string
  email: string
  api_key: string
  created_at: string
}

export async function getMe(token: string): Promise<MeResponse> {
  const response = await fetch(`${API_URL}/api/me`, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error('Failed to fetch user profile')
  }

  return response.json()
}

export async function refreshApiKey(token: string): Promise<MeResponse> {
  const response = await fetch(`${API_URL}/api/me/refresh-key`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error('Failed to refresh API key')
  }

  return response.json()
}
