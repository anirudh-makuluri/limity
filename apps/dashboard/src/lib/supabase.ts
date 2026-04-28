import { createClient } from '@supabase/supabase-js'

const supabaseUrl = process.env.VITE_SUPABASE_URL || ''
const supabaseKey = process.env.VITE_SUPABASE_ANON_KEY || ''

export const supabase = createClient(supabaseUrl, supabaseKey)

// Server-side client with admin privileges
export const createServerClient = () => {
  const serviceKey = process.env.SUPABASE_SERVICE_KEY || ''
  return createClient(supabaseUrl, serviceKey)
}
