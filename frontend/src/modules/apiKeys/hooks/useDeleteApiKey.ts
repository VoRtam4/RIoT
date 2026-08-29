/**
 * @file useDeleteApiKey.ts
 * @brief Hook pro odstranění API klíče přes GraphQL API.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useMutation } from "@apollo/client/react";
import {
  DeleteApiKeyDocument,
  type DeleteApiKeyMutation,
  type DeleteApiKeyMutationVariables,
} from "../../../generated/graphql";

import { useApiKeysStore } from "../stores/apiKeysStore";

export const useDeleteApiKey = () => {
  const refresh = useApiKeysStore((s) => s.refresh);

  const [mutate, state] = useMutation<
    DeleteApiKeyMutation,
    DeleteApiKeyMutationVariables
  >(DeleteApiKeyDocument, {
    onCompleted: async () => {
      await refresh();
    },
  });

  return {
    deleteApiKey: mutate,
    ...state,
  };
};
