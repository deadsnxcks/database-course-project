import { useState } from 'react'

type CreateVesselProps = {
  onSubmit?: (data: { title: string; vesselType: string; maxLoad: number }) => void
  onCancel?: () => void
}

export default function CreateVessel({ onSubmit, onCancel }: CreateVesselProps) {
  const [title, setTitle] = useState('')
  const [vesselType, setVesselType] = useState('')
  const [maxLoad, setMaxLoad] = useState<number | ''>('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!title || !vesselType || !maxLoad) {
      alert('Заполните все поля')
      return
    }

    const data = { title, vesselType, maxLoad: Number(maxLoad) }
    if (onSubmit) onSubmit(data)
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Создать новое судно</h2>
      
      <div>
        <label className="block mb-1">Название</label>
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="w-full border rounded px-2 py-1"
          required
        />
      </div>

      <div>
        <label className="block mb-1">Тип судна</label>
        <input
          type="text"
          value={vesselType}
          onChange={(e) => setVesselType(e.target.value)}
          className="w-full border rounded px-2 py-1"
          required
        />
      </div>

      <div>
        <label className="block mb-1">Максимальная загрузка</label>
        <input
          type="number"
          value={maxLoad}
          onChange={(e) => setMaxLoad(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
          min="0"
          step="0.01"
        />
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Создать
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