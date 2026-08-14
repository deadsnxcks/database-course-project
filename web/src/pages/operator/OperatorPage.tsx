import { Outlet } from 'react-router-dom'
import OperatorLayout from './OperatorLayout'

export default function OperatorPage() {
  return (
    <OperatorLayout>
      <Outlet />
    </OperatorLayout>
  )
}
