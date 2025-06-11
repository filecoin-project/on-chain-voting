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
	var unit, milliUnit, microUnit, nanoUnit, picoUnit, femtoUnit int64

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

	convertToString := func(unit int64) string {
		tempVal := new(big.Int).Div(value, big.NewInt(unit))
		decimalPart := new(big.Int).Mul(new(big.Int).Mod(value, big.NewInt(unit)), big.NewInt(100))
		decimalPart = new(big.Int).Div(decimalPart, big.NewInt(unit))

		return fmt.Sprintf("%d.%02d %s", tempVal, decimalPart, tokenName)
	}

	if value.Cmp(big.NewInt(unit)) >= 0 {
		return convertToString(unit)
	} else if value.Cmp(big.NewInt(milliUnit)) >= 0 {
		return convertToString(milliUnit)
	} else if value.Cmp(big.NewInt(microUnit)) >= 0 {
		return convertToString(microUnit)
	} else if value.Cmp(big.NewInt(nanoUnit)) >= 0 {
		return convertToString(nanoUnit)
	} else if value.Cmp(big.NewInt(picoUnit)) >= 0 {
		return convertToString(picoUnit)
	} else if value.Cmp(big.NewInt(femtoUnit)) >= 0 {
		return convertToString(femtoUnit)
	} else {
		return fmt.Sprintf("%d atto%s", value.Int64(), tokenName)
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
