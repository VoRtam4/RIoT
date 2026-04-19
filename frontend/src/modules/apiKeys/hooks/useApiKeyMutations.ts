import { useMutation } from "@apollo/client/react";
import {
  CreateApiKeyDocument,
  UpdateApiKeyDocument,
  DeleteApiKeyDocument,
  type CreateApiKeyMutation,
  type CreateApiKeyMutationVariables,
  type UpdateApiKeyMutation,
  type UpdateApiKeyMutationVariables,
  type DeleteApiKeyMutation,
  type DeleteApiKeyMutationVariables,
} from "../../../generated/graphql";

export const useApiKeyMutations = () => {
  const [createApiKey, createState] = useMutation<
    CreateApiKeyMutation,
    CreateApiKeyMutationVariables
  >(CreateApiKeyDocument);

  const [updateApiKey, updateState] = useMutation<
    UpdateApiKeyMutation,
    UpdateApiKeyMutationVariables
  >(UpdateApiKeyDocument);

  const [deleteApiKey, deleteState] = useMutation<
    DeleteApiKeyMutation,
    DeleteApiKeyMutationVariables
  >(DeleteApiKeyDocument);

  return {
    createApiKey,
    updateApiKey,
    deleteApiKey,

    loading: createState.loading || updateState.loading || deleteState.loading,

    error: createState.error || updateState.error || deleteState.error,
  };
};
