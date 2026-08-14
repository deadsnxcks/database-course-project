import type { ReactNode } from 'react'
import { Link, useNavigate } from 'react-router-dom'

type Props = {
  children: ReactNode
}

export default function OperatorLayout({ children }: Props) {
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
          justifyContent: 'space-between', // ← распределяем пространство
        }}
      >
        <div>
          <h2 style={{ marginBottom: 24 }}>Оператор порта</h2>

          <nav style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
            <Link to="vessels" style={linkStyle}>Суда</Link>
            <Link to="cargo-types" style={linkStyle}>Типы грузов</Link>
            <Link to="cargo" style={linkStyle}>Грузы</Link>
            <Link to="operations" style={linkStyle}>Операции</Link>
            <Link to="operation-cargo" style={linkStyle}>Операции ↔ Грузы</Link>
            <Link to="report-cargo-detail" style={linkStyle}>Отчет 1</Link>
            <Link to="report-cargo-type" style={linkStyle}>Отчет 2</Link>
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