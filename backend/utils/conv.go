// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"
)

// Convert string numbers to *big.Int type
// If the conversion fails, it returns a new big.Int with value 0
func StringConvToBigInt(s string) *big.Int {
	res, isSuccess := new(big.Int).SetString(s, 10)
	if !isSuccess {
		return big.NewInt(0)
	}
	return res
}

// DividedBy10To18 divides a big int by 10^18
// The result is to retain X decimal places
func DividedBy10To18(d *big.Int, x int32) string {
	return decimal.NewFromBigInt(d, 0).Div(decimal.New(1, 18)).StringFixed(x)
}

// bigIntDiv performs division operation on big integers and returns the result as a float64.
// Parameters x and y are pointers to big integers representing the dividend and divisor respectively.
// The return value z represents the division result as a float64.
// If the divisor is zero, it returns 0 to avoid division by zero error.
func BigIntDiv(x *big.Int, y *big.Int) decimal.Decimal {
	xd := decimal.NewFromBigInt(x, 0)
	yd := decimal.NewFromBigInt(y, 0)
	if yd.IsZero() {
		return decimal.Zero
	}

	// align precision
	return xd.Div(yd).Round(5)
}

// ParseStringToInt64 is a function that converts a string to an int64.
func ParseStringToInt64(v string) int64 {
	res, err := strconv.ParseInt(v, 10, 64)
	// If there is an error during the conversion, the function returns 0.
	if err != nil {
		return 0
	}

	return res
}

// ObjToString Serialize a structure into a string
func ObjToString(obj any) string {
	res, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(res)
}
