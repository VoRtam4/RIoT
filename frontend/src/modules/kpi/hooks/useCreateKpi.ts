/**
 * @file useCreateKpi.ts
 * @brief Hook pro vytvoření KPI definice z editoru podmínek.
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

    if (created && input.sdTypeUID) {
      await useKpiDefinitionsBySdTypeStore.getState().refresh(input.sdTypeUID);
    }

    return created;
  };

  return {
    createKpi,
    loading,
    error,
  };
};
