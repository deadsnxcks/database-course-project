import { useState } from 'react'

type CreateCargoProps = {
  onSubmit?: (data: { title: string; typeId: number; weight: number; volume: number; vesselId: number }) => void
  onCancel?: () => void
}

export default function CreateCargo({ onSubmit, onCancel }: CreateCargoProps) {
  const [title, setTitle] = useState('')
  const [typeId, setTypeId] = useState<number | ''>('')
  const [weight, setWeight] = useState<number | ''>('')
  const [volume, setVolume] = useState<number | ''>('')
  const [vesselId, setVesselId] = useState<number | ''>('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!title || typeId === '' || weight === '' || volume === '' || vesselId === '') {
      alert('Заполните все поля')
      return
    }

    const data = {
      title,
      typeId: Number(typeId),
      weight: Number(weight),
      volume: Number(volume),
      vesselId: Number(vesselId),
    }

    if (onSubmit) onSubmit(data)
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Создать новый груз</h2>

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
        <label className="block mb-1">ID типа</label>
        <input
          type="number"
          value={typeId}
          onChange={(e) => setTypeId(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
        />
      </div>

      <div>
        <label className="block mb-1">Вес</label>
        <input
          type="number"
          value={weight}
          onChange={(e) => setWeight(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
          min="0"
          step="0.01"
        />
      </div>

      <div>
        <label className="block mb-1">Объём</label>
        <input
          type="number"
          value={volume}
          onChange={(e) => setVolume(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
          min="0"
          step="0.01"
        />
      </div>

      <div>
        <label className="block mb-1">ID судна</label>
        <input
          type="number"
          value={vesselId}
          onChange={(e) => setVesselId(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
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
