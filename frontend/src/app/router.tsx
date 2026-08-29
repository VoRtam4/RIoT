/**
 * @file router.tsx
 * @brief Definice klientských rout a napojení stránek na chráněné části aplikace.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { Navigate, createBrowserRouter } from "react-router-dom";
import LoginPage from "../pages/LoginPage";
import DashboardPage from "../pages/DashboardPage";
import APIKeysPage from "../pages/APIKeysPage";
import KpiPage from "../pages/KpiPage";
import KpiEditorPage from "../pages/KpiEditorPage";
import TimeSeriesPage from "../pages/TimeSeriesPage";
import KpiDetailPage from "../pages/KpiDetailPage";
import SdInstancePage from "../pages/SdInstancePage";
import SdInstanceDetailPage from "../pages/SdInstanceDetailPage";
import { ProtectedRoute } from "./ProtectedRoute";
import Layout from "../components/layout/Layout";
import ApiDocsPage from "../pages/APIDocsPages";

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    element: (
      <ProtectedRoute>
        <Layout />
      </ProtectedRoute>
    ),
    children: [
      {
        path: "/",
        element: <DashboardPage />,
      },
      {
        path: "/api-keys",
        element: <APIKeysPage />,
      },
      {
        path: "/api-keys/docs",
        element: <ApiDocsPage />,
      },
      {
        path: "/kpi",
        element: <KpiPage />,
      },
      {
        path: "/kpi/:uid",
        element: <KpiDetailPage />,
      },
      {
        path: "/kpi/new",
        element: <KpiEditorPage />,
      },
      {
        path: "/kpi/edit/:uid",
        element: <KpiEditorPage />,
      },
      {
        path: "/sd-instance",
        element: <SdInstancePage />,
      },
      {
        path: "/sd-instance/:uid",
        element: <SdInstanceDetailPage />,
      },
      {
        path: "/history",
        element: <TimeSeriesPage />,
      },
      {
        path: "*",
        element: <Navigate to="/" replace />,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" replace />,
  },
]);
