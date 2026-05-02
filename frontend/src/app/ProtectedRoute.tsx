import type { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../modules/auth/hooks/useAuth";
import { useAuthStore } from "../modules/auth/stores/authStore";

type Props = {
  children: ReactNode;
};

export function ProtectedRoute({ children }: Props) {
  const { loading } = useAuth();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const location = useLocation();

  if (loading) {
    return (
      <div className="d-flex vh-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" role="status" />
      </div>
    );
  }

  if (!isAuthenticated) {
    const absoluteRedirect =
      window.location.origin + location.pathname + location.search;

    return (
      <Navigate
        to={`/login?redirect=${encodeURIComponent(absoluteRedirect)}`}
        replace
      />
    );
  }

  return <>{children}</>;
}
