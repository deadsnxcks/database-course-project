import { useEffect, useState } from 'react'

type EditCargoTypeFormProps = {
  cargoType: {
    id: number
    title: string
    processCost: number
  }
  onSubmit?: (data: {
    id: number
    title: string
    processCost: number
  }) => void
  onCancel?: () => void
}

export default function EditCargoTypeForm({
  cargoType,
  onSubmit,
  onCancel,
}: EditCargoTypeFormProps) {
  const [title, setTitle] = useState('')
  const [processCost, setProcessCost] = useState<number | ''>('')

  // Заполняем форму исходными данными
  useEffect(() => {
    setTitle(cargoType.title)
    setProcessCost(cargoType.processCost)
  }, [cargoType])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!title || processCost === '' || processCost <= 0) {
      alert('Заполните все поля')
      return
    }

    const data = {
      id: cargoType.id,
      title,
      processCost: Number(processCost),
    }

    if (onSubmit) onSubmit(data)
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Редактировать тип груза</h2>

      <div>
        <label className="block mb-1">Название</label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="w-full border rounded px-2 py-1"
        />
      </div>

      <div>
        <label className="block mb-1">Стоимость обработки</label>
        <input
          type="number"
          value={processCost}
          onChange={(e) => setProcessCost(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          min="0"
          step="0.01"
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
