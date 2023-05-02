package types

import (
	"github.com/hyperledger/fabric-private-chaincode/lib"
	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

type SLA struct {
	iso19086.SLA
}

type Approval struct {
	lib.Approval
}

type User struct {
	lib.User
}

type Metric struct {
	iso19086.Metrics
}
