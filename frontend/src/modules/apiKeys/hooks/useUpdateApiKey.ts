/**
 * @file useUpdateApiKey.ts
 * @brief Hook pro úpravu metadat, oprávnění a omezení existujícího API klíče.
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
  UpdateApiKeyDocument,
  type UpdateApiKeyMutation,
  type UpdateApiKeyMutationVariables,
} from "../../../generated/graphql";

import { useApiKeysStore } from "../stores/apiKeysStore";

export const useUpdateApiKey = () => {
  const refresh = useApiKeysStore((s) => s.refresh);

  const [mutate, state] = useMutation<
    UpdateApiKeyMutation,
    UpdateApiKeyMutationVariables
  >(UpdateApiKeyDocument, {
    onCompleted: async () => {
      await refresh();
    },
  });

  return {
    updateApiKey: mutate,
    ...state,
  };
};