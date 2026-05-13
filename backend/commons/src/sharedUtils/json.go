/**
 * @file json.go
 * @brief Pomocné funkce pro serializaci, deserializaci a porovnávání JSON hodnot.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní serializační a deserializační pomocné funkce.
 * - Vojtěch Hubáček: doplnění porovnávání JSON struktur pomocí CompareJSONs.
 *
 * @ingroup riot_commons
 */
package sharedUtils

import (
	"encoding/json"
)

func DeserializeFromJSON[T any](data []byte) Result[T] {
	var object T
	err := json.Unmarshal(data, &object)
	if err != nil {
		return NewFailureResult[T](err)
	}
	return NewSuccessResult[T](object)
}

func SerializeToJSON(object any) Result[[]byte] {
	data, err := json.Marshal(object)
	if err != nil {
		return NewFailureResult[[]byte](err)
	}
	return NewSuccessResult[[]byte](data)
}

func CompareJSONs(a interface{}, b interface{}) bool {
	aj := SerializeToJSON(a)
	bj := SerializeToJSON(b)
	if aj.IsFailure() || bj.IsFailure() {
		return false
	}
	return string(aj.GetPayload()) == string(bj.GetPayload())
}
