/**
 * @file sessionRedirect.ts
 * @brief Pomocná logika přesměrování podle stavu uživatelské session.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useAuthStore } from "../stores/authStore";

let redirectInFlight = false;

function isLoginRoute() {
  return window.location.pathname.startsWith("/login");
}

export function redirectToLoginForExpiredSession() {
  if (typeof window === "undefined") {
    return;
  }

  if (redirectInFlight || isLoginRoute()) {
    return;
  }

  redirectInFlight = true;
  useAuthStore.getState().logout();

  const redirect = window.location.href;
  window.location.replace(`/login?redirect=${encodeURIComponent(redirect)}`);
}

export function isUnauthorizedRuntimeError(error: unknown): boolean {
  if (!error) {
    return false;
  }

  if (typeof error === "string") {
    const normalized = error.toLowerCase();
    return (
      normalized.includes("unauthorized") ||
      normalized.includes("not authenticated") ||
      normalized.includes("missing auth")
    );
  }

  if (typeof error !== "object") {
    return false;
  }

  if ("statusCode" in error && error.statusCode === 401) {
    return true;
  }

  if ("message" in error && typeof error.message === "string") {
    return isUnauthorizedRuntimeError(error.message);
  }

  if ("errors" in error && Array.isArray(error.errors)) {
    return error.errors.some((entry) => isUnauthorizedRuntimeError(entry));
  }

  if ("cause" in error) {
    return isUnauthorizedRuntimeError(error.cause);
  }

  return false;
}
