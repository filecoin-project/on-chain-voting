package utils

import (
	"fmt"
	"math/big"
)

const (
	_  = iota
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB
	GB
	TB
	PiB
	EiB
)

func ConvertSize(size int64) string {
	switch {
	case size >= EiB:
		return fmt.Sprintf("%.2f EiB", float64(size)/float64(EiB))
	case size >= PiB:
		return fmt.Sprintf("%.2f PiB", float64(size)/float64(PiB))
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

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
