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