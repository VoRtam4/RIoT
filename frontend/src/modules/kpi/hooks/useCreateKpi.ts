import { useMutation } from "@apollo/client/react";
import {
  CreateKpiDefinitionDocument,
  type CreateKpiDefinitionMutation,
  type CreateKpiDefinitionMutationVariables,
} from "../../../generated/graphql";

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

    return res.data?.createKPIDefinition ?? null;
  };

  return {
    createKpi,
    loading,
    error,
  };
};
