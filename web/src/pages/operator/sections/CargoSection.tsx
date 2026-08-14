import { useEffect, useState } from 'react'
import { Table } from '../../../components/ui/Table'
import type { TableColumn } from '../../../components/ui/Table'
import Modal from '../../../components/ui/Modal'
import NotificationModal, { type NotificationType } from '../../../components/ui/NotificationModal'

import CreateCargoForm from '../forms/cargo/CreateCargoForm'
import EditCargoForm from '../forms/cargo/EditCargoForm'

import { cargoAPI, type Cargo } from '../../../api/cargoAPI'

interface NotificationState {
  isOpen: boolean
  type: NotificationType
  title: string
  message: string
  details?: string
  onAction?: () => void
  actionText?: string
}

const formatErrorMessage = (error: any) => {
  switch (error.code) {
    case 'AlreadyExists':
      return {
        title: 'Конфликт',
        message: 'Такой груз уже существует',
      }

    case 'NotFound':
      return {
        title: 'Не найдено',
        message: 'Груз не найден',
      }

    case 'InvalidArgument':
      return {
        title: 'Некорректные данные',
      }

    case 'FailedPrecondition':
      if (error.message?.includes('related')) {
        return {
          title: 'Некорректные данные',
          message: 'Указаны несуществующие связанные сущности',
        }
      }

      if (error.message?.includes('used')) {
        return {
          title: 'Невозможно удалить',
          message: 'Груз используется в других данных',

        }
      }

      return {
        title: 'Невозможно выполнить операцию',
        message: error.message,
      }


    default:
      return {
        title: 'Ошибка',
        message: error.message || 'Неизвестная ошибка',
      }
  }
}

export default function CargoSection() {
  const [search, setSearch] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editCargo, setEditCargo] = useState<Cargo | null>(null)
  const [cargos, setCargos] = useState<Cargo[]>([])
  const [loading, setLoading] = useState(true)

  const [notification, setNotification] = useState<NotificationState>({
    isOpen: false,
    type: 'info',
    title: '',
    message: '',
  })

  const closeNotification = () => {
    setNotification((prev) => ({ ...prev, isOpen: false }))
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

  const showSuccess = (message: string, title = 'Успешно') => {
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

  // ==================== LOAD ====================
  useEffect(() => {
    fetchCargos()
  }, [])

  const fetchCargos = async () => {
    try {
      setLoading(true)
      const data = await cargoAPI.list()
      setCargos(Array.isArray(data) ? data : [])
    } catch (err: any) {
      showError(err, fetchCargos)
    } finally {
      setLoading(false)
    }
  }

  const filtered = cargos.filter((c) =>
    c.title.toLowerCase().includes(search.toLowerCase())
  )

  // ==================== TABLE ====================
  const columns: TableColumn<Cargo>[] = [
    { key: 'id', title: 'ID' },
    { key: 'title', title: 'Название' },
    { key: 'typeId', title: 'Тип' },
    { key: 'weight', title: 'Вес', render: (v) => `${v} т` },
    { key: 'volume', title: 'Объём', render: (v) => `${v} м³` },
    { key: 'vesselId', title: 'Судно' },
    {
      key: 'actions',
      title: 'Действия',
      render: (_, c) => (
        <div className="flex gap-2">
          <button
            className="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
            onClick={() => setEditCargo(c)}
          >
            Редактировать
          </button>
          <button
            className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
            onClick={() => handleDelete(c.id)}
          >
            Удалить
          </button>
        </div>
      ),
    },
  ]

  // ==================== CREATE ====================
  const handleCreate = async (data: Omit<Cargo, 'id'>) => {
    try {
      const id = await cargoAPI.create(data)
      setCargos([...cargos, { id, ...data }])
      setShowCreateForm(false)
      showSuccess('Груз успешно создан')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== UPDATE ====================
  const handleEditSubmit = async (data: Cargo) => {
    try {
      await cargoAPI.update(data.id, data)
      setCargos(
        cargos.map((c) => (c.id === data.id ? data : c))
      )
      setEditCargo(null)
      showSuccess('Изменения сохранены')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== DELETE ====================
  const handleDelete = async (id: number) => {
    const cargo = cargos.find((c) => c.id === id)

    if (!window.confirm(`Удалить груз "${cargo?.title}"?`)) return

    try {
      await cargoAPI.delete(id)
      setCargos(cargos.filter((c) => c.id !== id))
      showSuccess('Груз удалён')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== RENDER ====================
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Загрузка...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Грузы</h2>

        <div className="flex items-center gap-4">
          <span className="text-gray-600">Всего: {cargos.length}</span>

          <input
            type="text"
            placeholder="Поиск по названию..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3 py-2 border rounded-md"
          />

          <button
            onClick={() => setShowCreateForm(true)}
            className="px-4 py-2 bg-green-500 text-white rounded-md"
          >
            Добавить груз
          </button>
        </div>
      </div>

      <Table columns={columns} data={filtered} />

      {/* CREATE */}
      <Modal isOpen={showCreateForm} onClose={() => setShowCreateForm(false)}>
        <div className="p-4">
          <CreateCargoForm
            onSubmit={handleCreate}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      </Modal>

      {/* EDIT */}
      <Modal isOpen={!!editCargo} onClose={() => setEditCargo(null)}>
        <div className="p-4">
          {editCargo && (
            <EditCargoForm
              cargo={editCargo}
              onSubmit={handleEditSubmit}
              onCancel={() => setEditCargo(null)}
            />
          )}
        </div>
      </Modal>

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
