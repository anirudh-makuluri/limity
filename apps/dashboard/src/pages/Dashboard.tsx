import React, { useEffect, useState } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import { useNavigate } from 'react-router-dom'
import * as keysAPI from '~/api/keys'
import { Activity, Database, Shield, TerminalSquare } from 'lucide-react'

interface ApiKey {
  id: string
  user_id: string
  key: string
  created_at: string
  revoked_at: string | null
}

export default function DashboardPage() {
  const { isLoading, isAuthenticated, user, getIdTokenClaims } = useAuth0()
  const navigate = useNavigate()
  const [keys, setKeys] = useState<ApiKey[]>([])
  const [isGenerating, setIsGenerating] = useState(false)
  const [isLoadingKeys, setIsLoadingKeys] = useState(true)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/')
    }
  }, [isLoading, isAuthenticated, navigate])

  // Load existing keys from backend
  useEffect(() => {
    if (isAuthenticated && user?.sub) {
      loadKeys()
    }
  }, [isAuthenticated, user?.sub])

  const loadKeys = async () => {
    try {
      setIsLoadingKeys(true)
      const claims = await getIdTokenClaims()
      const idToken = claims?.__raw
      console.log('Got ID token:', idToken ? `${idToken.substring(0, 20)}...` : 'empty')
      const fetchedKeys = await keysAPI.getKeys(idToken)
      setKeys(fetchedKeys || [])
    } catch (error) {
      console.error('Failed to load keys:', error)
      setKeys([])
    } finally {
      setIsLoadingKeys(false)
    }
  }

  const handleGenerateKey = async () => {
    setIsGenerating(true)
    try {
      const claims = await getIdTokenClaims()
      const idToken = claims?.__raw
      const newKey = await keysAPI.generateKey(idToken)
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
      const claims = await getIdTokenClaims()
      const idToken = claims?.__raw
      await keysAPI.revokeKey(keyId, idToken)
      setKeys(keys.map(k => k.id === keyId ? { ...k, revoked_at: new Date().toISOString() } : k))
    } catch (error) {
      console.error('Failed to revoke key:', error)
      alert('Failed to revoke key')
    }
  }

  if (isLoading || isLoadingKeys) {
    return <div className="retro-window p-8 text-center">Booting dashboard...</div>
  }

  const titleBar = (name: string, icon?: React.ReactNode) => (
    <div className="retro-titlebar">
      <span className="retro-title-text">{icon}{name}</span>
      <span className="retro-window-controls" aria-hidden>
        <i>_</i><i>□</i><i>X</i>
      </span>
    </div>
  )

  return (
    <div className="retro-desktop">
      <aside className="retro-sidebar">
        <h2 className="retro-sidebar-title">System</h2>
        <p className="retro-sidebar-subtitle">v2.0.98</p>
        <div className="retro-side-actions">
          <button className="retro-side-button">My Computer</button>
          <button className="retro-side-button">Network</button>
          <button className="retro-side-button">API Keys</button>
          <button className="retro-side-button">Usage</button>
        </div>
      </aside>

      <section className="retro-workspace">
        <div className="retro-grid-top">
          <div className="retro-window">
            {titleBar('Security_Properties.sec', <Shield size={14} />)}
            <div className="retro-window-body">
              <p className="retro-lead">Your personal API access credentials. Keep these secret.</p>
              <label className="retro-field-label">Your API Key:</label>
              <div className="retro-key-row">
                <div className="retro-key-mask">{keys.length ? '***********************' : 'No key generated yet'}</div>
                <button className="retro-button" onClick={handleGenerateKey} disabled={isGenerating}>
                  {isGenerating ? 'Working...' : 'Generate'}
                </button>
              </div>
              <div className="retro-toolbar">
                <button className="retro-button" onClick={handleGenerateKey} disabled={isGenerating}>
                  New Key
                </button>
              </div>
            </div>
          </div>

          <div className="retro-stack">
            <div className="retro-window compact">
              {titleBar('NODE_HEALTH')}
              <div className="retro-window-body">
                <div className="retro-status-item"><span><Database size={14} /> Core API</span><i className="dot green" /></div>
                <div className="retro-status-item"><span><Shield size={14} /> DB Cluster</span><i className="dot green" /></div>
                <div className="retro-status-item"><span><Activity size={14} /> Worker A</span><i className={`dot ${keys.some(k => !k.revoked_at) ? 'green' : 'red'}`} /></div>
              </div>
            </div>
            <div className="retro-window compact">
              {titleBar('OPS_LOG.TXT')}
              <div className="retro-window-body retro-log-body">
                <p>[14:00:01] System initializing...</p>
                <p>[14:00:12] Auth handshake complete</p>
                <p>[14:01:11] Key registry in sync</p>
                <p>[14:02:10] {keys.length} key(s) on record</p>
              </div>
            </div>
          </div>
        </div>

        <div className="retro-window">
          {titleBar('Real_Time_Traffic.mon', <Activity size={14} />)}
          <div className="retro-window-body">
            <div className="retro-terminal">
              <p>REQ/SEC: {4800 + keys.length * 7}</p>
              <div className="retro-bars" aria-hidden>
                {Array.from({ length: 18 }).map((_, i) => (
                  <span key={i} style={{ height: `${28 + Math.sin(i / 3) * 18 + i * 6}%` }} />
                ))}
              </div>
            </div>
          </div>
        </div>

        <div className="retro-window">
          {titleBar('API_KEYS.DAT', <TerminalSquare size={14} />)}
          <div className="retro-menubar">File Edit Search Help</div>
          <div className="retro-window-body no-padding">
            {keys.length === 0 ? (
              <div className="retro-empty-state">
                <p>No API keys yet. Generate one to get started.</p>
              </div>
            ) : (
              <div className="overflow-x-auto">
                <table className="retro-table">
                  <thead>
                    <tr>
                      <th>Key ID</th>
                      <th>Created</th>
                      <th>Status</th>
                      <th>Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {keys.map((key) => (
                      <tr key={key.id}>
                        <td>{key.id.substring(0, 16)}...</td>
                        <td>{new Date(key.created_at).toLocaleDateString()}</td>
                        <td>{key.revoked_at ? 'Revoked' : 'Active'}</td>
                        <td>
                          {!key.revoked_at && (
                            <button className="retro-button danger" onClick={() => handleRevokeKey(key.id)}>
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
        </div>
      </section>
      <div className="retro-taskbar">
        <button className="retro-start">Start</button>
        <span className="retro-task-item">Terminal</span>
        <span className="retro-task-item">Keys</span>
        <span className="retro-task-clock">{new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
      </div>
    </div>
  )
}
