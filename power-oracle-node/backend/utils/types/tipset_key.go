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

package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	mh "github.com/multiformats/go-multihash"

	typegen "github.com/whyrusleeping/cbor-gen"
)

var EmptyTSK = TipSetKey{}

// The length of a block header CID in bytes.
var blockHeaderCIDLen int

var (
	// HashFunction is the default hash function for computing CIDs.
	//
	// This is currently Blake2b-256.
	HashFunction = uint64(mh.BLAKE2B_MIN + 31)

	// CIDInlineLimit When producing a CID for an IPLD block less than or equal to CIDInlineLimit
	// bytes in length, the identity hash function will be used instead of
	// HashFunction. This will effectively "inline" the block into the CID, allowing
	// it to be extracted directly from the CID with no disk/network operations.
	//
	// This is currently -1 for "disabled".
	//
	// This is exposed for testing. Do not modify unless you know what you're doing.
	CIDInlineLimit = -1
)

//	type cidBuilder struct {
//		codec uint64
//	}
type V1Builder struct {
	Codec    uint64
	MhType   uint64
	MhLength int // MhLength <= 0 means the default length
}

func Sum(data []byte) (Cid, error) {
	hf := HashFunction
	if len(data) <= CIDInlineLimit {
		hf = mh.IDENTITY
	}
	return V1Builder{Codec: DagCBOR, MhType: hf}.Sum(data)
}

func (p V1Builder) Sum(data []byte) (Cid, error) {
	mhLen := p.MhLength
	if mhLen <= 0 {
		mhLen = -1
	}
	hash, err := mh.Sum(data, p.MhType, mhLen)
	if err != nil {
		return UndefCid, err
	}
	return NewCidV1(p.Codec, hash), nil
}

// var CidBuilder V1Builder = cidBuilder{codec: DagCBOR}
func init() {
	// hash a large string of zeros so we don't estimate based on inlined CIDs.
	var buf [256]byte
	c, err := Sum(buf[:])
	if err != nil {
		panic(err)
	}
	blockHeaderCIDLen = len(c.Bytes())
}

// A TipSetKey is an immutable collection of CIDs forming a unique key for a tipset.
// The CIDs are assumed to be distinct and in canonical order. Two keys with the same
// CIDs in a different order are not considered equal.
// TipSetKey is a lightweight value type, and may be compared for equality with ==.
type TipSetKey struct {
	// The internal representation is a concatenation of the bytes of the CIDs, which are
	// self-describing, wrapped as a string.
	// These gymnastics make the a TipSetKey usable as a map key.
	// The empty key has value "".
	value string
}

// TipSetKeyFromBytes wraps an encoded key, validating correct decoding.
func TipSetKeyFromBytes(encoded []byte) (TipSetKey, error) {
	_, err := decodeKey(encoded)
	if err != nil {
		return EmptyTSK, err
	}
	return TipSetKey{string(encoded)}, nil
}

// Cids returns a slice of the CIDs comprising this key.
func (k TipSetKey) Cids() []Cid {
	cids, err := decodeKey([]byte(k.value))
	if err != nil {
		panic("invalid tipset key: " + err.Error())
	}
	return cids
}

// String() returns a human-readable representation of the key.
func (k TipSetKey) String() string {
	b := strings.Builder{}
	b.WriteString("{")
	cids := k.Cids()
	for i, c := range cids {
		b.WriteString(c.String())
		if i < len(cids)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("}")
	return b.String()
}

// Bytes() returns a binary representation of the key.
func (k TipSetKey) Bytes() []byte {
	return []byte(k.value)
}

func (k TipSetKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Cids())
}

func (k *TipSetKey) UnmarshalJSON(b []byte) error {
	var cids []Cid
	if err := json.Unmarshal(b, &cids); err != nil {
		return err
	}
	k.value = string(encodeKey(cids))
	return nil
}

func (k TipSetKey) MarshalCBOR(writer io.Writer) error {
	if err := typegen.WriteMajorTypeHeader(writer, typegen.MajByteString, uint64(len(k.Bytes()))); err != nil {
		return err
	}

	_, err := writer.Write(k.Bytes())
	return err
}

func (k *TipSetKey) UnmarshalCBOR(reader io.Reader) error {
	cr := typegen.NewCborReader(reader)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if extra > typegen.ByteArrayMaxLen {
		return fmt.Errorf("t.Binary: byte array too large (%d)", extra)
	}
	if maj != typegen.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	b := make([]uint8, extra)

	if _, err := io.ReadFull(cr, b); err != nil {
		return err
	}

	*k, err = TipSetKeyFromBytes(b)
	return err
}

func (k TipSetKey) IsEmpty() bool {
	return len(k.value) == 0
}

func encodeKey(cids []Cid) []byte {
	buffer := new(bytes.Buffer)
	for _, c := range cids {
		// bytes.Buffer.Write() err is documented to be always nil.
		_, _ = buffer.Write(c.Bytes())
	}
	return buffer.Bytes()
}

func decodeKey(encoded []byte) ([]Cid, error) {
	// To avoid reallocation of the underlying array, estimate the number of CIDs to be extracted
	// by dividing the encoded length by the expected CID length.
	estimatedCount := len(encoded) / blockHeaderCIDLen
	cids := make([]Cid, 0, estimatedCount)
	nextIdx := 0
	for nextIdx < len(encoded) {
		nr, c, err := CidFromBytes(encoded[nextIdx:])
		if err != nil {
			return nil, err
		}
		cids = append(cids, c)
		nextIdx += nr
	}
	return cids, nil
}

var (
	_ typegen.CBORMarshaler   = &TipSetKey{}
	_ typegen.CBORUnmarshaler = &TipSetKey{}
)
