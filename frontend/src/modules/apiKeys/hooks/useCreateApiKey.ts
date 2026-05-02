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
