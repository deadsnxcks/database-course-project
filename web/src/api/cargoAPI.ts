// src/api/cargoAPI.ts

export type Cargo = {
  id: number
  title: string
  typeId: number
  weight: number
  volume: number
  vesselId: number
}

type ApiErrorResponse = {
  code?: string
  message?: string
}

const BASE_URL = 'http://localhost:8081/cargo'

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

export const cargoAPI = {
  list: async (): Promise<Cargo[]> => {
    const resp = await fetch(BASE_URL)
    return handleResponse<Cargo[]>(resp)
  },

  get: async (id: number): Promise<Cargo> => {
    const resp = await fetch(`${BASE_URL}/${id}`)
    return handleResponse<Cargo>(resp)
  },

  create: async (data: Omit<Cargo, 'id'>): Promise<number> => {
    const payload = {
      title: data.title.trim(),
      typeId: Number(data.typeId),
      weight: Number(data.weight),
      volume: Number(data.volume),
      vesselId: Number(data.vesselId),
    }

    // клиентская валидация
    if (!payload.title) {
      throw new Error('Название груза обязательно')
    }
    if (payload.typeId <= 0 || isNaN(payload.typeId)) {
      throw new Error('Некорректный тип груза')
    }
    if (payload.vesselId <= 0 || isNaN(payload.vesselId)) {
      throw new Error('Некорректное судно')
    }
    if (payload.weight <= 0 || isNaN(payload.weight)) {
      throw new Error('Вес должен быть больше 0')
    }
    if (payload.volume <= 0 || isNaN(payload.volume)) {
      throw new Error('Объём должен быть больше 0')
    }

    let resp: Response
    try {
      resp = await fetch(BASE_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
    } catch {
      throw new Error('Не удалось связаться с сервером')
    }

    const result = await handleResponse<{ id: number }>(resp)

    if (result.id === undefined) {
      throw new Error('Сервер не вернул ID созданного груза')
    }

    return result.id
  },

  update: async (
    id: number,
    data: Partial<Omit<Cargo, 'id'>>
  ): Promise<void> => {
    const payload: any = {}

    if (data.title !== undefined) {
      payload.title = data.title.trim()
      if (!payload.title) {
        throw new Error('Название груза не может быть пустым')
      }
    }

    if (data.typeId !== undefined) {
      payload.typeId = Number(data.typeId)
      if (payload.typeId <= 0 || isNaN(payload.typeId)) {
        throw new Error('Некорректный тип груза')
      }
    }

    if (data.vesselId !== undefined) {
      payload.vesselId = Number(data.vesselId)
      if (payload.vesselId <= 0 || isNaN(payload.vesselId)) {
        throw new Error('Некорректное судно')
      }
    }

    if (data.weight !== undefined) {
      payload.weight = Number(data.weight)
      if (payload.weight <= 0 || isNaN(payload.weight)) {
        throw new Error('Вес должен быть больше 0')
      }
    }

    if (data.volume !== undefined) {
      payload.volume = Number(data.volume)
      if (payload.volume <= 0 || isNaN(payload.volume)) {
        throw new Error('Объём должен быть больше 0')
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
