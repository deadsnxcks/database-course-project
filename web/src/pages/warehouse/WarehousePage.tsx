import { Outlet } from 'react-router-dom'
import WarehouseLayout from './WarehouseLayout'

export default function WarehousePage() {
  return (
    <WarehouseLayout>
      <Outlet />
    </WarehouseLayout>
  )
}
