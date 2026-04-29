import { createServerClient } from '~/lib/supabase'
import { generateRandomKey } from '~/lib/utils'

const pool = createServerClient()

export async function generateKey(userId: string) {
  const key = generateRandomKey()

  const result = await pool.query(
    'INSERT INTO api_keys (user_id, key) VALUES ($1, $2) RETURNING *',
    [userId, key]
  )

  return result.rows[0]
}

export async function revokeKey(keyId: string, userId: string) {
  const result = await pool.query(
    'UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND user_id = $2 RETURNING *',
    [keyId, userId]
  )

  return result.rows[0]
}

export async function getUserKeys(userId: string) {
  const result = await pool.query(
    'SELECT * FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC',
    [userId]
  )

  return result.rows
}

export async function verifyKey(key: string) {
  const result = await pool.query(
    'SELECT user_id FROM api_keys WHERE key = $1 AND revoked_at IS NULL',
    [key]
  )

  if (result.rows.length === 0) return null
  return result.rows[0].user_id
}
