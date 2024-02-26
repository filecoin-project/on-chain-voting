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

package models

import "math/big"

type Power struct {
	FipEditorPower   *big.Int `json:"fipEditorPower"`
	DeveloperPower   *big.Int `json:"developerPower"`
	SpPower          [][]byte `json:"spPower"`
	ClientPower      [][]byte `json:"clientPower"`
	TokenHolderPower *big.Int `json:"tokenHolderPower"`
}
