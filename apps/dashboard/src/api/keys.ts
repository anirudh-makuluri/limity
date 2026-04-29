const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

interface ApiKey {
  id: string
  user_id: string
  key: string
  created_at: string
  revoked_at: string | null
}

export async function generateKey(token: string): Promise<ApiKey> {
  const response = await fetch(`${API_URL}/api/keys/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error('Failed to generate key')
  }

  return response.json()
}

export async function revokeKey(keyId: string, token: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/keys/${keyId}/revoke`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    throw new Error('Failed to revoke key')
  }
}

export async function getKeys(token: string): Promise<ApiKey[]> {
  console.log('getKeys: sending request with token:', token ? `${token.substring(0, 20)}...` : 'empty')
  const response = await fetch(`${API_URL}/api/keys`, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })

  if (!response.ok) {
    const errorData = await response.text()
    console.error('getKeys error response:', response.status, errorData)
    throw new Error('Failed to fetch keys')
  }

  return response.json()
}
