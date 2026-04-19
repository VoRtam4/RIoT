import { useMutation } from "@apollo/client/react";
import {
  UpdateKpiDefinitionDocument,
  type UpdateKpiDefinitionMutation,
  type UpdateKpiDefinitionMutationVariables,
} from "../../../generated/graphql";

export const useUpdateKpi = () => {
  const [mutation, { loading, error }] = useMutation<
    UpdateKpiDefinitionMutation,
    UpdateKpiDefinitionMutationVariables
  >(UpdateKpiDefinitionDocument);

  const updateKpi = async (
    id: string,
    input: UpdateKpiDefinitionMutationVariables["input"],
  ) => {
    const res = await mutation({
      variables: { id, input },
    });

    return res.data?.updateKPIDefinition ?? null;
  };

  return {
    updateKpi,
    loading,
    error,
  };
};
