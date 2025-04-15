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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringConvToBigInt(t *testing.T) {
	strNumber := "1234567890123456789012345678901234567890"
	bigNumber := StringConvToBigInt(strNumber)
	if bigNumber.String() != strNumber {
		t.Errorf("Expected %s, got %s", strNumber, bigNumber.String())
	}
}

func TestBigIntConvToString(t *testing.T) {
	// 1 Ether = 10^18 Wei
	strNumber := "1000000000000000000"
	bigNumber := StringConvToBigInt(strNumber)
	strResult := DividedBy10To18(bigNumber, 0)
	if strResult != "1" {
		t.Errorf("Expected %s, got %s", "1", strResult)
	}
}

func TestStringToBase64URL(t *testing.T) {
	res := StringToBase64URL("test")
	assert.NotEmpty(t, res)
}
