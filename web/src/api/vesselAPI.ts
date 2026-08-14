// src/api/vesselAPI.ts

export type VesselType = {
  id: number
  title: string
  vesselType: string
  maxLoad: number
}

type ApiErrorResponse = {
  code?: string
  message?: string
}

const BASE_URL = '/api/vessel'

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

export const vesselApi = {
  listVessels: async (): Promise<VesselType[]> => {
    const resp = await fetch(BASE_URL)
    return handleResponse<VesselType[]>(resp)
  },

  getVessel: async (id: number): Promise<VesselType> => {
    const resp = await fetch(`${BASE_URL}/${id}`)
    return handleResponse<VesselType>(resp)
  },

  createVessel: async (data: Omit<VesselType, 'id'>): Promise<number> => {
    const payload = {
      title: data.title.trim(),
      vesselType: data.vesselType.trim(),
      maxLoad: Number(data.maxLoad),
    }

    // клиентская валидация
    if (!payload.title) {
      throw new Error('Название судна обязательно')
    }
    if (!payload.vesselType) {
      throw new Error('Тип судна обязателен')
    }
    if (payload.maxLoad <= 0 || isNaN(payload.maxLoad)) {
      throw new Error('Максимальная нагрузка должна быть больше 0')
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

    // 🔥 ВОТ ЗДЕСЬ обработка ошибок
    const result = await handleResponse<{ id: number }>(resp)

    if (result.id === undefined) {
      throw new Error('Сервер не вернул ID созданного судна')
    }

    return result.id
  },


  updateVessel: async (
    id: number,
    data: Partial<Omit<VesselType, 'id'>>
  ): Promise<void> => {
    const payload: any = {}

    if (data.title !== undefined) {
      payload.title = data.title.trim()
      if (!payload.title) {
        throw new Error('Название судна не может быть пустым')
      }
    }

    if (data.vesselType !== undefined) {
      payload.vesselType = data.vesselType.trim()
      if (!payload.vesselType) {
        throw new Error('Тип судна не может быть пустым')
      }
    }

    if (data.maxLoad !== undefined) {
      payload.maxLoad = Number(data.maxLoad)
      if (payload.maxLoad <= 0 || isNaN(payload.maxLoad)) {
        throw new Error('Максимальная нагрузка должна быть больше 0')
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

  deleteVessel: async (id: number): Promise<void> => {
    const resp = await fetch(`${BASE_URL}/${id}`, {
      method: 'DELETE',
    })

    await handleResponse<void>(resp)
  },
}
