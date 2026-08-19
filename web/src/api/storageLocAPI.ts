import { toServerDate } from '../lib/date'

const BASE_URL = '/api/storageloc'

/* =======================
   Types
======================= */

export type StorageLocType = {
  id: number
  cargoTypeId: number
  maxWeight: number
  maxVolume: number
  cargoId?: number | null
  dateOfPlacement?: string | null
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
  /** Значение из <input type="datetime-local">, без часового пояса. */
  dateOfPlacement?: string
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

const json = (method: string, body?: unknown): RequestInit => ({
  method,
  headers: { 'Content-Type': 'application/json' },
  body: body === undefined ? undefined : JSON.stringify(body),
})

/* =======================
   API
======================= */

// Сервер отдаёт даты строками RFC3339 — конвертировать на приёме нечего.
// Форматирование для показа живёт в компонентах (lib/date.ts).
export const storageLocAPI = {
  list: async (): Promise<StorageLocType[]> => {
    const resp = await fetch(BASE_URL)

    return handleResponse<StorageLocType[]>(resp)
  },

  get: async (id: number): Promise<StorageLocType> => {
    const resp = await fetch(`${BASE_URL}/${id}`)

    return handleResponse<StorageLocType>(resp)
  },

  create: async (data: StorageLocCreate): Promise<{ id: number }> => {
    if (!data.cargoTypeId || data.cargoTypeId <= 0) {
      throw new Error('cargoTypeId обязателен')
    }

    const resp = await fetch(BASE_URL, json('POST', data))

    return handleResponse<{ id: number }>(resp)
  },

  update: async (id: number, data: StorageLocUpdate): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}`, json('PUT', data))

    await handleResponse<void>(resp)
  },

  delete: async (id: number): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}`, json('DELETE'))

    await handleResponse<void>(resp)
  },

  use: async (id: number, data: StorageLocUse): Promise<void> => {
    if (!data.cargoId || data.cargoId <= 0) {
      throw new Error('cargoId обязателен')
    }

    // Дата приводится к UTC здесь, а не в форме:
    // одно место, через которое проходят все вызовы
    const resp = await fetch(
      `${BASE_URL}/${id}/use`,
      json('POST', {
        cargoId: data.cargoId,
        dateOfPlacement: toServerDate(data.dateOfPlacement),
      }),
    )

    await handleResponse<void>(resp)
  },

  reset: async (id: number): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}/reset`, json('POST'))

    await handleResponse<void>(resp)
  },
}