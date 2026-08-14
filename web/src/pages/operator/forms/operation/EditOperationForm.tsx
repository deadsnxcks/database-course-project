import { useState, useEffect } from 'react'

type EditOperationProps = {
  operation: {
    id: number
    title: string
  }
  onSubmit?: (data: { id: number; title: string }) => void
  onCancel?: () => void
}

export default function EditOperationForm({
  operation,
  onSubmit,
  onCancel,
}: EditOperationProps) {
  const [title, setTitle] = useState('')

  // Заполняем форму исходными данными
  useEffect(() => {
    setTitle(operation.title)
  }, [operation])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!title) {
      alert('Заполните название операции')
      return
    }

    const data = { id: operation.id, title }

    if (onSubmit) {
      onSubmit(data)
    } else {
      alert(`Обновляем операцию: ${JSON.stringify(data)}`)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Редактировать операцию</h2>

      <div>
        <label className="block mb-1">Название</label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="w-full border rounded px-2 py-1"
        />
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Сохранить
        </button>
        <button
          type="button"
          className="bg-gray-300 text-black px-4 py-2 rounded hover:bg-gray-400"
          onClick={onCancel}
        >
          Отмена
        </button>
      </div>
    </form>
  )
}
