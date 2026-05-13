/**
 * @file sdTypes_test.go
 * @brief Test mapování typů zdrojů dat z databázového modelu do doménového modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní test mapování typů zdrojů dat.
 * - Vojtěch Hubáček: úprava testovaných očekávání pro labely a parametry typů zdrojů.
 *
 * @ingroup riot_backend_core
 */
package db2dll

import (
	"testing"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/go-cmp/cmp"
)

func TestToDLLModelSDType(t *testing.T) {
	sdTypeEntity := dbModel.SDTypeEntity{
		ID:  1,
		UID: "shelly1pro",
		Parameters: sharedUtils.SliceOf(dbModel.SDParameterEntity{
			ID:         1,
			SDTypeID:   1,
			Denotation: "relay_0_temperature",
			Type:       "number",
		}, dbModel.SDParameterEntity{
			ID:         2,
			SDTypeID:   1,
			Denotation: "relay_0_source",
			Type:       "string",
		}),
	}
	expected := dllModel.SDType{
		ID:  sharedUtils.NewOptionalOf[uint32](1),
		UID: "shelly1pro",
		Parameters: sharedUtils.SliceOf(dllModel.SDParameter{
			ID:         sharedUtils.NewOptionalOf[uint32](1),
			Denotation: "relay_0_temperature",
			Type:       dllModel.SDParameterTypeNumber,
		}, dllModel.SDParameter{
			ID:         sharedUtils.NewOptionalOf[uint32](2),
			Denotation: "relay_0_source",
			Type:       dllModel.SDParameterTypeString,
		}),
	}
	actual := ToDLLModelSDType(sdTypeEntity)
	if diff := cmp.Diff(expected, actual, cmp.Comparer(sharedUtils.OptionalComparer[uint32])); diff != "" {
		t.Fatalf("ToDLLModelSDType() mismatch (-expected +actual):\n%s", diff)
	}
}
