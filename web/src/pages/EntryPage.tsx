import { useNavigate } from "react-router-dom";

export default function EntryPage() {
  const navigate = useNavigate();

  return (
    <div style={{
      height: '100vh',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      gap: '20px',
      backgroundColor: '#f5f5f5'
    }}>
      <h1 style={{
        fontSize: '2.5rem',
        marginBottom: '2rem',
        textAlign: 'center',
        color: '#333'
      }}>
        Система управления портом
      </h1>

      <button 
        style={{
          padding: '12px 24px',
          fontSize: '16px',
          cursor: 'pointer',
          minWidth: '300px',
          backgroundColor: '#2196F3',
          color: 'white',
          border: 'none',
          borderRadius: '8px',
          fontWeight: 'bold',
          transition: 'background-color 0.3s'
        }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#1976D2'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#2196F3'}
        onClick={() => navigate("/operator")}
      >
        Войти как оператор порта
      </button>

      <button 
        style={{
          padding: '12px 24px',
          fontSize: '16px',
          cursor: 'pointer',
          minWidth: '300px',
          backgroundColor: '#4CAF50',
          color: 'white',
          border: 'none',
          borderRadius: '8px',
          fontWeight: 'bold',
          transition: 'background-color 0.3s'
        }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#388E3C'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#4CAF50'}
        onClick={() => navigate("/warehouse")}
      >
        Войти как сотрудник склада
      </button>
    </div>
  );
}