import { useEffect, useState } from 'react'
import { useAuth } from '~/lib/useAuth'

export default function HomePage() {
  const { isLoading } = useAuth()
  const [viewport, setViewport] = useState({ width: 800, height: 600 })

  useEffect(() => {
    const updateViewport = () => {
      setViewport({
        width: window.innerWidth,
        height: window.innerHeight,
      })
    }

    updateViewport()
    window.addEventListener('resize', updateViewport)
    return () => window.removeEventListener('resize', updateViewport)
  }, [])

  if (isLoading) {
    return <div className="retro-window p-8 text-center">Booting homepage...</div>
  }

  return (
    <section className="retro-landing retro-landing-fit">
      <div className="retro-landing-hero">
        <div className="retro-hero-copy">
          <div className="retro-pill" aria-label="New version announcement">
            <div className="retro-pill-track">
              <span>*** NEW! VERSION 0.1.1 IS HERE! ***</span>
            </div>
          </div>
          <h2 className="retro-hero-title alt">
            Rate Limiting,<br />
            <em>Reimagined</em> for<br />
            Developers
          </h2>
          <p className="retro-hero-sub alt">
            Experience ~1ms latency with zero-config memory-based limiting.
            Built for high-velocity engineering teams who demand precision and performance.
            <strong> Compatible with all major browsers!</strong>
          </p>
          <div className="retro-cta-row">
            <button className="retro-cta">DOWNLOAD NOW</button>
            <button className="retro-cta secondary">Help Contents</button>
          </div>
          <ul className="retro-bullets">
            <li>ULTRA LOW LATENCY SYSTEM</li>
            <li>IN-MEMORY SPEED TECHNOLOGY</li>
          </ul>
        </div>

        <div className="retro-hero-code-col">
          <article className="retro-window">
            <div className="retro-titlebar">
              <span className="retro-title-text">C:\WINDOWS\NOTEPAD.EXE - integration.ts</span>
              <span className="retro-window-controls" aria-hidden><i>_</i><i>□</i><i>X</i></span>
            </div>
            <div className="retro-window-body">
              <pre className="retro-code terminal">{`C:\\> import { Limity } from '@limity/core';

// Initialize with ~1ms latency overhead
const limiter = new Limity({
  rate: 10,
  window: '1m',
  strategy: 'sliding-window'
});

// Seamless middleware integration
export const middleware = async (req, res) => {
  const { success } = await limiter.check(req.ip);

  if (!success) {
    return res.status(429).json({ error: 'Limit exceeded' });
  }
};`}</pre>
            </div>
          </article>
          <div className="retro-stats-row">
            <div className="retro-stat-box">
              <strong>0.82ms</strong>
              <span>AVERAGE SPEED</span>
            </div>
            <div className="retro-stat-box">
              <strong>10M+</strong>
              <span>HITS PER SEC</span>
            </div>
          </div>
        </div>
      </div>

      <div className="retro-feature-grid alt">
        <article className="retro-window">
          <div className="retro-window-body">
            <h3>Zero-Config</h3>
            <p>Drop-in memory-based limiting that scales horizontally without external DB dependencies.</p>
          </div>
        </article>
        <article className="retro-window">
          <div className="retro-window-body">
            <h3>Total Control</h3>
            <p>Native monitoring support and precise thresholds so your team stays fully in charge.</p>
          </div>
        </article>
        <article className="retro-window">
          <div className="retro-window-body">
            <h3>Options</h3>
            <p>Switch between Token Bucket, Sliding Window, and Leaky Bucket with one-line config.</p>
          </div>
        </article>
      </div>

      <footer className="retro-site-footer">
        <div>
          <h4>Limity Engineering Corp.</h4>
          <p>
            Copyright © {new Date().getFullYear()} Limity. All Rights Reserved. Optimized for {viewport.width}x{viewport.height} resolution.
          </p>
        </div>
        {/* <nav>
          <a href="#">Privacy Policy</a>
          <a href="#">Terms of Use</a>
          <a href="#">Site Map</a>
          <a href="#">Webmaster</a>
        </nav> */}
      </footer>
    </section>
  )
}
