import { useState, useEffect } from 'react'
import { type StorageLocType, type StorageLocUpdate } from '../../../../api/storageLocAPI'

type EditStorageLocFormProps = {
  storageLoc: StorageLocType
  onSubmit?: (data: StorageLocUpdate) => void
  onCancel?: () => void
}

export default function EditStorageLocForm({
  storageLoc,
  onSubmit,
  onCancel,
}: EditStorageLocFormProps) {
  const [cargoTypeId, setCargoTypeId] = useState<number | ''>('')
  const [maxWeight, setMaxWeight] = useState<number | ''>('')
  const [maxVolume, setMaxVolume] = useState<number | ''>('')

  // Инициализация
  useEffect(() => {
    setCargoTypeId(storageLoc.cargoTypeId)
    setMaxWeight(storageLoc.maxWeight)
    setMaxVolume(storageLoc.maxVolume)
  }, [storageLoc])

  // ==================== HELPERS ====================
  const isValidNumber = (v: number | ''): v is number =>
    v !== '' && Number.isFinite(v)

  // ==================== SUBMIT ====================
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    const isCargoTypeChanged =
      isValidNumber(cargoTypeId) && cargoTypeId !== storageLoc.cargoTypeId

    const isWeightChanged =
      isValidNumber(maxWeight) && maxWeight !== storageLoc.maxWeight

    const isVolumeChanged =
      isValidNumber(maxVolume) && maxVolume !== storageLoc.maxVolume

    if (!isCargoTypeChanged && !isWeightChanged && !isVolumeChanged) {
      alert('Нет изменений для сохранения')
      return
    }

    // Валидация
    if (isCargoTypeChanged && cargoTypeId <= 0) {
      alert('ID типа груза должен быть больше 0')
      return
    }

    if (isWeightChanged && maxWeight <= 0) {
      alert('Максимальный вес должен быть больше 0')
      return
    }

    if (isVolumeChanged && maxVolume <= 0) {
      alert('Максимальный объем должен быть больше 0')
      return
    }

    const data: StorageLocUpdate = {}

    if (isCargoTypeChanged) data.cargoTypeId = cargoTypeId
    if (isWeightChanged) data.maxWeight = maxWeight
    if (isVolumeChanged) data.maxVolume = maxVolume

    onSubmit?.(data)
  }

  // ==================== RESET ====================
  const handleReset = () => {
    setCargoTypeId(storageLoc.cargoTypeId)
    setMaxWeight(storageLoc.maxWeight)
    setMaxVolume(storageLoc.maxVolume)
  }

  // ==================== INPUT HANDLERS ====================
  const handleNumberChange =
    (setter: (v: number | '') => void) =>
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const value = e.target.value
      setter(value === '' ? '' : Number(value))
    }

  // ==================== RENDER ====================
  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Редактировать место хранения</h2>

      <div className="mb-2">
        <p className="text-sm text-gray-600">ID: {storageLoc.id}</p>
        {storageLoc.cargoId && (
          <p className="text-sm text-gray-600">
            Занято грузом ID: {storageLoc.cargoId} (размещен:{' '}
            {storageLoc.dateOfPlacement})
          </p>
        )}
      </div>

      <div>
        <label className="block mb-1">ID типа груза</label>
        <input
          type="number"
          value={cargoTypeId}
          onChange={handleNumberChange(setCargoTypeId)}
          className="w-full border rounded px-2 py-1"
          min="1"
        />
        <p className="text-sm text-gray-500 mt-1">
          Текущее значение: {storageLoc.cargoTypeId}
        </p>
      </div>

      <div>
        <label className="block mb-1">Максимальный вес</label>
        <input
          type="number"
          value={maxWeight}
          onChange={handleNumberChange(setMaxWeight)}
          className="w-full border rounded px-2 py-1"
          min="0.01"
          step="0.01"
        />
        <p className="text-sm text-gray-500 mt-1">
          Текущее значение: {storageLoc.maxWeight} кг
        </p>
      </div>

      <div>
        <label className="block mb-1">Максимальный объем</label>
        <input
          type="number"
          value={maxVolume}
          onChange={handleNumberChange(setMaxVolume)}
          className="w-full border rounded px-2 py-1"
          min="0.01"
          step="0.01"
        />
        <p className="text-sm text-gray-500 mt-1">
          Текущее значение: {storageLoc.maxVolume} м³
        </p>
      </div>

      <div className="flex flex-wrap gap-2">
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Сохранить
        </button>

        <button
          type="button"
          className="bg-gray-300 text-black px-4 py-2 rounded hover:bg-gray-400"
          onClick={handleReset}
        >
          Сбросить изменения
        </button>

        <button
          type="button"
          className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
          onClick={onCancel}
        >
          Отмена
        </button>
      </div>
    </form>
  )
}
