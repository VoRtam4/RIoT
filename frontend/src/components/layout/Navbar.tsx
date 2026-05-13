/**
 * @file Navbar.tsx
 * @brief Horní navigace aplikace s odkazy na hlavní sekce a uživatelské akce.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { NavLink } from "react-router-dom";

export default function Navbar() {
  const handleLogout = () => {
    window.location.href = `/auth/logout?redirect=${encodeURIComponent(window.location.origin)}`;
  };

  return (
    <div>
      <style>
        {`
          .navbar {
            border-top-left-radius: 0;
            border-top-right-radius: 0;
          }
        `}
      </style>
      <nav className="navbar navbar-dark px-3 d-flex justify-content-between">
        <div className="d-flex align-items-center gap-3">
          <span className="navbar-brand mb-0">RIoT</span>

          <div className="d-flex gap-2">
            <NavLink
              to="/"
              end
              className={({ isActive }) =>
                `btn btn-sm ${isActive ? "btn-primary" : "btn-outline-light"}`
              }
            >
              Dashboard
            </NavLink>

            <NavLink
              to="/sd-instance"
              className={({ isActive }) =>
                `btn btn-sm ${isActive ? "btn-primary" : "btn-outline-light"}`
              }
            >
              Devices
            </NavLink>

            <NavLink
              to="/kpi"
              className={({ isActive }) =>
                `btn btn-sm ${isActive ? "btn-primary" : "btn-outline-light"}`
              }
            >
              KPI
            </NavLink>

            <NavLink
              to="/history"
              className={({ isActive }) =>
                `btn btn-sm ${isActive ? "btn-primary" : "btn-outline-light"}`
              }
            >
              History
            </NavLink>

            <NavLink
              to="/api-keys"
              className={({ isActive }) =>
                `btn btn-sm ${isActive ? "btn-primary" : "btn-outline-light"}`
              }
            >
              API Keys
            </NavLink>
          </div>
        </div>

        <button className="btn btn-outline-light" onClick={handleLogout}>
          Logout
        </button>
      </nav>
    </div>
  );
}
