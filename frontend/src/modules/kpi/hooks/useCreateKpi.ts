import { useMutation } from "@apollo/client/react";
import {
  CreateKpiDefinitionDocument,
  type CreateKpiDefinitionMutation,
  type CreateKpiDefinitionMutationVariables,
} from "../../../generated/graphql";
import { useKpiDefinitionsBySdTypeStore } from "../stores/kpiDefinitionsBySdTypeStore";

export const useCreateKpi = () => {
  const [createKpiMutation, { loading, error }] = useMutation<
    CreateKpiDefinitionMutation,
    CreateKpiDefinitionMutationVariables
  >(CreateKpiDefinitionDocument);

  const createKpi = async (
    input: CreateKpiDefinitionMutationVariables["input"],
  ) => {
    const res = await createKpiMutation({
      variables: { input },
    });

    const created = res.data?.createKPIDefinition ?? null;

    if (created && input.sdTypeID) {
      await useKpiDefinitionsBySdTypeStore
        .getState()
        .refresh(input.sdTypeID);
    }

    return created;
  };

  return {
    createKpi,
    loading,
    error,
  };
};