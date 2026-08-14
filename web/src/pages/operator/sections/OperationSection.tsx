import { useState, useEffect } from 'react'
import { Table } from '../../../components/ui/Table'
import type { TableColumn } from '../../../components/ui/Table'
import CreateOperationForm from '../forms/operation/CreateOperationForm'
import EditOperationForm from '../forms/operation/EditOperationForm'
import Modal from '../../../components/ui/Modal'
import NotificationModal, { type NotificationType } from '../../../components/ui/NotificationModal'
import * as operationApi from '../../../api/operationAPI'
import type { Operation } from '../../../api/operationAPI'

interface NotificationState {
  isOpen: boolean;
  type: NotificationType;
  title: string;
  message: string;
  details?: string;
  onAction?: () => void;
  actionText?: string;
}

const formatErrorMessage = (error: any) => {
  switch (error.code) {
    case 'AlreadyExists':
      return {
        title: 'Конфликт',
        message: 'Такая операция уже существует',
      }

    case 'NotFound':
      return {
        title: 'Не найдено',
        message: 'Операция не найдена',
      }

    case 'InvalidArgument':
      return {
        title: 'Некорректные данные',
        message: error.message,
      }

    case 'FailedPrecondition':
      return {
        title: 'Невозможно удалить',
        message: 'Операция используется',
      }

    default:
      return {
        title: 'Ошибка',
        message: error.message || 'Неизвестная ошибка',
      }
  }
}

const formatDate = (value: any): string => {
  if (!value) return '-';

  console.log('Formatting date value:', value);
  
  if (value && typeof value === 'object' && 'seconds' in value) {
    const seconds = Number(value.seconds);
    const nanos = Number(value.nanos ?? 0); // 🔥 ВОТ ЭТО ГЛАВНОЕ

    const timestamp = seconds * 1000 + Math.floor(nanos / 1_000_000);
    const date = new Date(timestamp);

    if (isNaN(date.getTime())) {
      return 'Неверная дата';
    }

    return new Intl.DateTimeFormat('ru-RU', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    }).format(date);
  }

  
  if (typeof value === 'string') {
  // Превращаем postgres timestamp в валидный ISO
  const isoDate = value
    .replace(' ', 'T')
    + 'Z'; // или +03:00 если у тебя локальное время

  const date = new Date(isoDate);

  if (!isNaN(date.getTime())) {
      return new Intl.DateTimeFormat('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      }).format(date);
    }
  }

  
  // Если это число (timestamp)
  if (typeof value === 'number') {
    const date = new Date(value);
    if (!isNaN(date.getTime())) {
      return new Intl.DateTimeFormat('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      }).format(date);
    }
  }
  
  return 'Неверный формат даты';
};

export default function OperationSection() {
  const [search, setSearch] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editOperation, setEditOperation] = useState<Operation | null>(null)
  const [operations, setOperations] = useState<Operation[]>([])
  const [loading, setLoading] = useState(true)
  const [notification, setNotification] = useState<NotificationState>({
    isOpen: false,
    type: 'info',
    title: '',
    message: '',
  })

  const closeNotification = () => {
    setNotification(prev => ({ ...prev, isOpen: false }))
  }

  const showNotification = (
    type: NotificationType,
    title: string,
    message: string,
    details?: string,
    onAction?: () => void,
    actionText?: string
  ) => {
    setNotification({
      isOpen: true,
      type,
      title,
      message,
      details,
      onAction,
      actionText,
    })
  }

  const showSuccess = (message: string, title: string = 'Успешно') => {
    showNotification('success', title, message)
  }

  const showError = (error: any, onRetry?: () => void) => {
    const formatted = formatErrorMessage(error)

    setNotification({
      isOpen: true,
      type: 'error',
      title: formatted.title,
      message: formatted.message,
      onAction: onRetry,
      actionText: onRetry ? 'Повторить' : undefined,
    })
  }

  useEffect(() => {
    fetchOperations()
  }, [])

  const fetchOperations = async () => {
    try {
      setLoading(true)
      const data = await operationApi.listOperations()
      setOperations(Array.isArray(data) ? data : [])
    } catch (err: any) {
      showError(err, fetchOperations)
    } finally {
      setLoading(false)
    }
  }

  const filtered = (operations || []).filter((op) =>
    op.title.toLowerCase().includes(search.toLowerCase())
  )

  const columns: TableColumn<Operation>[] = [
    { key: 'id', title: 'ID' },
    { 
      key: 'title', 
      title: 'Название',
      render: (value) => <span>{value}</span>
    },
    { 
      key: 'createdAt', 
      title: 'Дата создания',
      render: (value) => formatDate(value as string)
    },
    {
      key: 'actions',
      title: 'Действия',
      render: (_, op) => (
        <div className="flex gap-2">
          <button
            className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
            onClick={() => setEditOperation(op)}
          >
            Редактировать
          </button>
          <button
            className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition text-sm"
            onClick={() => handleDelete(op.id)}
          >
            Удалить
          </button>
        </div>
      ),
    },
  ]

  const handleCreate = async (data: { title: string }) => {
    try {
      const { id } = await operationApi.createOperation(data.title)
      const newOperation: Operation = {
        id,
        title: data.title,
        createdAt: new Date().toISOString()
      }
      setOperations([...operations, newOperation])
      setShowCreateForm(false)
      showSuccess('Операция успешно создана')
    } catch (err: any) {
      showError(err)
    }
  }

  const handleEditSubmit = async (data: { id: number; title: string }) => {
    try {
      await operationApi.updateOperation(data.id, { title: data.title })
      setOperations(
        operations.map((op) => (op.id === data.id ? { ...op, title: data.title } : op))
      )
      setEditOperation(null)
      showSuccess('Изменения сохранены')
    } catch (err: any) {
      showError(err)
    }
  }

  const handleDelete = async (id: number) => {
    const operationToDelete = operations.find(op => op.id === id)
    
    const confirmed = window.confirm(
      `Вы уверены, что хотите удалить операцию "${operationToDelete?.title}"?`
    )
    
    if (!confirmed) return

    try {
      await operationApi.deleteOperation(id)
      setOperations(operations.filter((op) => op.id !== id))
      showSuccess('Операция удалена')
    } catch (err: any) {
      showError(err)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Загрузка операций...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Операции</h2>
        <div className="flex items-center gap-4">
          <span className="text-gray-600">Всего: {operations.length}</span>
          <input
            type="text"
            placeholder="Поиск по названию..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 w-64"
          />
          <button
            onClick={() => setShowCreateForm(true)}
            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 transition"
          >
            Добавить операцию
          </button>
        </div>
      </div>

      <Table
        columns={columns}
        data={filtered}
      />

      {/* Модальное окно создания */}
      <Modal isOpen={showCreateForm} onClose={() => setShowCreateForm(false)}>
        <div className="p-4">
          <CreateOperationForm 
            onSubmit={handleCreate}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      </Modal>

      {/* Модальное окно редактирования */}
      <Modal isOpen={!!editOperation} onClose={() => setEditOperation(null)}>
        <div className="p-4">
          {editOperation && (
            <EditOperationForm
              operation={editOperation}
              onSubmit={(data) => handleEditSubmit({ id: editOperation.id, title: data.title })}
              onCancel={() => setEditOperation(null)}
            />
          )}
        </div>
      </Modal>

      {/* Универсальное модальное окно уведомлений */}
      <NotificationModal
        isOpen={notification.isOpen}
        onClose={closeNotification}
        type={notification.type}
        title={notification.title}
        message={notification.message}
        details={notification.details}
        onAction={notification.onAction}
        actionText={notification.actionText}
      />
    </div>
  )
}