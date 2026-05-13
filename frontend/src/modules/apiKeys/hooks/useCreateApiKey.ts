/**
 * @file useCreateApiKey.ts
 * @brief Hook pro vytvoření API klíče včetně oprávnění a IP omezení.
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
  CreateApiKeyDocument,
  type CreateApiKeyMutation,
  type CreateApiKeyMutationVariables,
} from "../../../generated/graphql";

import { useApiKeysStore } from "../stores/apiKeysStore";

export const useCreateApiKey = () => {
  const refresh = useApiKeysStore((s) => s.refresh);

  const [mutate, state] = useMutation<
    CreateApiKeyMutation,
    CreateApiKeyMutationVariables
  >(CreateApiKeyDocument, {
    onCompleted: async () => {
      await refresh();
    },
  });

  return {
    createApiKey: mutate,
    ...state,
  };
};
