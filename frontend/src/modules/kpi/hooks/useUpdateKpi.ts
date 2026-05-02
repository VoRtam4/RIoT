import { useMutation } from "@apollo/client/react";
import {
  UpdateKpiDefinitionDocument,
  type UpdateKpiDefinitionMutation,
  type UpdateKpiDefinitionMutationVariables,
} from "../../../generated/graphql";
import { useKpiDefinitionsBySdTypeStore } from "../stores/kpiDefinitionsBySdTypeStore";

export const useUpdateKpi = () => {
  const [mutation, { loading, error }] = useMutation<
    UpdateKpiDefinitionMutation,
    UpdateKpiDefinitionMutationVariables
  >(UpdateKpiDefinitionDocument);

  const updateKpi = async (
    id: string,
    previousSdTypeID: string,
    input: UpdateKpiDefinitionMutationVariables["input"],
  ) => {
    const res = await mutation({
      variables: { id, input },
    });

    const updated = res.data?.updateKPIDefinition ?? null;

    if (updated) {
      const store = useKpiDefinitionsBySdTypeStore.getState();

      await store.refresh(previousSdTypeID);

      if (input.sdTypeID && input.sdTypeID !== previousSdTypeID) {
        await store.refresh(input.sdTypeID);
      }
    }

    return updated;
  };

  return {
    updateKpi,
    loading,
    error,
  };
};