package contracts

import (
	"strings"
)

func SLAInUserContracts(providerOfList, clientOfList, slaId string) bool {
	// Slice of size 1 means that the delimiter was not found in the string
	for _, sla := range GetIDsFromString(clientOfList) {
		if sla == slaId {
			return true
		}
	}

	// Slice of size 1 means that the delimiter was not found in the string
	for _, sla := range GetIDsFromString(providerOfList) {
		if sla == slaId {
			return true
		}
	}
	return false
}

func GetIDsFromString(idString string) []string {
	return strings.Split(idString, ",")
}
