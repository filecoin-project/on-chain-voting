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

const (
	KTBLS       = "bls"
	KTSecp256k1 = "secp256k1"

	SigTypeSecp256k1 = 1
	SigTypeBLS       = 2
)

// Payload ucan payload
type Payload struct {
	Iss string `json:"iss" bson:"iss"`
	Aud string `json:"aud" bson:"aud"`
	Act string `json:"act" bson:"act"`
	Prf string `json:"prf" bson:"prf"`
}

// Header ucan header
type Header struct {
	Alg     string `json:"alg" bson:"alg"`
	Type    string `json:"type" bson:"type"`
	Version string `json:"version" bson:"version"`
}
