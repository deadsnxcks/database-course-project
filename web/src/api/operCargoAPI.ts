export type OperCargoType = {
  operationId: number
  cargoId: number
}

type ApiErrorResponse = {
  code?: string
  message?: string
}

const BASE_URL = '/api/opercargo'

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

export const operCargoApi = {
  list: async (): Promise<OperCargoType[]> => {
    const resp = await fetch(BASE_URL)
    return handleResponse<OperCargoType[]>(resp)
  },

  create: async (data: OperCargoType): Promise<void> => {
    const payload = {
      operationId: Number(data.operationId),
      cargoId: Number(data.cargoId),
    }

    if (payload.operationId <= 0) {
      throw new Error('Некорректный ID операции')
    }
    if (payload.cargoId <= 0) {
      throw new Error('Некорректный ID груза')
    }

    const resp = await fetch(BASE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    await handleResponse<void>(resp)
  },

  delete: async (operationId: number, cargoId: number): Promise<void> => {
    const payload = {
      operationId: Number(operationId),
      cargoId: Number(cargoId),
    }

    if (payload.operationId <= 0) {
      throw new Error('Некорректный ID операции')
    }
    if (payload.cargoId <= 0) {
      throw new Error('Некорректный ID груза')
    }

    const resp = await fetch(BASE_URL, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    await handleResponse<void>(resp)
  },
}
