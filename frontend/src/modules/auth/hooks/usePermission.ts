import { useEffect } from "react";
import { useQuery } from "@apollo/client/react";
import { RoleDocument, type RoleQuery } from "../../../generated/graphql";
import { useAuthStore } from "../store/authStore";

export const usePermissions = () => {
  const userId = useAuthStore((s) => s.userId);
  const permissions = useAuthStore((s) => s.permissions);
  const setPermissions = useAuthStore((s) => s.setPermissions);

  const { data, loading, error } = useQuery<RoleQuery>(RoleDocument, {
    skip: !userId || !!permissions,
    fetchPolicy: "cache-first",
  });

  useEffect(() => {
    if (data?.role?.permissions) {
      setPermissions(data.role.permissions);
    }
  }, [data, setPermissions]);

  return {
    permissions: permissions ?? [],
    loading,
    error,
  };
};
