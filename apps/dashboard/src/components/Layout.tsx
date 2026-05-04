import React, { ReactNode } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '~/lib/useAuth'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { logout, user, isLoading } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  return (
    <div className="retro-shell">
      <nav className="retro-topbar">
        <div className="retro-topbar-inner">
          <div className="retro-topbar-row">
            <div className="retro-nav-left">
              <h1 className="retro-brand">Limity™</h1>
              {/* <a href="#" className="retro-nav-link">Docs</a>
              <a href="#" className="retro-nav-link">Changelog</a>
              <a href="#" className="retro-nav-link">Pricing</a>
              <a href="#" className="retro-nav-link">Enterprise</a> */}
            </div>
            <div>
              {isLoading ? (
                <p className="retro-user-text">Loading...</p>
              ) : user ? (
                <div className="retro-user-controls">
                  <span className="retro-user-text">{user.email}</span>
                  {location.pathname === '/' ? (
                    <button
                      onClick={() => navigate('/home')}
                      className="retro-button"
                    >
                      Home
                    </button>
                  ) : (
                    <button
                      onClick={() => logout()}
                      className="retro-button"
                    >
                      Logout
                    </button>
                  )}
                </div>
              ) : (
                <div className="retro-user-controls">
                  <button onClick={() => navigate('/auth')} className="retro-button">Sign In</button>
                  <button onClick={() => navigate('/auth')} className="retro-button retro-join">Join Now!</button>
                </div>
              )}
            </div>
          </div>
        </div>
      </nav>

      <main className="retro-main">
        {children}
      </main>
    </div>
  )
}
