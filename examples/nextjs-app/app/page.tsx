'use client';

import { useState } from 'react';

export default function Home() {
  const [loading, setLoading] = useState(false);
  const [response, setResponse] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  const testEndpoint = async (endpoint: string) => {
    setLoading(true);
    setError(null);
    setResponse(null);

    try {
      const res = await fetch(endpoint);
      const data = await res.json();

      if (!res.ok) {
        setError(`${res.status}: ${data.error}`);
      } else {
        setResponse(data);
      }
    } catch (err) {
      setError(String(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '2rem' }}>
      <h1>🚀 Limity Next.js Example</h1>

      <section style={{ marginTop: '2rem' }}>
        <h2>Test Endpoints</h2>

        <div style={{ display: 'grid', gap: '1rem' }}>
          <button
            onClick={() => testEndpoint('/api/data')}
            disabled={loading}
            style={{
              padding: '0.75rem',
              fontSize: '1rem',
              cursor: loading ? 'not-allowed' : 'pointer',
              backgroundColor: '#0070f3',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              opacity: loading ? 0.5 : 1,
            }}
          >
            GET /api/data (100 req/min)
          </button>

          <button
            onClick={() => testEndpoint('/api/limited')}
            disabled={loading}
            style={{
              padding: '0.75rem',
              fontSize: '1rem',
              cursor: loading ? 'not-allowed' : 'pointer',
              backgroundColor: '#f5a623',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              opacity: loading ? 0.5 : 1,
            }}
          >
            GET /api/limited (10 req/min)
          </button>
        </div>
      </section>

      {error && (
        <div
          style={{
            marginTop: '2rem',
            padding: '1rem',
            backgroundColor: '#fee',
            color: '#c33',
            borderRadius: '4px',
            fontFamily: 'monospace',
            whiteSpace: 'pre-wrap',
            wordBreak: 'break-all',
          }}
        >
          ❌ Error: {error}
        </div>
      )}

      {response && (
        <div
          style={{
            marginTop: '2rem',
            padding: '1rem',
            backgroundColor: '#f0f0f0',
            borderRadius: '4px',
            fontFamily: 'monospace',
            whiteSpace: 'pre-wrap',
            wordBreak: 'break-all',
          }}
        >
          ✅ Response:
          {'\n'}
          {JSON.stringify(response, null, 2)}
        </div>
      )}

      <section style={{ marginTop: '3rem' }}>
        <h2>How It Works</h2>
        <ul>
          <li>Each endpoint has its own rate limit</li>
          <li>Limits are per IP address</li>
          <li>Exceeding limit returns 429 status</li>
          <li>Check response headers for rate limit info</li>
          <li>Supports both in-memory and hosted modes</li>
        </ul>
      </section>

      <section style={{ marginTop: '3rem', color: '#666' }}>
        <h3>To test rate limiting:</h3>
        <ol>
          <li>Click a button multiple times rapidly</li>
          <li>After exceeding the limit, you'll see a 429 error</li>
          <li>Wait for the window to reset (~60 seconds)</li>
          <li>You can make requests again</li>
        </ol>
      </section>
    </div>
  );
}
