import React, { useEffect, useState } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { useNavigate } from 'react-router-dom'

interface ApiKey {
  id: string
  key: string
  created_at: string
  revoked_at: string | null
}

export default function DashboardPage() {
  const { isLoading, isAuthenticated, user, getAccessTokenSilently } = useAuth0()
  const navigate = useNavigate()
  const [keys, setKeys] = useState<ApiKey[]>([
    {
      id: '550e8400-e29b-41d4-a716-446655440000',
      key: 'limity_8fkd9sj2k1j29k1j2k1j29k',
      created_at: new Date().toISOString(),
      revoked_at: null,
    },
  ])
  const [isGenerating, setIsGenerating] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/')
    }
  }, [isLoading, isAuthenticated, navigate])

  const handleGenerateKey = async () => {
    setIsGenerating(true)
    try {
      const newKey: ApiKey = {
        id: Math.random().toString(36).substring(7),
        key: `limity_${Math.random().toString(36).substring(2, 28)}`,
        created_at: new Date().toISOString(),
        revoked_at: null,
      }
      setKeys([newKey, ...keys])
      alert(`Your new API key: ${newKey.key}\n\nMake sure to copy it now - you won't be able to see it again!`)
    } catch (error) {
      console.error('Failed to generate key:', error)
      alert('Failed to generate key')
    } finally {
      setIsGenerating(false)
    }
  }

  const handleRevokeKey = async (keyId: string) => {
    if (!confirm('Are you sure you want to revoke this key?')) return

    try {
      setKeys(keys.map(k => k.id === keyId ? { ...k, revoked_at: new Date().toISOString() } : k))
    } catch (error) {
      console.error('Failed to revoke key:', error)
      alert('Failed to revoke key')
    }
  }

  if (isLoading) {
    return <div className="text-center py-12">Loading...</div>
  }

  return (
    <div>
      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900">API Keys</h2>
        <p className="text-gray-600">Manage your rate limiting API keys</p>
      </div>

      <div className="mb-8">
        <button
          onClick={handleGenerateKey}
          disabled={isGenerating}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 disabled:opacity-50"
        >
          {isGenerating ? 'Generating...' : 'Generate New Key'}
        </button>
      </div>

      {keys.length === 0 ? (
        <div className="text-center py-12 border border-gray-200 rounded-lg bg-gray-50">
          <p className="text-gray-600">No API keys yet. Generate one to get started!</p>
        </div>
      ) : (
        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Key ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Created
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {keys.map((key) => (
                <tr key={key.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                    {key.id.substring(0, 16)}...
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                    {new Date(key.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        key.revoked_at
                          ? 'bg-red-100 text-red-800'
                          : 'bg-green-100 text-green-800'
                      }`}
                    >
                      {key.revoked_at ? 'Revoked' : 'Active'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    {!key.revoked_at && (
                      <button
                        onClick={() => handleRevokeKey(key.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Revoke
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
