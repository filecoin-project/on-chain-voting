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

// Constants defining key types and signature types.
const (
	KTBLS       = "bls"       // BLS key type.
	KTSecp256k1 = "secp256k1" // Secp256k1 key type.

	SigTypeSecp256k1 = 1 // Signature type for Secp256k1.
	SigTypeBLS       = 2 // Signature type for BLS.
)

// Payload represents the payload of a UCAN (UnixFS Content Addressed Network) token.
type Payload struct {
	Iss string `json:"iss" bson:"iss"` // Issuer of the token.
	Aud string `json:"aud" bson:"aud"` // Audience of the token.
	Act string `json:"act" bson:"act"` // Action permitted by the token.
	Prf string `json:"prf" bson:"prf"` // Profile of the token.
}

// Header represents the header of a UCAN token.
type Header struct {
	Alg     string `json:"alg" bson:"alg"`         // Algorithm used for signing.
	Type    string `json:"type" bson:"type"`       // Type of the token.
	Version string `json:"version" bson:"version"` // Version of the token.
}
