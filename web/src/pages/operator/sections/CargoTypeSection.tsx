import { useEffect, useState } from 'react'
import { Table } from '../../../components/ui/Table'
import type { TableColumn } from '../../../components/ui/Table'
import Modal from '../../../components/ui/Modal'
import NotificationModal, {
  type NotificationType,
} from '../../../components/ui/NotificationModal'
import { cargoTypeApi, type CargoType } from '../../../api/cargoTypeAPI'
import CreateCargoType from '../forms/cargotype/CreateCargoTypeForm'
import EditCargoType from '../forms/cargotype/EditCargoTypeForm'

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
        message: 'Такой тип груза уже существует',
        details: 'Измените название',
      }

    case 'NotFound':
      return {
        title: 'Не найдено',
        message: 'Тип груза не найден',
      }

    case 'InvalidArgument':
      return {
        title: 'Некорректные данные',
        message: error.message,
        details: 'Проверьте введённые значения',
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

export default function CargoTypeSection() {
  const [items, setItems] = useState<CargoType[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [editItem, setEditItem] = useState<CargoType | null>(null)
  const [notification, setNotification] = useState<NotificationState>({
    isOpen: false,
    type: 'info',
    title: '',
    message: '',
  })

  const showError = (error: any, onRetry?: () => void) => {
    const f = formatErrorMessage(error)
    setNotification({
      isOpen: true,
      type: 'error',
      title: f.title,
      message: f.message,
      details: f.details,
      onAction: onRetry,
      actionText: onRetry ? 'Повторить' : undefined,
    })
  }

  const showSuccess = (message: string) => {
    setNotification({
      isOpen: true,
      type: 'success',
      title: 'Успешно',
      message,
    })
  }

  const fetchItems = async () => {
    try {
      setLoading(true)
      setItems(await cargoTypeApi.list())
    } catch (e: any) {
      showError(e, fetchItems)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchItems()
  }, [])

  const filtered = items.filter((i) =>
    i.title.toLowerCase().includes(search.toLowerCase())
  )

  const columns: TableColumn<CargoType>[] = [
    { key: 'id', title: 'ID' },
    { key: 'title', title: 'Название' },
    {
      key: 'processCost',
      title: 'Стоимость обработки',
      render: (v) => `${v} ₽`,
    },
    {
      key: 'actions',
      title: 'Действия',
      render: (_, item) => (
        <div className="flex gap-2">
          <button
            className="px-2 py-1 bg-blue-500 text-white rounded"
            onClick={() => setEditItem(item)}
          >
            Редактировать
          </button>
          <button
            className="px-2 py-1 bg-red-500 text-white rounded"
            onClick={() => handleDelete(item.id)}
          >
            Удалить
          </button>
        </div>
      ),
    },
  ]

  const handleCreate = async (data: Omit<CargoType, 'id'>) => {
    try {
      const id = await cargoTypeApi.create(data)
      setItems([...items, { id, ...data }])
      setShowCreate(false)
      showSuccess('Тип груза создан')
    } catch (e: any) {
      showError(e)
    }
  }

  const handleUpdate = async (data: CargoType) => {
    try {
      await cargoTypeApi.update(data.id, data)
      setItems(items.map((i) => (i.id === data.id ? data : i)))
      setEditItem(null)
      showSuccess('Изменения сохранены')
    } catch (e: any) {
      showError(e)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить тип груза?')) return

    try {
      await cargoTypeApi.delete(id)
      setItems(items.filter((i) => i.id !== id))
      showSuccess('Тип груза удалён')
    } catch (e: any) {
      showError(e)
    }
  }

  if (loading) {
    return <div>Загрузка...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Типы грузов</h2>
        <div className="flex items-center gap-4">
          <span className="text-gray-600">Всего: {items.length}</span>
          <input
            type="text"
            placeholder="Поиск по названию..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 transition"
            onClick={() => setShowCreate(true)}
          >
            Добавить тип груза
          </button>
        </div>
      </div>

      <Table columns={columns} data={filtered} />

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)}>
        <CreateCargoType onSubmit={handleCreate} onCancel={() => setShowCreate(false)} />
      </Modal>

      <Modal isOpen={!!editItem} onClose={() => setEditItem(null)}>
        {editItem && (
          <EditCargoType
            cargoType={editItem}
            onSubmit={handleUpdate}
            onCancel={() => setEditItem(null)}
          />
        )}
      </Modal>

      <NotificationModal
        {...notification}
        onClose={() => setNotification({ ...notification, isOpen: false })}
      />
    </div>
  )
}
