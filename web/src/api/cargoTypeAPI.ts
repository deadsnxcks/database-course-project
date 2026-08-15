// src/api/cargoTypeAPI.ts

export type CargoType = {
  id: number
  title: string
  processCost: number
}

type ApiErrorResponse = {
  code?: string
  message?: string
}

const BASE_URL = '/api/cargotype'

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

export const cargoTypeApi = {
  list: async (): Promise<CargoType[]> => {
    const resp = await fetch(BASE_URL)
    return handleResponse<CargoType[]>(resp)
  },

  get: async (id: number): Promise<CargoType> => {
    const resp = await fetch(`${BASE_URL}/${id}`)
    return handleResponse<CargoType>(resp)
  },

  create: async (data: Omit<CargoType, 'id'>): Promise<number> => {
    const payload = {
      title: data.title.trim(),
      processCost: Number(data.processCost),
    }

    // client-side validation
    if (!payload.title) {
      throw new Error('Название типа груза обязательно')
    }
    if (payload.processCost <= 0 || isNaN(payload.processCost)) {
      throw new Error('Стоимость обработки должна быть больше 0')
    }

    const resp = await fetch(BASE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    const result = await handleResponse<{ id: number }>(resp)

    if (result.id === undefined) {
      throw new Error('Сервер не вернул ID созданного типа груза')
    }

    return result.id
  },

  update: async (
    id: number,
    data: Partial<Omit<CargoType, 'id'>>
  ): Promise<void> => {
    const payload: any = {}

    if (data.title !== undefined) {
      payload.title = data.title.trim()
      if (!payload.title) {
        throw new Error('Название типа груза не может быть пустым')
      }
    }

    if (data.processCost !== undefined) {
      payload.processCost = Number(data.processCost)
      if (payload.processCost <= 0 || isNaN(payload.processCost)) {
        throw new Error('Стоимость обработки должна быть больше 0')
      }
    }

    if (Object.keys(payload).length === 0) {
      throw new Error('Нет данных для обновления')
    }

    const resp = await fetch(`${BASE_URL}/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    await handleResponse<void>(resp)
  },

  delete: async (id: number): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}`, {
      method: 'DELETE',
    })

    await handleResponse<void>(resp)
  },
}
