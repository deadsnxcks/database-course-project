import { createBrowserRouter } from "react-router-dom";
import { Navigate } from 'react-router-dom'

import EntryPage from "../pages/EntryPage";
import OperatorPage from "../pages/operator/OperatorPage";
import WarehousePage from "../pages/warehouse/WarehousePage";

import VesselSection from '../pages/operator/sections/VesselSection'
import CargoTypeSection from '../pages/operator/sections/CargoTypeSection'
import CargoSection from '../pages/operator/sections/CargoSection'
import OperationSection from '../pages/operator/sections/OperationSection'
import OperationCargoSection from '../pages/operator/sections/OperationCargoSection'
import CargoDetailReportSection from '../pages/operator/sections/CargoDetailReportSection'
import CargoTypeReportSection from '../pages/operator/sections/CargoTypeReportSection'
import StorageLocSection from '../pages/warehouse/sections/StorageLocSection'



export const router = createBrowserRouter([
  {
    path: "/",
    element: <EntryPage />,
  },
  {
    path: "/operator",
    element: <OperatorPage />,
    children: [
        { index: true, element: <Navigate to="vessels" /> },
        { path: 'vessels', element: <VesselSection /> },
        { path: 'cargo-types', element: <CargoTypeSection /> },
        { path: 'cargo', element: <CargoSection /> },
        { path: 'operations', element: <OperationSection /> },
        { path: 'operation-cargo', element: <OperationCargoSection /> },
        { path: 'report-cargo-detail', element: <CargoDetailReportSection /> },
        { path: 'report-cargo-type', element: <CargoTypeReportSection /> },
    ],
  },
  {
    path: "/warehouse",
    element: <WarehousePage />,
    children: [
        { index: true, element: <Navigate to="storageloc" /> },
        { path: 'storageloc', element: <StorageLocSection /> },
    ],
  },
]);
