import { useState } from 'react'
import { type StorageLocCreate } from '../../../../api/storageLocAPI'

type CreateStorageLocFormProps = {
  onSubmit?: (data: StorageLocCreate) => void
  onCancel?: () => void
}

export default function CreateStorageLocForm({ onSubmit, onCancel }: CreateStorageLocFormProps) {
  const [cargoTypeId, setCargoTypeId] = useState<number | ''>('')
  const [maxWeight, setMaxWeight] = useState<number | ''>('')
  const [maxVolume, setMaxVolume] = useState<number | ''>('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (cargoTypeId === '' || maxWeight === '' || maxVolume === '') {
      alert('Заполните все поля')
      return
    }

    if (cargoTypeId <= 0) {
      alert('ID типа груза должен быть больше 0')
      return
    }

    if (Number(maxWeight) <= 0) {
      alert('Максимальный вес должен быть больше 0')
      return
    }

    if (Number(maxVolume) <= 0) {
      alert('Максимальный объем должен быть больше 0')
      return
    }

    const data: StorageLocCreate = {
      cargoTypeId: Number(cargoTypeId),
      maxWeight: Number(maxWeight),
      maxVolume: Number(maxVolume),
    }

    if (onSubmit) onSubmit(data)
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Создать новое место хранения</h2>

      <div>
        <label className="block mb-1">ID типа груза</label>
        <input
          type="number"
          value={cargoTypeId}
          onChange={(e) => setCargoTypeId(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
          min="1"
        />
        <p className="text-sm text-gray-500 mt-1">Обязательное поле, должно быть больше 0</p>
      </div>

      <div>
        <label className="block mb-1">Максимальный вес</label>
        <input
          type="number"
          value={maxWeight}
          onChange={(e) => setMaxWeight(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
          min="0.01"
          step="0.01"
        />
        <p className="text-sm text-gray-500 mt-1">В килограммах, должно быть больше 0</p>
      </div>

      <div>
        <label className="block mb-1">Максимальный объем</label>
        <input
          type="number"
          value={maxVolume}
          onChange={(e) => setMaxVolume(e.target.valueAsNumber)}
          className="w-full border rounded px-2 py-1"
          required
          min="0.01"
          step="0.01"
        />
        <p className="text-sm text-gray-500 mt-1">В кубических метрах, должно быть больше 0</p>
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