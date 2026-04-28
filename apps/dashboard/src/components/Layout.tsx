import React, { ReactNode } from 'react'
import { useAuth0 } from '@auth0/auth0-react'

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { loginWithRedirect, logout, user, isLoading } = useAuth0()

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-gray-900">Limity</h1>
            </div>
            <div>
              {isLoading ? (
                <p className="text-sm text-gray-700">Loading...</p>
              ) : user ? (
                <div className="flex items-center gap-4">
                  <span className="text-sm text-gray-700">{user.email}</span>
                  <button
                    onClick={() =>
                      logout({ logoutParams: { returnTo: window.location.origin } })
                    }
                    className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700"
                  >
                    Logout
                  </button>
                </div>
              ) : (
                <button
                  onClick={() => loginWithRedirect()}
                  className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
                >
                  Login
                </button>
              )}
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  )
}
