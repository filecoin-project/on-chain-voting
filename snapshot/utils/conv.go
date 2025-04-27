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
	"math/big"
	"strconv"

	"go.uber.org/zap"
)

func StringToBigInt(v string) *big.Int {
	if len(v) == 0 {
		return big.NewInt(0)
	}

	res, ok := big.NewInt(0).SetString(v, 10)
	if ok {
		return res
	}

	zap.L().Warn("failed to convert string to big.Int", zap.Any("convert value", v))
	return big.NewInt(0)
}

func SafeParseInt(v string) int64 {
	res, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return res
}
