import { useEffect } from "react";
import { useQuery } from "@apollo/client/react";
import {
  UserConfigDocument,
  type UserConfigQuery,
} from "../../../generated/graphql";
import { useAuthStore } from "../store/authStore";

export const useAuth = () => {
  const setAuth = useAuthStore((s) => s.setAuth);

  const { data, loading, error } = useQuery<UserConfigQuery>(
    UserConfigDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const userId = data?.userConfig?.userID ?? null;

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
