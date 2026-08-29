/**
 * @file useUpdateKpi.ts
 * @brief Hook pro uložení změn existující KPI definice.
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
    uid: string,
    previousSdTypeUID: string,
    input: UpdateKpiDefinitionMutationVariables["input"],
  ) => {
    const res = await mutation({
      variables: { uid, input },
    });

    const updated = res.data?.updateKPIDefinition ?? null;

    if (updated) {
      const store = useKpiDefinitionsBySdTypeStore.getState();

      await store.refresh(previousSdTypeUID);

      if (input.sdTypeUID && input.sdTypeUID !== previousSdTypeUID) {
        await store.refresh(input.sdTypeUID);
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
