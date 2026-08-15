import { useEffect, useState } from 'react'
import { Table } from '../../../components/ui/Table'
import type { TableColumn } from '../../../components/ui/Table'
import NotificationModal, { type NotificationType } from '../../../components/ui/NotificationModal'
import { reportAPI, type CargoTypeReportItem } from '../../../api/reportAPI'

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
    case 'NotFound':
      return {
        title: 'Не найдено',
        message: 'Данные для отчета не найдены',
      }
    default:
      return {
        title: 'Ошибка',
        message: error.message || 'Неизвестная ошибка',
      }
  }
}

export default function CargoTypeReportSection() {
  const [data, setData] = useState<CargoTypeReportItem[]>([])
  const [loading, setLoading] = useState(true)
  const [exporting, setExporting] = useState(false)

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

  // ==================== LOAD DATA ====================
  useEffect(() => {
    loadReport()
  }, [])

  const loadReport = async () => {
    try {
      setLoading(true)
      const reportData = await reportAPI.getCargoTypeReport()
      console.log('Loaded cargo type data:', reportData)
      setData(Array.isArray(reportData) ? reportData : [])
    } catch (err: any) {
      console.error('Load error:', err)
      showError(err, loadReport)
    } finally {
      setLoading(false)
    }
  }

  // ==================== EXPORT ====================
  const exportToCSV = async () => {
    try {
      setExporting(true)
      
      const csvContent = generateCSV(data)
      
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
      const link = document.createElement('a')
      const url = URL.createObjectURL(blob)
      
      link.setAttribute('href', url)
      link.setAttribute('download', `отчет_типы_грузов_${new Date().toISOString().split('T')[0]}.csv`)
      link.style.visibility = 'hidden'
      
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      showSuccess('Отчет успешно экспортирован в CSV')
    } catch (err: any) {
      showError(err)
    } finally {
      setExporting(false)
    }
  }

  const generateCSV = (items: CargoTypeReportItem[]): string => {
    if (!items.length) return ''
    
    const headers = ['Тип груза', 'Количество грузов', 'Суммарный вес, т', 'Стоимость обработки']
    const csvRows = []
    
    csvRows.push(headers.join(';'))
    
    items.forEach(item => {
      const row = [
        item.cargoTypeName,
        item.count,
        item.weight.toFixed(1),
        item.processCost || '—'
      ]
      csvRows.push(row.join(';'))
    })
    
    return csvRows.join('\n')
  }

  // ==================== TABLE COLUMNS ====================
  const columns: TableColumn<CargoTypeReportItem>[] = [
    { key: 'cargoTypeName', title: 'Тип груза' },
    { 
      key: 'count', 
      title: 'Количество грузов', 
      render: (v: string | number | undefined) => {
        const num = Number(v)
        return !isNaN(num) ? num.toString() : '—'
      }
    },
    { 
      key: 'weight', 
      title: 'Суммарный вес, т', 
      render: (v: string | number | undefined) => {
        const num = Number(v)
        return !isNaN(num) ? num.toFixed(1) : '—'
      }
    },
    { 
      key: 'processCost', 
      title: 'Стоимость обработки', 
      render: (v: string | number | undefined) => {
        if (!v) return '—'
        return String(v)
      }
    },
  ]

  // ==================== CALCULATE TOTALS ====================
  const calculateTotals = () => {
    const totalCargos = data.reduce((sum, item) => sum + item.count, 0)
    const totalWeight = data.reduce((sum, item) => sum + item.weight, 0)
    
    // Извлекаем числа из строки стоимости (например: "210 000 руб." -> 210000)
    // Просто суммируй числа, не парси строки!
    const totalCost = data.reduce((sum, item) => {
        return sum + (item.processCost || 0)
    }, 0)
    
    return {
      totalCargoTypes: data.length,
      totalCargos,
      totalWeight: totalWeight.toFixed(1),
      totalCost: totalCost.toLocaleString('ru-RU') + ' руб.',
      avgCostPerTon: totalWeight > 0 ? (totalCost / totalWeight).toFixed(0) : '0'
    }
  }

  const totals = calculateTotals()

  // ==================== RENDER ====================
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Загрузка отчета...</div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <h2 className="text-2xl font-bold">Отчет по типам грузов</h2>

        <div className="flex flex-col md:flex-row items-start md:items-center gap-4 w-full md:w-auto">
          <div className="flex items-center gap-4">
            <button
              onClick={loadReport}
              className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 whitespace-nowrap"
            >
              Обновить отчет
            </button>
            
            <button
              onClick={exportToCSV}
              disabled={exporting || data.length === 0}
              className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
            >
              {exporting ? 'Экспорт...' : 'Экспорт в CSV'}
            </button>
          </div>
        </div>
      </div>

      {/* Summary info */}
      <div className="mb-6 p-4 bg-gray-50 rounded-md">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-sm text-gray-500">Типов грузов</div>
            <div className="text-xl font-bold">{totals.totalCargoTypes}</div>
          </div>
          <div className="text-center">
            <div className="text-sm text-gray-500">Всего грузов</div>
            <div className="text-xl font-bold">{totals.totalCargos}</div>
          </div>
          <div className="text-center">
            <div className="text-sm text-gray-500">Общий вес</div>
            <div className="text-xl font-bold">{totals.totalWeight} т</div>
          </div>
          <div className="text-center">
            <div className="text-sm text-gray-500">Общая стоимость</div>
            <div className="text-xl font-bold">{totals.totalCost}</div>
          </div>
        </div>
        <div className="mt-4 pt-4 border-t text-center">
          <div className="text-sm text-gray-500">Средняя стоимость обработки</div>
          <div className="text-lg font-semibold">{totals.avgCostPerTon} руб./т</div>
        </div>
      </div>

      {/* Table */}
      <div className="mb-8">
        {data.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            Нет данных для отображения
          </div>
        ) : (
          <Table columns={columns} data={data} />
        )}
      </div>

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