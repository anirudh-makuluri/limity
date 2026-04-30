import React, { useMemo, useState } from 'react'

interface ApiKeyEnvSnippetProps {
  apiKey: string
  onRefresh: () => Promise<void>
  isRefreshing: boolean
}

export default function ApiKeyEnvSnippet({ apiKey, onRefresh, isRefreshing }: ApiKeyEnvSnippetProps) {
  const [copied, setCopied] = useState(false)

  const envText = useMemo(() => {
    if (!apiKey) {
      return '# API key unavailable yet\nLIMITY_API_KEY='
    }

    return [
      '# Limity API Configuration',
      `LIMITY_API_KEY=${apiKey}`,
      'LIMITY_BASE_URL=https://api.limity.dev',
    ].join('\n')
  }, [apiKey])

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(envText)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1500)
    } catch (error) {
      console.error('Failed to copy env snippet:', error)
    }
  }

  return (
    <div className="retro-window">
      <div className="retro-titlebar">
        <span className="retro-title-text">ENV_SNIPPET.env</span>
        <span className="retro-window-controls" aria-hidden>
          <i>_</i><i>□</i><i>X</i>
        </span>
      </div>
      <div className="retro-window-body">
        <p className="retro-lead">Copy and paste this into your app's `.env` file.</p>
        <pre className="retro-env-block">{envText}</pre>
        <div className="retro-toolbar">
          <button className="retro-button" type="button" onClick={onRefresh} disabled={isRefreshing}>
            {isRefreshing ? 'Refreshing...' : 'Refresh Key'}
          </button>
          <button className="retro-button" type="button" onClick={handleCopy}>
            {copied ? 'Copied!' : 'Copy .env'}
          </button>
        </div>
      </div>
    </div>
  )
}
