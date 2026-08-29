/**
 * @file useDeleteKpi.ts
 * @brief Hook pro odstranění KPI definice.
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

  const deleteKpi = async (uid: string, sdTypeUID: string) => {
    const res = await deleteMutation({
      variables: { uid },
    });

    const ok = res.data?.deleteKPIDefinition ?? false;

    if (ok) {
      await useKpiDefinitionsBySdTypeStore.getState().refresh(sdTypeUID);
    }

    return ok;
  };

  return {
    deleteKpi,
    loading,
    error,
  };
};
