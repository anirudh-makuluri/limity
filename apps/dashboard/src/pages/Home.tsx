import React from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { useNavigate } from 'react-router-dom'
import { useEffect } from 'react'

export default function HomePage() {
  const { isLoading, isAuthenticated, error } = useAuth0()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      // Small delay to ensure auth state is fully settled
      const timer = setTimeout(() => {
        navigate('/dashboard')
      }, 500)
      return () => clearTimeout(timer)
    }
  }, [isLoading, isAuthenticated, navigate])

  if (isLoading) {
    return <div className="text-center py-12">Loading...</div>
  }

  if (error) {
    return <div className="text-center py-12 text-red-600">Error: {error.message}</div>
  }

  return (
    <div className="bg-white">
      <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 md:py-16 lg:px-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Limity Dashboard</h1>
          <p className="text-xl text-gray-600 mb-8">
            Manage your rate limiting API keys
          </p>
          <p className="text-gray-600">Login to get started</p>
        </div>
      </div>
    </div>
  )
}
