/**
 * @file useAuth.ts
 * @brief Hook pro přihlášení, odhlášení a načtení aktuální session.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useEffect } from "react";
import { useQuery } from "@apollo/client/react";
import {
  UserConfigDocument,
  type UserConfigQuery,
} from "../../../generated/graphql";
import { useAuthStore } from "../stores/authStore";

export const useAuth = () => {
  const setAuth = useAuthStore((s) => s.setAuth);

  const { data, loading, error } = useQuery<UserConfigQuery>(
    UserConfigDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const userId = data?.userConfig?.userUID ?? null;

  useEffect(() => {
    if (userId) {
      setAuth(userId);
    }
  }, [userId, setAuth]);

  return {
    userId,
    loading,
    error,
    isAuthenticated: !!userId,
  };
};
