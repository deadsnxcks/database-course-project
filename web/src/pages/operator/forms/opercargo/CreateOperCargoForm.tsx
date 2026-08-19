import { useState } from 'react'

type CreateOperCargoFormProps = {
  onSubmit?: (data: {
    operationId: number
    cargoId: number
  }) => void
  onCancel?: () => void
}

export default function CreateOperCargoForm({
  onSubmit,
  onCancel,
}: CreateOperCargoFormProps) {
  const [operationId, setOperationId] = useState<number | ''>('')
  const [cargoId, setCargoId] = useState<number | ''>('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!operationId || !cargoId) {
      alert('Заполните все поля')
      return
    }

    onSubmit?.({
      operationId: Number(operationId),
      cargoId: Number(cargoId),
    })
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">
        Добавить груз в операцию
      </h2>

      <div>
        <label className="block mb-1">ID операции</label>
        <input
          type="number"
          value={operationId}
          onChange={(e) => setOperationId(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          min="1"
          required
        />
      </div>

      <div>
        <label className="block mb-1">ID груза</label>
        <input
          type="number"
          value={cargoId}
          onChange={(e) => setCargoId(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          min="1"
          required
        />
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Добавить
        </button>
        <button
          type="button"
          className="bg-gray-300 px-4 py-2 rounded hover:bg-gray-400"
          onClick={onCancel}
        >
          Отмена
        </button>
      </div>
    </form>
  )
}
