/**
 * @file LoginPage.tsx
 * @brief Přihlašovací stránka napojená na autentizační vrstvu Backend Core.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { Navigate } from "react-router-dom";
import { useAuth } from "../modules/auth/hooks/useAuth";
import { useAuthStore } from "../modules/auth/stores/authStore";
import { colors } from "../theme/colors";

export default function LoginPage() {
  const { loading } = useAuth();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  const params = new URLSearchParams(window.location.search);
  const redirect = params.get("redirect") || window.location.origin;

  if (loading) {
    return (
      <div className="d-flex vh-100 justify-content-center align-items-center">
        <div className="spinner-border text-primary" role="status" />
      </div>
    );
  }

  if (isAuthenticated) {
    try {
      const url = new URL(redirect);
      return <Navigate to={url.pathname + url.search} replace />;
    } catch {
      return <Navigate to="/" replace />;
    }
  }

  const handleLogin = () => {
    window.location.href = `/auth/login?redirect=${encodeURIComponent(
      redirect,
    )}`;
  };

  return (
    <div className="d-flex vh-100 justify-content-center align-items-center flex-column">
      <h1 className="mb-4" style={{ color: colors.textMain }}>
        RIoT
      </h1>
      <button className="btn btn-primary btn-lg" onClick={handleLogin}>
        <i className="fab fa-google me-2"></i>
        Sign in with Google
      </button>
    </div>
  );
}
