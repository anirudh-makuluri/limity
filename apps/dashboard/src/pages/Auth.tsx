import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '~/lib/useAuth'

export default function AuthPage() {
  const { signInWithPassword, signUpWithPassword } = useAuth()
  const navigate = useNavigate()
  const [mode, setMode] = useState<'signin' | 'signup'>('signin')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    setMessage('')
    setError('')
    setIsSubmitting(true)

    try {
      if (mode === 'signin') {
        await signInWithPassword(email, password)
        navigate('/home')
      } else {
        await signUpWithPassword(email, password)
        setMessage('Account created. Check your inbox for email confirmation.')
      }
    } catch (err) {
      const text = err instanceof Error ? err.message : 'Authentication failed'
      setError(text)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <section className="retro-auth-wrap">
      <div className="retro-window retro-auth-window">
        <div className="retro-titlebar">
          <span className="retro-title-text">AUTHENTICATION</span>
          <span className="retro-window-controls" aria-hidden>
            <i>_</i><i>□</i><i>X</i>
          </span>
        </div>
        <div className="retro-window-body">
          <div className="retro-auth-tabs">
            <button
              className={`retro-button ${mode === 'signin' ? 'retro-auth-tab-active' : ''}`}
              onClick={() => setMode('signin')}
              type="button"
            >
              Sign In
            </button>
            <button
              className={`retro-button ${mode === 'signup' ? 'retro-auth-tab-active' : ''}`}
              onClick={() => setMode('signup')}
              type="button"
            >
              Create Account
            </button>
          </div>

          <form onSubmit={handleSubmit} className="retro-auth-form">
            <label className="retro-auth-label" htmlFor="auth-email">Email</label>
            <input
              id="auth-email"
              className="retro-auth-input"
              type="email"
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />

            <label className="retro-auth-label" htmlFor="auth-password">Password</label>
            <input
              id="auth-password"
              className="retro-auth-input"
              type="password"
              autoComplete={mode === 'signin' ? 'current-password' : 'new-password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={8}
            />

            <button className="retro-button" type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Working...' : mode === 'signin' ? 'Sign In' : 'Create Account'}
            </button>
          </form>

          {message && <p className="retro-auth-message">{message}</p>}
          {error && <p className="retro-auth-error">{error}</p>}
        </div>
      </div>
    </section>
  )
}
