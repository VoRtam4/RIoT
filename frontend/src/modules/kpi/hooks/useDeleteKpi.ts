import { useMutation } from "@apollo/client/react";
import {
  DeleteKpiDefinitionDocument,
  type DeleteKpiDefinitionMutation,
  type DeleteKpiDefinitionMutationVariables,
} from "../../../generated/graphql";

export const useDeleteKpi = () => {
  const [deleteMutation, { loading, error }] = useMutation<
    DeleteKpiDefinitionMutation,
    DeleteKpiDefinitionMutationVariables
  >(DeleteKpiDefinitionDocument);

  const deleteKpi = async (id: string) => {
    const res = await deleteMutation({
      variables: { id },
    });

    return res.data?.deleteKPIDefinition ?? false;
  };

  return {
    deleteKpi,
    loading,
    error,
  };
};