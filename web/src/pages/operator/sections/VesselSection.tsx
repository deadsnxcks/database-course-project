import { useState, useEffect } from 'react'
import { Table } from '../../../components/ui/Table'
import type { TableColumn } from '../../../components/ui/Table'
import CreateVessel from '../forms/vessel/CreateVesselForm'
import EditVessel from '../forms/vessel/EditVesselForm'
import Modal from '../../../components/ui/Modal'
import NotificationModal, { type NotificationType } from '../../../components/ui/NotificationModal'
import { vesselApi, type VesselType } from '../../../api/vesselAPI'

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
        message: 'Судно с таким названием уже существует',
      }

    case 'NotFound':
      return {
        title: 'Не найдено',
        message: 'Судно не найдено',
      }

    case 'InvalidArgument':
      return {
        title: 'Некорректные данные',
        message: error.message,
      }

    case 'FailedPrecondition':
      return {
        title: 'Невозможно удалить',
        message: 'Объект используется в других данных',
      }

    default:
      return {
        title: 'Ошибка',
        message: error.message || 'Неизвестная ошибка',
      }
  }
}


export default function VesselSection() {
  const [search, setSearch] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editVessel, setEditVessel] = useState<VesselType | null>(null)
  const [vessels, setVessels] = useState<VesselType[]>([])
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
    fetchVessels()
  }, [])

  const fetchVessels = async () => {
    try {
      setLoading(true)
      const data = await vesselApi.listVessels()
      setVessels(Array.isArray(data) ? data : [])
    } catch (err: any) {
      showError(err, fetchVessels)
    } finally {
      setLoading(false)
    }
  }

  const filtered = (vessels || []).filter((v) =>
    v.title.toLowerCase().includes(search.toLowerCase())
  )

  const columns: TableColumn<VesselType>[] = [
    { key: 'id', title: 'ID' },
    { key: 'title', title: 'Название' },
    { key: 'vesselType', title: 'Тип судна' },
    { key: 'maxLoad', title: 'Макс. загрузка', render: (value) => `${value} т` },
    {
      key: 'actions',
      title: 'Действия',
      render: (_, v) => (
        <div className="flex gap-2">
          <button
            className="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
            onClick={() => setEditVessel(v)}
          >
            Редактировать
          </button>
          <button
            className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition"
            onClick={() => handleDelete(v.id)}
          >
            Удалить
          </button>
        </div>
      ),
    },
  ]

  // ==================== CREATE ====================
  const handleCreate = async (data: Omit<VesselType, 'id'>) => {
    try {
      const id = await vesselApi.createVessel(data)
      setVessels([...vessels, { id, ...data }])
      setShowCreateForm(false)
      showSuccess('Судно успешно создано')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== UPDATE ====================
  const handleEditSubmit = async (data: VesselType) => {
    try {
      await vesselApi.updateVessel(data.id, {
        title: data.title,
        vesselType: data.vesselType,
        maxLoad: data.maxLoad,
      })
      setVessels(
        vessels.map((v) => (v.id === data.id ? { ...v, ...data } : v))
      )
      setEditVessel(null)
      showSuccess('Изменения сохранены')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== DELETE ====================
  const handleDelete = async (id: number) => {
    const vesselToDelete = vessels.find(v => v.id === id)
    
    const confirmed = window.confirm(
      `Вы уверены, что хотите удалить судно "${vesselToDelete?.title}"?`
    )
    
    if (!confirmed) return

    try {
      await vesselApi.deleteVessel(id)
      setVessels(vessels.filter((item) => item.id !== id))
      showSuccess('Судно удалено')
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
        <h2 className="text-2xl font-bold">Суда</h2>
        <div className="flex items-center gap-4">
          <span className="text-gray-600">Всего: {vessels.length}</span>
          <input
            type="text"
            placeholder="Поиск по названию..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            onClick={() => setShowCreateForm(true)}
            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 transition"
          >
            Добавить судно
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
          <CreateVessel 
            onSubmit={handleCreate}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      </Modal>

      {/* Модальное окно редактирования */}
      <Modal isOpen={!!editVessel} onClose={() => setEditVessel(null)}>
        <div className="p-4">
          {editVessel && (
            <EditVessel
              vessel={editVessel}
              onSubmit={handleEditSubmit}
              onCancel={() => setEditVessel(null)}
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