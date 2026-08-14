import { type StorageLocType } from '../../../../api/storageLocAPI'

type ResetStorageLocFormProps = {
  storageLoc: StorageLocType
  onSubmit?: () => void
  onCancel?: () => void
}

export default function ResetStorageLocForm({ storageLoc, onSubmit, onCancel }: ResetStorageLocFormProps) {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!storageLoc.cargoId) {
      alert('Место хранения уже свободно')
      return
    }

    if (onSubmit) onSubmit()
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2">Сбросить место хранения</h2>
      
      <div className="mb-2">
        <p className="font-semibold">Информация:</p>
        <p className="text-sm">ID: {storageLoc.id}</p>
        <p className="text-sm">Занято грузом ID: {storageLoc.cargoId}</p>
      </div>

      <div className="p-3 bg-red-50 border border-red-200 rounded">
        <p className="font-semibold text-red-800">Внимание!</p>
        <p className="text-sm text-red-700">
          После сброса место хранения будет освобождено.
        </p>
      </div>

      <div className="flex gap-2">
        <button
          type="submit"
          className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
        >
          Сбросить
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