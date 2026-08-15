import { useState } from 'react'

type UseStorageLocFormProps = {
  onSubmit?: (data: { cargoId: number; dateOfPlacement: string }) => void
  onCancel?: () => void
}

export default function UseStorageLocForm({ onSubmit, onCancel }: UseStorageLocFormProps) {
  const [cargoId, setCargoId] = useState<number | ''>('')
  const [dateOfPlacement, setDateOfPlacement] = useState<string>('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (cargoId === '' || !dateOfPlacement) {
      alert('Заполните все поля')
      return
    }


    const dateWithSeconds = dateOfPlacement + ':00Z';

    const data = {
      cargoId: Number(cargoId),
      dateOfPlacement: dateWithSeconds,
    }

    if (onSubmit) onSubmit(data)
    else alert(`Используем место хранения: ${JSON.stringify(data)}`)
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Использовать место хранения</h2>

      <div>
        <label className="block mb-1">ID груза</label>
        <input
          type="number"
          value={cargoId}
          onChange={(e) => setCargoId(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
        />
      </div>

      <div>
        <label className="block mb-1">Дата размещения</label>
        <input
          type="datetime-local"
          value={dateOfPlacement}
          onChange={(e) => setDateOfPlacement(e.target.value)}
          className="w-full border rounded px-2 py-1"
        />
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
        >
          Использовать
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
