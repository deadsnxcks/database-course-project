import type { ReactNode } from 'react'
import { Link, useNavigate } from 'react-router-dom'

type Props = {
  children: ReactNode
}

export default function WarehouseLayout({ children }: Props) {
  const navigate = useNavigate()

  const handleLogout = () => {
    navigate('/')
  }

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      {/* Sidebar */}
      <aside
        style={{
          width: 240,
          background: '#1e293b',
          color: '#fff',
          padding: 20,
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
        }}
      >
        <div>
          <h2 style={{ marginBottom: 24 }}>Работник склада</h2> {/* ← Изменил текст */}

          <nav style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            <Link to="storageloc" style={linkStyle}>Места хранения</Link>
          </nav>
        </div>

        {/* Кнопка выхода */}
        <button
          onClick={handleLogout}
          style={{
            background: '#ef4444',
            color: 'white',
            border: 'none',
            padding: '10px 16px',
            borderRadius: '6px',
            cursor: 'pointer',
            fontSize: '14px',
            fontWeight: '500',
            marginTop: 'auto',
          }}
        >
          Выйти
        </button>
      </aside>

      {/* Content */}
      <main style={{ flex: 1, padding: 24, overflow: 'auto' }}>
        {children}
      </main>
    </div>
  )
}

const linkStyle: React.CSSProperties = {
  color: '#e5e7eb',
  textDecoration: 'none',
  padding: '8px 12px',
  borderRadius: 6,
}