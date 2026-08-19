import { useEffect, useState } from 'react'
import { Table, type TableColumn } from '../../../components/ui/Table'
import Modal from '../../../components/ui/Modal'
import NotificationModal, {
  type NotificationType,
} from '../../../components/ui/NotificationModal'
import CreateOperCargoForm from '../forms/opercargo/CreateOperCargoForm'
import {
  operCargoApi,
  type OperCargoType,
} from '../../../api/operCargoAPI'

/**
 * UI-модель для таблицы
 * id нужен ТОЛЬКО для Table
 */
type OperCargoRow = {
  id: string
  operationId: number
  cargoId: number
}

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
        message: 'Связь уже существует',
        details: 'Этот груз уже добавлен в операцию',
      }

    case 'NotFound':
      return {
        title: 'Не найдено',
        message: 'Операция или груз не найдены',
      }

    case 'InvalidArgument':
      return {
        title: 'Некорректные данные',
        message: error.message,
      }

    case 'FailedPrecondition':
      return {
        title: 'Невозможно добавить',
        message: 'Такого груза или операции не существует',
      }

    default:
      return {
        title: 'Ошибка',
        message: error.message || 'Неизвестная ошибка',
      }
  }
}


export default function OperationCargoSection() {
  const [items, setItems] = useState<OperCargoRow[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)

  const [notification, setNotification] = useState<NotificationState>({
    isOpen: false,
    type: 'info',
    title: '',
    message: '',
  })

  const closeNotification = () =>
    setNotification((p) => ({ ...p, isOpen: false }))

  const showError = (error: any, retry?: () => void) => {
  const formatted = formatErrorMessage(error)

  setNotification({
    isOpen: true,
    type: 'error',
    title: formatted.title,
    message: formatted.message,
    details: formatted.details,
    onAction: retry,
    actionText: retry ? 'Повторить' : undefined,
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

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      setLoading(true)

      const data = await operCargoApi.list()

      const rows: OperCargoRow[] = data.map((x) => ({
        id: `${x.operationId}-${x.cargoId}`,
        operationId: x.operationId,
        cargoId: x.cargoId,
      }))

      setItems(rows)
    } catch (e: any) {
      showError(e, fetchData)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (data: OperCargoType) => {
    try {
      await operCargoApi.create(data)

      setItems([
        ...items,
        {
          id: `${data.operationId}-${data.cargoId}`,
          operationId: data.operationId,
          cargoId: data.cargoId,
        },
      ])

      setShowCreate(false)
      showSuccess('Груз добавлен в операцию')
    } catch (e: any) {
      showError(e)
    }
  }

  const handleDelete = async (row: OperCargoRow) => {
    const confirmed = window.confirm(
      `Удалить груз ${row.cargoId} из операции ${row.operationId}?`
    )
    if (!confirmed) return

    try {
      await operCargoApi.delete(row.operationId, row.cargoId)
      setItems(items.filter((i) => i.id !== row.id))
      showSuccess('Связь удалена')
    } catch (e: any) {
      showError(e)
    }
  }

  const columns: TableColumn<OperCargoRow>[] = [
    { key: 'id', title: 'ID' },
    { key: 'operationId', title: 'Операция ID' },
    { key: 'cargoId', title: 'Груз ID' },
    {
      key: 'actions',
      title: 'Действия',
      render: (_, row) => (
        <button
          className="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600"
          onClick={() => handleDelete(row)}
        >
          Удалить
        </button>
      ),
    },
  ]

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
        <h2 className="text-2xl font-bold">Грузы в операциях</h2>
        <button
          onClick={() => setShowCreate(true)}
          className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        >
          Добавить связь
        </button>
      </div>

      <Table columns={columns} data={items} />

      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)}>
        <div className="p-4">
          <CreateOperCargoForm
            onSubmit={handleCreate}
            onCancel={() => setShowCreate(false)}
          />
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
