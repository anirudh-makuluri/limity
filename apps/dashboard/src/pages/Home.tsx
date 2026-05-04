import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import * as keysAPI from '~/api/keys'
import { Activity, Database, Shield } from 'lucide-react'
import { useAuth } from '~/lib/useAuth'
import type { MeResponse } from '~/api/keys'
import ApiKeyEnvSnippet from '~/components/ApiKeyEnvSnippet'

export default function DashboardPage() {
  const { isLoading, isAuthenticated, user, getAccessToken } = useAuth()
  const navigate = useNavigate()
  const [me, setMe] = useState<MeResponse | null>(null)
  const [isLoadingMe, setIsLoadingMe] = useState(true)
  const [isRefreshingKey, setIsRefreshingKey] = useState(false)

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/')
    }
  }, [isLoading, isAuthenticated, navigate])

  // Load existing keys from backend
  useEffect(() => {
    if (isAuthenticated && user?.id) {
      loadMe()
    }
  }, [isAuthenticated, user?.id])

  const loadMe = async () => {
    try {
      setIsLoadingMe(true)
      const accessToken = await getAccessToken()
      const userProfile = await keysAPI.getMe(accessToken)
      setMe(userProfile)
    } catch (error) {
      console.error('Failed to load user profile:', error)
      setMe(null)
    } finally {
      setIsLoadingMe(false)
    }
  }

  const handleRefreshKey = async () => {
    try {
      setIsRefreshingKey(true)
      const accessToken = await getAccessToken()
      const updatedProfile = await keysAPI.refreshApiKey(accessToken)
      setMe(updatedProfile)
    } catch (error) {
      console.error('Failed to refresh API key:', error)
    } finally {
      setIsRefreshingKey(false)
    }
  }

  if (isLoading || isLoadingMe) {
    return <div className="retro-window p-8 text-center">Booting dashboard...</div>
  }
  const primaryKey = me?.api_key ?? ''

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
        <p className="retro-sidebar-subtitle">v0.1.2</p>
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
                <div className="retro-key-mask">{primaryKey || 'No key available'}</div>
              </div>
            </div>
          </div>

          <div className="retro-stack">
            <div className="retro-window compact">
              {titleBar('NODE_HEALTH')}
              <div className="retro-window-body">
                <div className="retro-status-item"><span><Database size={14} /> Core API</span><i className="dot green" /></div>
                <div className="retro-status-item"><span><Shield size={14} /> DB Cluster</span><i className="dot green" /></div>
                <div className="retro-status-item"><span><Activity size={14} /> Worker A</span><i className={`dot ${primaryKey ? 'green' : 'red'}`} /></div>
              </div>
            </div>
            <div className="retro-window compact">
              {titleBar('OPS_LOG.TXT')}
              <div className="retro-window-body retro-log-body">
                <p>[14:00:01] System initializing...</p>
                <p>[14:00:12] Auth handshake complete</p>
                <p>[14:01:11] Key registry in sync</p>
                <p>[14:02:10] {primaryKey ? '1 key on record' : '0 keys on record'}</p>
              </div>
            </div>
          </div>
        </div>
        <ApiKeyEnvSnippet
          apiKey={primaryKey}
          onRefresh={handleRefreshKey}
          isRefreshing={isRefreshingKey}
        />
      </section>
      <div className="retro-taskbar absolute bottom-0 left-0 right-0 flex items-center gap-4 px-4 py-2 bg-gray-800 text-white">
        <button className="retro-start">Start</button>
        <span className="retro-task-item">Terminal</span>
        <span className="retro-task-item">Keys</span>
        <span className="retro-task-clock">{new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
      </div>
    </div>
  )
}
