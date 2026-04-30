import React, { ReactNode } from 'react'
import { useAuth0 } from '@auth0/auth0-react'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { loginWithRedirect, logout, user, isLoading } = useAuth0()

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
                  <button
                    onClick={() =>
                      logout({ logoutParams: { returnTo: window.location.origin } })
                    }
                    className="retro-button"
                  >
                    Logout
                  </button>
                </div>
              ) : (
                <div className="retro-user-controls">
                  <button onClick={() => loginWithRedirect()} className="retro-button">Sign In</button>
                  <button onClick={() => loginWithRedirect()} className="retro-button retro-join">Join Now!</button>
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
