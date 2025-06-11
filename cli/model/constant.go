package model

import (
	"fmt"
	"math/big"
)

const (
	InvokeContract = 3844450837

	Page     = 1
	PageSize = 10

	Approved = "Approved"
	Rejected = "Rejected"

	Approve             = "approve"
	Reject              = "reject"
	BaseProposalAPIPath = "/api"

	GithubAPI = "https://api.github.com/users/"
)

const (
	ProposalStatusPending ProposalStatus = iota + 1
	ProposalStatusInProgress
	ProposalStatusCounting
	ProposalStatusCompleted
)

const (
	attoFIL   = 1
	femtoFIL  = attoFIL * 1_000
	picoFIL   = femtoFIL * 1_000
	nanoFIL   = picoFIL * 1_000
	microFIL  = nanoFIL * 1_000
	milliFIL  = microFIL * 1_000
	FIL       = milliFIL * 1_000
	tFIL      = FIL
	milliTFIL = microFIL * 1_000
	microTFIL = nanoFIL * 1_000
	nanoTFIL  = picoFIL * 1_000
	picoTFIL  = femtoFIL * 1_000
	femtoTFIL = attoFIL * 1_000
)

func ConvertToFIL(value *big.Int, chainId int) string {
	var tokenName string
	var unit int64
	var milliUnit, microUnit, nanoUnit, picoUnit, femtoUnit int64

	switch chainId {
	case 314159:
		tokenName = "tFIL"
		unit = tFIL
		milliUnit = milliTFIL
		microUnit = microTFIL
		nanoUnit = nanoTFIL
		picoUnit = picoTFIL
		femtoUnit = femtoTFIL
	case 314:
		tokenName = "FIL"
		unit = FIL
		milliUnit = milliFIL
		microUnit = microFIL
		nanoUnit = nanoFIL
		picoUnit = picoFIL
		femtoUnit = femtoFIL
	}

	valFloat := new(big.Float).SetInt(value)

	if value.Cmp(big.NewInt(unit)) >= 0 {
		return fmt.Sprintf("%s %s", new(big.Float).Quo(valFloat, new(big.Float).SetFloat64(float64(unit))).Text('f', 2), tokenName)
	} else if value.Cmp(big.NewInt(milliUnit)) >= 0 {
		return fmt.Sprintf("%s milli%s", new(big.Float).Quo(valFloat, new(big.Float).SetFloat64(float64(milliUnit))).Text('f', 2), tokenName)
	} else if value.Cmp(big.NewInt(microUnit)) >= 0 {
		return fmt.Sprintf("%s micro%s", new(big.Float).Quo(valFloat, new(big.Float).SetFloat64(float64(microUnit))).Text('f', 2), tokenName)
	} else if value.Cmp(big.NewInt(nanoUnit)) >= 0 {
		return fmt.Sprintf("%s nano%s", new(big.Float).Quo(valFloat, new(big.Float).SetFloat64(float64(nanoUnit))).Text('f', 2), tokenName)
	} else if value.Cmp(big.NewInt(picoUnit)) >= 0 {
		return fmt.Sprintf("%s pico%s", new(big.Float).Quo(valFloat, new(big.Float).SetFloat64(float64(picoUnit))).Text('f', 2), tokenName)
	} else if value.Cmp(big.NewInt(femtoUnit)) >= 0 {
		return fmt.Sprintf("%s femto%s", new(big.Float).Quo(valFloat, new(big.Float).SetFloat64(float64(femtoUnit))).Text('f', 2), tokenName)
	} else {
		return fmt.Sprintf("%s atto%s", value.String(), tokenName)
	}
}

// ProposalStatus defines the various states a proposal can have.
type ProposalStatus int

// String returns the string representation of the ProposalStatus.
func (status ProposalStatus) String() string {
	switch status {
	case ProposalStatusPending:
		return "Upcoming"
	case ProposalStatusInProgress:
		return "In Progress"
	case ProposalStatusCounting:
		return "Vote Counting"
	case ProposalStatusCompleted:
		return "Complete"
	default:
		return "Unknown"
	}
}
