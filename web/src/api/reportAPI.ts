const BASE_URL = '/api'

export type CargoDetailReportItem = {
  id: string | number
  cargoName: string
  weight: number
  cargoType: string
  vesselName: string
  unloadDate: string
}

export type CargoTypeReportItem = {
  id: string | number
  cargoTypeName: string
  count: number
  weight: number
  volume: number
  processCost: number
}

export const reportAPI = {
  getCargoDetailReport: async (): Promise<CargoDetailReportItem[]> => {
    const resp = await fetch(`${BASE_URL}/report-cargo-detail`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    })
    
    console.log(resp)
    if (!resp.ok) {
      throw new Error(`Ошибка HTTP: ${resp.status}`)
    }
    
    const items = await resp.json()
    
    // Просто добавляем id к каждому элементу
    return items.map((item: any, index: number) => ({
      ...item,
      id: index + 1,
    }))
  },

  getCargoTypeReport: async (): Promise<CargoTypeReportItem[]> => {
    const resp = await fetch(`${BASE_URL}/report-cargo-type`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    })
    
    console.log(resp)
    if (!resp.ok) {
      throw new Error(`Ошибка HTTP: ${resp.status}`)
    }
    
    const items = await resp.json()
    
    return items.map((item: any, index: number) => ({
      ...item,
      id: index + 1,
    }))
  },
}