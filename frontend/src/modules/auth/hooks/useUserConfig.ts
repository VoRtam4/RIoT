import { useMemo } from "react";
import { useQuery, useMutation } from "@apollo/client/react";

import {
  UserConfigDocument,
  UpdateUserConfigDocument,
  type UserConfigQuery,
} from "../../../generated/graphql";

type Config = {
  favoriteSdInstances: string[];
  favoriteKpis: string[];
};

const DEFAULT_CONFIG: Config = {
  favoriteSdInstances: [],
  favoriteKpis: [],
};

export const useUserConfig = () => {
  const { data, loading, refetch } = useQuery<UserConfigQuery>(
    UserConfigDocument,
    {
      fetchPolicy: "cache-and-network",
    },
  );

  const [update] = useMutation(UpdateUserConfigDocument);

  const config: Config = useMemo(() => {
    try {
      const raw = data?.userConfig?.config;

      if (!raw) return DEFAULT_CONFIG;

      const parsed =
        typeof raw === "string" ? JSON.parse(raw) : raw;

      return {
        ...DEFAULT_CONFIG,
        ...parsed,
      };
    } catch {
      return DEFAULT_CONFIG;
    } 
  }, [data]);

  const save = async (newConfig: Config) => {
    await update({
      variables: {
        input: {
          config: JSON.stringify(newConfig),
        },
      },
    });
    await refetch();
  };

  const toggleSdInstance = async (id: string) => {
    const exists = config.favoriteSdInstances.includes(id);

    const updated: Config = {
      ...config,
      favoriteSdInstances: exists
        ? config.favoriteSdInstances.filter((x) => x !== id)
        : [...config.favoriteSdInstances, id],
    };

    await save(updated);
  };

  const toggleKpi = async (id: string) => {
    const exists = config.favoriteKpis.includes(id);

    const updated: Config = {
      ...config,
      favoriteKpis: exists
        ? config.favoriteKpis.filter((x) => x !== id)
        : [...config.favoriteKpis, id],
    };

    await save(updated);
  };

  return {
    config,
    loading,
    toggleSdInstance,
    toggleKpi,
  };
};