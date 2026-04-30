import { supabase } from '~/lib/supabase'
 
export async function getUserByAPIKey(key: string) {
  const { data, error } = await supabase
    .from('users')
    .select('id')
    .eq('api_key', key)
    .maybeSingle()
  if (error) throw error
  return data?.id ?? null
}
