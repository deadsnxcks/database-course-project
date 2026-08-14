const BASE_URL = '/api/storageloc'

/* =======================
   Types
======================= */

export type StorageLocType = {
  id: number
  cargoTypeId: number
  maxWeight: number
  maxVolume: number
  cargoId?: number
  dateOfPlacement?: string 
}

export type StorageLocCreate = {
  cargoTypeId: number
  maxWeight: number
  maxVolume: number
}

export type StorageLocUpdate = {
  cargoTypeId?: number
  maxWeight?: number
  maxVolume?: number
}

export type StorageLocUse = {
  cargoId: number
  dateOfPlacement: string
}

/* =======================
   Error handling
======================= */

type ApiErrorResponse = {
  code?: string
  message?: string
}

const handleResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    let body: ApiErrorResponse = {}

    try {
      body = await response.json()
    } catch {}

    const error = new Error(body.message || response.statusText)
    ;(error as any).code = body.code
    ;(error as any).status = response.status

    throw error
  }

  if (response.status === 204) {
    return {} as T
  }

  return response.json() as Promise<T>
}

/* =======================
   API
======================= */

export const storageLocAPI = {
  /* ---------- GET ---------- */

  list: async (): Promise<StorageLocType[]> => {
  const resp = await fetch(BASE_URL)
  const data = await handleResponse<any[]>(resp)
  
  // Преобразуй timestamp в строку
  return data.map(item => ({
      ...item,
      dateOfPlacement: item.dateOfPlacement 
        ? convertTimestampToString(item.dateOfPlacement)
        : undefined
    }))
  },

  get: async (id: number): Promise<StorageLocType> => {
    const resp = await fetch(`${BASE_URL}/${id}`)
    const data = await handleResponse<any>(resp)
    
    return {
      ...data,
      dateOfPlacement: data.dateOfPlacement 
        ? convertTimestampToString(data.dateOfPlacement)
        : undefined
    }
  },

  /* ---------- CREATE ---------- */

  create: async (data: StorageLocCreate): Promise<void> => {
    if (!data.cargoTypeId || data.cargoTypeId <= 0) {
      throw new Error('cargoTypeId обязателен')
    }

    const resp = await fetch(BASE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })

    await handleResponse<void>(resp)
  },

  /* ---------- UPDATE ---------- */

  update: async (id: number, data: StorageLocUpdate): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })

    await handleResponse<void>(resp)
  },

  /* ---------- DELETE ---------- */

  delete: async (id: number): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}`, {
      method: 'DELETE',
    })

    await handleResponse<void>(resp)
  },

  /* ---------- USE ---------- */

  use: async (id: number, data: StorageLocUse): Promise<void> => {
    if (!data.cargoId || data.cargoId <= 0) {
      throw new Error('cargoId обязателен')
    }

    const resp = await fetch(`${BASE_URL}/${id}/use`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })

    await handleResponse<void>(resp)
  },

  /* ---------- RESET ---------- */

  reset: async (id: number): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}/reset`, {
      method: 'POST',
    })

    await handleResponse<void>(resp)
  },
}

const convertTimestampToString = (timestamp: any): string => {
  if (!timestamp) return ''
  
  try {
    let date: Date
    
    if (typeof timestamp === 'string') {
      date = new Date(timestamp)
    } else if (timestamp && typeof timestamp === 'object' && 'seconds' in timestamp) {
      // gRPC timestamp уже в UTC
      const seconds = Number(timestamp.seconds)
      const nanos = timestamp.nanos ? Number(timestamp.nanos) / 1000000 : 0
      date = new Date(seconds * 1000 + nanos)
    } else {
      date = new Date(timestamp)
    }
    
    if (isNaN(date.getTime())) {
      return ''
    }
    
    // Указываем UTC чтобы убрать смещение +10 часов
    return new Intl.DateTimeFormat('ru-RU', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC' // ← ВАЖНО!
    }).format(date)
    
  } catch (error) {
    console.error('Error converting timestamp:', error, timestamp)
    return ''
  }
}