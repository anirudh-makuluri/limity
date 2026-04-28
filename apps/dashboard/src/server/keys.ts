import { createServerClient } from '~/lib/supabase'
import { generateRandomKey } from '~/lib/utils'

const supabase = createServerClient()

export async function generateKey(userId: string) {
  const key = generateRandomKey()

  const { data, error } = await supabase
    .from('api_keys')
    .insert({
      user_id: userId,
      key,
    })
    .select()
    .single()

  if (error) throw error
  return data
}

export async function revokeKey(keyId: string, userId: string) {
  const { data, error } = await supabase
    .from('api_keys')
    .update({ revoked_at: new Date().toISOString() })
    .eq('id', keyId)
    .eq('user_id', userId)
    .select()
    .single()

  if (error) throw error
  return data
}

export async function getUserKeys(userId: string) {
  const { data, error } = await supabase
    .from('api_keys')
    .select('*')
    .eq('user_id', userId)
    .order('created_at', { ascending: false })

  if (error) throw error
  return data
}

export async function verifyKey(key: string) {
  const { data, error } = await supabase
    .from('api_keys')
    .select('user_id')
    .eq('key', key)
    .is('revoked_at', null)
    .single()

  if (error) return null
  return data?.user_id
}
