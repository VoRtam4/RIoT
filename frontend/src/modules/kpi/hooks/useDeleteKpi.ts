import { useMutation } from "@apollo/client/react";
import {
  DeleteKpiDefinitionDocument,
  type DeleteKpiDefinitionMutation,
  type DeleteKpiDefinitionMutationVariables,
} from "../../../generated/graphql";
import { useKpiDefinitionsBySdTypeStore } from "../stores/kpiDefinitionsBySdTypeStore";

export const useDeleteKpi = () => {
  const [deleteMutation, { loading, error }] = useMutation<
    DeleteKpiDefinitionMutation,
    DeleteKpiDefinitionMutationVariables
  >(DeleteKpiDefinitionDocument);

  const deleteKpi = async (id: string, sdTypeID: string) => {
    const res = await deleteMutation({
      variables: { id },
    });

    const ok = res.data?.deleteKPIDefinition ?? false;

    if (ok) {
      await useKpiDefinitionsBySdTypeStore.getState().refresh(sdTypeID);
    }

    return ok;
  };

  return {
    deleteKpi,
    loading,
    error,
  };
};