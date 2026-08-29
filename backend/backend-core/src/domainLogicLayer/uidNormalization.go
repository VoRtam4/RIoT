/**
 * @file uidNormalization.go
 * @brief Pomocné normalizace veřejných UID v doménové vrstvě.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import "github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"

func normalizeSDTypeUID(uid string) (string, error) {
	normalizedUID, _, err := sharedUtils.NormalizePrefixedUID(uid, "sdt", "SD type")
	return normalizedUID, err
}

func normalizeAPIKeyUID(uid string) (string, error) {
	normalizedUID, _, err := sharedUtils.NormalizePrefixedUID(uid, "ak", "API key")
	return normalizedUID, err
}

func normalizeUserUID(uid string) (string, error) {
	normalizedUID, _, err := sharedUtils.NormalizePrefixedUID(uid, "usr", "user")
	return normalizedUID, err
}

func normalizeRoleUID(uid string) (string, error) {
	normalizedUID, _, err := sharedUtils.NormalizePrefixedUID(uid, "role", "role")
	return normalizedUID, err
}

func normalizeExportUID(uid string) (string, error) {
	normalizedUID, _, err := sharedUtils.NormalizePrefixedUID(uid, "exp", "time series export")
	return normalizedUID, err
}

func normalizeUserSessionUID(uid string) (string, error) {
	normalizedUID, _, err := sharedUtils.NormalizePrefixedUID(uid, "ses", "user session")
	return normalizedUID, err
}
