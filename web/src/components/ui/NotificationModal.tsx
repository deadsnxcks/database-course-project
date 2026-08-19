// components/ui/NotificationModal.tsx
import { useEffect } from 'react';
import Modal from './Modal';

export type NotificationType = 'success' | 'error' | 'info' | 'warning';

interface NotificationModalProps {
  isOpen: boolean;
  onClose: () => void;
  type: NotificationType;
  title: string;
  message: string;
  details?: string;
  onAction?: () => void;
  actionText?: string;
  autoClose?: number;
}

export default function NotificationModal({
  isOpen,
  onClose,
  type,
  title,
  message,
  details,
  onAction,
  actionText,
  autoClose = type === 'success' ? 3000 : 10000
}: NotificationModalProps) {
  // Автозакрытие
  useEffect(() => {
    if (isOpen && autoClose > 0) {
      const timer = setTimeout(() => {
        onClose();
      }, autoClose);
      
      return () => clearTimeout(timer);
    }
  }, [isOpen, autoClose, onClose]);

  // Стили в зависимости от типа
  const typeStyles = {
    success: {
      bg: 'bg-green-50',
      border: 'border-green-200',
      text: 'text-green-800',
      iconBg: 'bg-green-100',
      iconColor: 'text-green-600',
      button: 'bg-green-500 hover:bg-green-600',
    },
    error: {
      bg: 'bg-red-50',
      border: 'border-red-200',
      text: 'text-red-800',
      iconBg: 'bg-red-100',
      iconColor: 'text-red-600',
      button: 'bg-red-500 hover:bg-red-600',
    },
    info: {
      bg: 'bg-blue-50',
      border: 'border-blue-200',
      text: 'text-blue-800',
      iconBg: 'bg-blue-100',
      iconColor: 'text-blue-600',
      button: 'bg-blue-500 hover:bg-blue-600',
    },
    warning: {
      bg: 'bg-yellow-50',
      border: 'border-yellow-200',
      text: 'text-yellow-800',
      iconBg: 'bg-yellow-100',
      iconColor: 'text-yellow-600',
      button: 'bg-yellow-500 hover:bg-yellow-600',
    },
  };

  const currentStyle = typeStyles[type];

  // Иконки
  const icons = {
    success: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
      </svg>
    ),
    error: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
      </svg>
    ),
    info: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    ),
    warning: (
      <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.998-.833-2.732 0L4.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
      </svg>
    ),
  };

  // Заголовки по умолчанию
  const defaultTitles = {
    success: 'Успешно',
    error: 'Ошибка',
    info: 'Информация',
    warning: 'Внимание',
  };

  const displayTitle = title || defaultTitles[type];

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className={`p-6 ${currentStyle.bg} border ${currentStyle.border} rounded-lg`}>
        <div className="flex items-start mb-4">
          <div className={`${currentStyle.iconBg} p-2 rounded-full mr-3 flex-shrink-0`}>
            <div className={currentStyle.iconColor}>
              {icons[type]}
            </div>
          </div>
          <div className="flex-1">
            <h3 className={`text-xl font-bold ${currentStyle.text}`}>
              {displayTitle}
            </h3>
            <p className={`mt-2 ${currentStyle.text}`}>
              {message}
            </p>
            
            {details && (
              <div className="mt-4">
                <p className="text-sm text-gray-600 mb-1">Детали:</p>
                <div className="bg-gray-100 p-3 rounded text-sm font-mono overflow-auto max-h-32">
                  {details}
                </div>
              </div>
            )}
          </div>
          
          
        </div>
        
        {/* Прогресс-бар для автозакрытия */}
        {autoClose > 0 && (
          <div className="w-full bg-gray-200 rounded-full h-1 mb-4">
            <div 
              className={`h-1 rounded-full ${currentStyle.iconBg} transition-all duration-100 ease-linear`}
              style={{ 
                width: isOpen ? '100%' : '0%',
                transitionDuration: `${autoClose}ms`
              }}
            />
          </div>
        )}
        
        {/* Кнопки действий */}
        <div className="flex justify-end gap-3">
          {onAction && (
            <button
              onClick={() => {
                onAction();
                onClose();
              }}
              className={`px-4 py-2 ${currentStyle.button} text-white rounded transition`}
            >
              {actionText || (type === 'error' ? 'Повторить' : 'ОК')}
            </button>
          )}
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400 transition"
          >
            Закрыть
          </button>
        </div>
      </div>
    </Modal>
  );
}

// Экспортируем тип отдельно
export type { NotificationModalProps };