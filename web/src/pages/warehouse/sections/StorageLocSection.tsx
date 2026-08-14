import { useEffect, useState } from 'react'
import { Table } from '../../../components/ui/Table'
import type { TableColumn } from '../../../components/ui/Table'
import Modal from '../../../components/ui/Modal'
import NotificationModal, { type NotificationType } from '../../../components/ui/NotificationModal'

import CreateStorageLocForm from '../forms/storageloc/CreateStorageLocForm'
import EditStorageLocForm from '../forms/storageloc/EditStorageLocForm'
import UseStorageLocForm from '../forms/storageloc/UseStorageLocForm'
import ResetStorageLocForm from '../forms/storageloc/ResetStorageLocForm'

import { storageLocAPI, type StorageLocType, type StorageLocCreate, type StorageLocUpdate, type StorageLocUse } from '../../../api/storageLocAPI'

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
        message: 'Такое место хранения уже существует',
      }

    case 'NotFound':
      if (error.message?.includes('cargo')) {
        return {
          title: 'Не найдено',
          message: 'Груз с таким ID не найден',
        }
      }
      if (error.message?.includes('related')) {
        return {
          title: 'Не найдено',
          message: 'Тип груза с таким ID не найден',
        }
      }
      return {
        title: 'Не найдено',
        message: 'Место хранения не найдено',
      }

    case 'InvalidArgument':
      if (error.message?.includes('cargoTypeId')) {
        return {
          title: 'Некорректные данные',
          message: 'Указан несуществующий ID типа груза',
        }
      }
      if (error.message?.includes('cargo_id')) {
        return {
          title: 'Некорректные данные',
          message: 'Указан несуществующий ID груза',
        }
      }
      if (error.message?.includes('date_of_placement')) {
        return {
          title: 'Некорректные данные',
          message: 'Указана дата в будующем',
        }
      }
      if (error.message?.includes('maxWeight')) {
        return {
          title: 'Некорректные данные',
          message: 'Максимальный вес должен быть положительным числом',
        }
      }
      if (error.message?.includes('maxVolume')) {
        return {
          title: 'Некорректные данные',
          message: 'Максимальный объем должен быть положительным числом',
        }
      }
      return {
        title: 'Некорректные данные',
        message: error.message,
      }

    case 'FailedPrecondition':
      if (error.message?.includes('occupied')) {
        return {
          title: 'Место занято',
          message: 'Место хранения уже занято другим грузом',
        }
      }
      if (error.message?.includes('placed')) {
        return {
          title: 'Груз уже размещен',
          message: 'Груз с таким ID уже размещен на складе',
        }
      }
      if (error.message?.includes('is used')) {
        return {
          title: 'Место занято',
          message: 'Сбросьте груз перед удалением',
        }
      }
      if (error.message?.includes('free')) {
        return {
          title: 'Место свободно',
          message: 'Место хранения уже свободно',
        }
      }
      if (error.message?.includes('type not suitable')) {
        return {
          title: 'Груз не поддерживается',
          message: 'Указанный груз не соответствует типу места хранения',
        }
      }
      if (error.message?.includes('not suitable')) {
        return {
          title: 'Груз не поддерживается',
          message: 'Указанный груз превышает габариты места хранения',
        }
      }
      if (error.message?.includes('capacity')) {
        return {
          title: 'Превышена вместимость',
          message: 'Груз не помещается в место хранения',
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

const checkDate = (value: any) => {
  if (value === "") {
    return "-"
  }

  return value
}

export default function StorageLocSection() {
  const [search, setSearch] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [editStorageLoc, setEditStorageLoc] = useState<StorageLocType | null>(null)
  const [useStorageLoc, setUseStorageLoc] = useState<StorageLocType | null>(null)
  const [resetStorageLoc, setResetStorageLoc] = useState<StorageLocType | null>(null)
  const [storageLocs, setStorageLocs] = useState<StorageLocType[]>([])
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
    fetchStorageLocs()
  }, [])

  const fetchStorageLocs = async () => {
    try {
      setLoading(true)
      const data = await storageLocAPI.list()
      setStorageLocs(Array.isArray(data) ? data : [])
    } catch (err: any) {
      showError(err, fetchStorageLocs)
    } finally {
      setLoading(false)
    }
  }

  const filtered = storageLocs.filter((loc) =>
    loc.id.toString().includes(search) ||
    loc.cargoTypeId.toString().includes(search) ||
    (loc.cargoId && loc.cargoId.toString().includes(search))
  )

  // ==================== TABLE ====================
  const columns: TableColumn<StorageLocType>[] = [
    { key: 'id', title: 'ID' },
    { key: 'cargoTypeId', title: 'Тип груза' },
    { 
      key: 'maxWeight', 
      title: 'Макс. вес', 
      render: (v) => `${v} т` 
    },
    { 
      key: 'maxVolume', 
      title: 'Макс. объем', 
      render: (v) => `${v} м³` 
    },
    { 
      key: 'cargoId', 
      title: 'Статус', 
      render: (v) => v ? `Занято (ID: ${v})` : 'Свободно'
    },
    { 
      key: 'dateOfPlacement', 
      title: 'Дата размещения', 
      render: (v) => checkDate(v)
    },
    {
      key: 'actions',
      title: 'Действия',
      render: (_: any, loc: StorageLocType) => (
        <div className="flex flex-wrap gap-2">
          {!loc.cargoId ? (
            // Две кнопки для свободного места
            <>
              <button
                className="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                onClick={() => setEditStorageLoc(loc)}
              >
                Редактировать
              </button>
              <button
                className="px-2 py-1 bg-green-500 text-white rounded hover:bg-green-600"
                onClick={() => setUseStorageLoc(loc)}
              >
                Разместить груз
              </button>
            </>
          ) : (
            // Две кнопки для занятого места
            
              
              <button
                className="px-2 py-1 bg-orange-500 text-white rounded hover:bg-orange-600"
                onClick={() => setResetStorageLoc(loc)}
              >
                Сбросить
              </button>
            
          )}
          <button
            className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
            onClick={() => handleDelete(loc.id)}
          >
            Удалить
          </button>
        </div>
      ),
    },
  ]

  // ==================== CREATE ====================
  const handleCreate = async (data: StorageLocCreate) => {
    try {
      await storageLocAPI.create(data)
      await fetchStorageLocs() // Перезагружаем список
      setShowCreateForm(false)
      showSuccess('Место хранения успешно создано')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== UPDATE ====================
  // ==================== UPDATE ====================
  const handleEditSubmit = async (data: StorageLocUpdate) => {
    if (!editStorageLoc) return

    // Проверяем, что есть хотя бы одно поле для обновления
    if (Object.keys(data).length === 0) {
      alert('Нет изменений для сохранения')
      return
    }

    try {
      await storageLocAPI.update(editStorageLoc.id, data)
      await fetchStorageLocs() // Перезагружаем список
      setEditStorageLoc(null)
      showSuccess('Изменения сохранены')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== USE ====================
  const handleUseSubmit = async (data: StorageLocUse) => {
    if (!useStorageLoc) return

    try {
      await storageLocAPI.use(useStorageLoc.id, data)
      await fetchStorageLocs() // Перезагружаем список
      setUseStorageLoc(null)
      showSuccess('Груз успешно размещен в месте хранения')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== RESET ====================
  const handleResetSubmit = async () => {
    if (!resetStorageLoc) return

    try {
      await storageLocAPI.reset(resetStorageLoc.id)
      await fetchStorageLocs() // Перезагружаем список
      setResetStorageLoc(null)
      showSuccess('Место хранения успешно сброшено')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== DELETE ====================
  const handleDelete = async (id: number) => {

    if (!window.confirm(`Удалить место хранения ID: ${id}?`)) return

    try {
      await storageLocAPI.delete(id)
      setStorageLocs(storageLocs.filter((loc) => loc.id !== id))
      showSuccess('Место хранения удалено')
    } catch (err: any) {
      showError(err)
    }
  }

  // ==================== RENDER ====================
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Загрузка мест хранения...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <h2 className="text-2xl font-bold">Места хранения</h2>

        <div className="flex flex-col md:flex-row items-start md:items-center gap-4 w-full md:w-auto">
          <div className="flex items-center gap-2">
            <span className="text-gray-600 whitespace-nowrap">
              Всего: {storageLocs.length} 
              ({storageLocs.filter(l => l.cargoId).length} занято)
            </span>
          </div>

          <input
            type="text"
            placeholder="Поиск по ID или ID груза..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3 py-2 border rounded-md w-full md:w-auto"
          />

          <button
            onClick={() => setShowCreateForm(true)}
            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 whitespace-nowrap"
          >
            Добавить место хранения
          </button>
        </div>
      </div>

      <Table columns={columns} data={filtered} />

      {/* CREATE MODAL */}
      <Modal isOpen={showCreateForm} onClose={() => setShowCreateForm(false)}>
        <div className="p-4">
          <CreateStorageLocForm
            onSubmit={handleCreate}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      </Modal>

      {/* EDIT MODAL */}
      <Modal isOpen={!!editStorageLoc} onClose={() => setEditStorageLoc(null)}>
        <div className="p-4">
          {editStorageLoc && (
            <EditStorageLocForm
              storageLoc={editStorageLoc}
              onSubmit={handleEditSubmit}
              onCancel={() => setEditStorageLoc(null)}
            />
          )}
        </div>
      </Modal>

      {/* USE MODAL */}
      <Modal isOpen={!!useStorageLoc} onClose={() => setUseStorageLoc(null)}>
        <div className="p-4">
          {useStorageLoc && (
            <UseStorageLocForm
              onSubmit={handleUseSubmit}
              onCancel={() => setUseStorageLoc(null)}
            />
          )}
        </div>
      </Modal>

      {/* RESET MODAL */}
      <Modal isOpen={!!resetStorageLoc} onClose={() => setResetStorageLoc(null)}>
        <div className="p-4">
          {resetStorageLoc && (
            <ResetStorageLocForm
              storageLoc={resetStorageLoc}
              onSubmit={handleResetSubmit}
              onCancel={() => setResetStorageLoc(null)}
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