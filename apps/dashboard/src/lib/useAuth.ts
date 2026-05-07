import { useCallback, useEffect, useState } from 'react'
import type { Session, User } from '@supabase/supabase-js'
import { supabase } from './supabase'

function getEmailRedirectUrl(): string {
  const configured = import.meta.env.VITE_AUTH_REDIRECT_URL
  if (configured && configured.trim().length > 0) {
    return configured.trim()
  }

  if (typeof window !== 'undefined' && window.location?.origin) {
    return `${window.location.origin}/auth`
  }

  return 'http://localhost:5173/auth'
}

export function useAuth() {
  const [session, setSession] = useState<Session | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    let isMounted = true

    supabase.auth.getSession().then(({ data }) => {
      if (!isMounted) return
      setSession(data.session ?? null)
      setIsLoading(false)
    })

    const { data } = supabase.auth.onAuthStateChange((_event, nextSession) => {
      setSession(nextSession)
      setIsLoading(false)
    })

    return () => {
      isMounted = false
      data.subscription.unsubscribe()
    }
  }, [])

  const logout = useCallback(async () => {
    const { error } = await supabase.auth.signOut()
    if (error) throw error
  }, [])

  const getAccessToken = useCallback(async () => {
    const { data } = await supabase.auth.getSession()
    return data.session?.access_token ?? ''
  }, [])

  const signInWithPassword = useCallback(async (email: string, password: string) => {
    const { error } = await supabase.auth.signInWithPassword({ email, password })
    if (error) throw error
  }, [])

  const signUpWithPassword = useCallback(async (email: string, password: string) => {
    const { error } = await supabase.auth.signUp({
      email,
      password,
      options: {
        emailRedirectTo: getEmailRedirectUrl(),
      },
    })
    if (error) throw error
  }, [])

  const user: User | null = session?.user ?? null

  return {
    isLoading,
    isAuthenticated: !!session,
    user,
    logout,
    getAccessToken,
    signInWithPassword,
    signUpWithPassword,
  }
}
