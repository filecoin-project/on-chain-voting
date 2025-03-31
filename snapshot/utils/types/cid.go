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
	"encoding"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	mbase "github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
	varint "github.com/multiformats/go-varint"
)

var (
	// ErrCidTooShort means that the cid passed to decode was not long
	// enough to be a valid Cid
	ErrCidTooShort = errors.New("cid too short")

	// ErrInvalidBaseEncoding means that selected encoding is not supported
	// by this Cid version
	ErrInvalidBaseEncoding = errors.New("invalid base encoding")
)

// Consts below are DEPRECATED and left only for legacy reasons:
// <https://github.com/ipfs/go-cid/pull/137>
// Modern code should use consts from go-multicodec instead:
// <https://github.com/multiformats/go-multicodec>
const (
	DagProtobuf = 0x70 // https://ipld.io/docs/codecs/known/dag-pb/
	DagCBOR     = 0x71 // https://ipld.io/docs/codecs/known/dag-cbor/
)

// tryNewCidV0 tries to convert a multihash into a CIDv0 CID and returns an
// error on failure.
func tryNewCidV0(mhash mh.Multihash) (Cid, error) {
	// Need to make sure hash is valid for CidV0 otherwise we will
	// incorrectly detect it as CidV1 in the Version() method
	dec, err := mh.Decode(mhash)
	if err != nil {
		return UndefCid, err
	}
	if dec.Code != mh.SHA2_256 || dec.Length != 32 {
		return UndefCid, fmt.Errorf("invalid hash for cidv0 %d-%d", dec.Code, dec.Length)
	}
	return Cid{string(mhash)}, nil
}

// NewCidV0 returns a Cid-wrapped multihash.
// They exist to allow IPFS to work with Cids while keeping
// compatibility with the plain-multihash format used used in IPFS.
// NewCidV1 should be used preferentially.
//
// Panics if the multihash isn't sha2-256.
func NewCidV0(mhash mh.Multihash) Cid {
	c, err := tryNewCidV0(mhash)
	if err != nil {
		panic(err)
	}
	return c
}

// NewCidV1 returns a new Cid using the given multicodec-packed
// content type.
//
// Panics if the multihash is invalid.
func NewCidV1(codecType uint64, mhash mh.Multihash) Cid {
	hashlen := len(mhash)

	// Two 8 bytes (max) numbers plus hash.
	// We use strings.Builder to only allocate once.
	var b strings.Builder
	b.Grow(1 + varint.UvarintSize(codecType) + hashlen)

	b.WriteByte(1)

	var buf [binary.MaxVarintLen64]byte
	n := varint.PutUvarint(buf[:], codecType)
	b.Write(buf[:n])

	cn, _ := b.Write(mhash)
	if cn != hashlen {
		panic("copy hash length is inconsistent")
	}

	return Cid{b.String()}
}

var (
	_ encoding.BinaryMarshaler   = Cid{}
	_ encoding.BinaryUnmarshaler = (*Cid)(nil)
	_ encoding.TextMarshaler     = Cid{}
	_ encoding.TextUnmarshaler   = (*Cid)(nil)
)

// Cid represents a self-describing content addressed
// identifier. It is formed by a Version, a Codec (which indicates
// a multicodec-packed content type) and a Multihash.
type Cid struct{ str string }

// Undef can be used to represent a nil or undefined Cid, using Cid{}
// directly is also acceptable.
var UndefCid = Cid{}

// Defined returns true if a Cid is defined
// Calling any other methods on an undefined Cid will result in
// undefined behavior.
func (c Cid) Defined() bool {
	return c.str != ""
}

// Parse is a short-hand function to perform Decode, Cast etc... on
// a generic interface{} type.
func Parse(v interface{}) (Cid, error) {
	switch v2 := v.(type) {
	case string:
		if strings.Contains(v2, "/ipfs/") {
			return Decode(strings.Split(v2, "/ipfs/")[1])
		}
		return Decode(v2)
	case []byte:
		return Cast(v2)
	case mh.Multihash:
		return tryNewCidV0(v2)
	case Cid:
		return v2, nil
	default:
		return UndefCid, fmt.Errorf("can't parse %+v as Cid", v2)
	}
}

// Decode parses a Cid-encoded string and returns a Cid object.
// For CidV1, a Cid-encoded string is primarily a multibase string:
//
//	<multibase-type-code><base-encoded-string>
//
// The base-encoded string represents a:
//
// <version><codec-type><multihash>
//
// Decode will also detect and parse CidV0 strings. Strings
// starting with "Qm" are considered CidV0 and treated directly
// as B58-encoded multihashes.
func Decode(v string) (Cid, error) {
	if len(v) < 2 {
		return UndefCid, ErrCidTooShort
	}

	if len(v) == 46 && v[:2] == "Qm" {
		hash, err := mh.FromB58String(v)
		if err != nil {
			return UndefCid, err
		}

		return tryNewCidV0(hash)
	}

	_, data, err := mbase.Decode(v)
	if err != nil {
		return UndefCid, err
	}

	return Cast(data)
}

// Extract the encoding from a Cid.  If Decode on the same string did
// not return an error neither will this function.
func ExtractEncoding(v string) (mbase.Encoding, error) {
	if len(v) < 2 {
		return -1, ErrCidTooShort
	}

	if len(v) == 46 && v[:2] == "Qm" {
		return mbase.Base58BTC, nil
	}

	encoding := mbase.Encoding(v[0])

	// check encoding is valid
	_, err := mbase.NewEncoder(encoding)
	if err != nil {
		return -1, err
	}

	return encoding, nil
}

// Cast takes a Cid data slice, parses it and returns a Cid.
// For CidV1, the data buffer is in the form:
//
//	<version><codec-type><multihash>
//
// CidV0 are also supported. In particular, data buffers starting
// with length 34 bytes, which starts with bytes [18,32...] are considered
// binary multihashes.
//
// Please use decode when parsing a regular Cid string, as Cast does not
// expect multibase-encoded data. Cast accepts the output of Cid.Bytes().
func Cast(data []byte) (Cid, error) {
	nr, c, err := CidFromBytes(data)
	if err != nil {
		return UndefCid, err
	}

	if nr != len(data) {
		return UndefCid, fmt.Errorf("trailing bytes in data buffer passed to cid Cast")
	}

	return c, nil
}

// UnmarshalBinary is equivalent to Cast(). It implements the
// encoding.BinaryUnmarshaler interface.
func (c *Cid) UnmarshalBinary(data []byte) error {
	casted, err := Cast(data)
	if err != nil {
		return err
	}
	c.str = casted.str
	return nil
}

// UnmarshalText is equivalent to Decode(). It implements the
// encoding.TextUnmarshaler interface.
func (c *Cid) UnmarshalText(text []byte) error {
	decodedCid, err := Decode(string(text))
	if err != nil {
		return err
	}
	c.str = decodedCid.str
	return nil
}

// Version returns the Cid version.
func (c Cid) Version() uint64 {
	if len(c.str) == 34 && c.str[0] == 18 && c.str[1] == 32 {
		return 0
	}
	return 1
}

// Type returns the multicodec-packed content type of a Cid.
func (c Cid) Type() uint64 {
	if c.Version() == 0 {
		return DagProtobuf
	}
	_, n, _ := uvarint(c.str)
	codec, _, _ := uvarint(c.str[n:])
	return codec
}

// String returns the default string representation of a
// Cid. Currently, Base32 is used for CIDV1 as the encoding for the
// multibase string, Base58 is used for CIDV0.
func (c Cid) String() string {
	switch c.Version() {
	case 0:
		return c.Hash().B58String()
	case 1:
		mbstr, err := mbase.Encode(mbase.Base32, c.Bytes())
		if err != nil {
			panic("should not error with hardcoded mbase: " + err.Error())
		}

		return mbstr
	default:
		panic("not possible to reach this point")
	}
}

// String returns the string representation of a Cid
// encoded is selected base
func (c Cid) StringOfBase(base mbase.Encoding) (string, error) {
	switch c.Version() {
	case 0:
		if base != mbase.Base58BTC {
			return "", ErrInvalidBaseEncoding
		}
		return c.Hash().B58String(), nil
	case 1:
		return mbase.Encode(base, c.Bytes())
	default:
		panic("not possible to reach this point")
	}
}

// Encode return the string representation of a Cid in a given base
// when applicable.  Version 0 Cid's are always in Base58 as they do
// not take a multibase prefix.
func (c Cid) Encode(base mbase.Encoder) string {
	switch c.Version() {
	case 0:
		return c.Hash().B58String()
	case 1:
		return base.Encode(c.Bytes())
	default:
		panic("not possible to reach this point")
	}
}

// Hash returns the multihash contained by a Cid.
func (c Cid) Hash() mh.Multihash {
	bytes := c.Bytes()

	if c.Version() == 0 {
		return mh.Multihash(bytes)
	}

	// skip version length
	_, n1, _ := varint.FromUvarint(bytes)
	// skip codec length
	_, n2, _ := varint.FromUvarint(bytes[n1:])

	return mh.Multihash(bytes[n1+n2:])
}

// Bytes returns the byte representation of a Cid.
// The output of bytes can be parsed back into a Cid
// with Cast().
func (c Cid) Bytes() []byte {
	return []byte(c.str)
}

// ByteLen returns the length of the CID in bytes.
// It's equivalent to `len(c.Bytes())`, but works without an allocation,
// and should therefore be preferred.
//
// (See also the WriteTo method for other important operations that work without allocation.)
func (c Cid) ByteLen() int {
	return len(c.str)
}

// WriteBytes writes the CID bytes to the given writer.
// This method works without incurring any allocation.
//
// (See also the ByteLen method for other important operations that work without allocation.)
func (c Cid) WriteBytes(w io.Writer) (int, error) {
	n, err := io.WriteString(w, c.str)
	if err != nil {
		return n, err
	}
	if n != len(c.str) {
		return n, fmt.Errorf("failed to write entire cid string")
	}
	return n, nil
}

// MarshalBinary is equivalent to Bytes(). It implements the
// encoding.BinaryMarshaler interface.
func (c Cid) MarshalBinary() ([]byte, error) {
	return c.Bytes(), nil
}

// MarshalText is equivalent to String(). It implements the
// encoding.TextMarshaler interface.
func (c Cid) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// Equals checks that two Cids are the same.
// In order for two Cids to be considered equal, the
// Version, the Codec and the Multihash must match.
func (c Cid) Equals(o Cid) bool {
	return c == o
}

// UnmarshalJSON parses the JSON representation of a Cid.
func (c *Cid) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("invalid cid json blob")
	}
	obj := struct {
		CidTarget string `json:"/"`
	}{}
	objptr := &obj
	err := json.Unmarshal(b, &objptr)
	if err != nil {
		return err
	}
	if objptr == nil {
		*c = Cid{}
		return nil
	}

	if obj.CidTarget == "" {
		return fmt.Errorf("cid was incorrectly formatted")
	}

	out, err := Decode(obj.CidTarget)
	if err != nil {
		return err
	}

	*c = out

	return nil
}

// MarshalJSON procudes a JSON representation of a Cid, which looks as follows:
//
//	{ "/": "<cid-string>" }
//
// Note that this formatting comes from the IPLD specification
// (https://github.com/ipld/specs/tree/master/ipld)
func (c Cid) MarshalJSON() ([]byte, error) {
	if !c.Defined() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("{\"/\":\"%s\"}", c.String())), nil
}

// KeyString returns the binary representation of the Cid as a string
func (c Cid) KeyString() string {
	return c.str
}

// Loggable returns a Loggable (as defined by
// https://godoc.org/github.com/ipfs/go-log).
func (c Cid) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"cid": c,
	}
}

// Prefix builds and returns a Prefix out of a Cid.
func (c Cid) Prefix() Prefix {
	if c.Version() == 0 {
		return Prefix{
			MhType:   mh.SHA2_256,
			MhLength: 32,
			Version:  0,
			Codec:    DagProtobuf,
		}
	}

	offset := 0
	version, n, _ := uvarint(c.str[offset:])
	offset += n
	codec, n, _ := uvarint(c.str[offset:])
	offset += n
	mhtype, n, _ := uvarint(c.str[offset:])
	offset += n
	mhlen, _, _ := uvarint(c.str[offset:])

	return Prefix{
		MhType:   mhtype,
		MhLength: int(mhlen),
		Version:  version,
		Codec:    codec,
	}
}

// Prefix represents all the metadata of a Cid,
// that is, the Version, the Codec, the Multihash type
// and the Multihash length. It does not contains
// any actual content information.
// NOTE: The use -1 in MhLength to mean default length is deprecated,
//
//	use the V0Builder or V1Builder structures instead
type Prefix struct {
	Version  uint64
	Codec    uint64
	MhType   uint64
	MhLength int
}

// Sum uses the information in a prefix to perform a multihash.Sum()
// and return a newly constructed Cid with the resulting multihash.
func (p Prefix) Sum(data []byte) (Cid, error) {
	length := p.MhLength
	if p.MhType == mh.IDENTITY {
		length = -1
	}

	if p.Version == 0 && (p.MhType != mh.SHA2_256 ||
		(p.MhLength != 32 && p.MhLength != -1)) {
		return UndefCid, fmt.Errorf("invalid v0 prefix")
	}

	hash, err := mh.Sum(data, p.MhType, length)
	if err != nil {
		return UndefCid, err
	}

	switch p.Version {
	case 0:
		return NewCidV0(hash), nil
	case 1:
		return NewCidV1(p.Codec, hash), nil
	default:
		return UndefCid, fmt.Errorf("invalid cid version")
	}
}

// Bytes returns a byte representation of a Prefix. It looks like:
//
//	<version><codec><mh-type><mh-length>
func (p Prefix) Bytes() []byte {
	size := varint.UvarintSize(p.Version)
	size += varint.UvarintSize(p.Codec)
	size += varint.UvarintSize(p.MhType)
	size += varint.UvarintSize(uint64(p.MhLength))

	buf := make([]byte, size)
	n := varint.PutUvarint(buf, p.Version)
	n += varint.PutUvarint(buf[n:], p.Codec)
	n += varint.PutUvarint(buf[n:], p.MhType)
	n += varint.PutUvarint(buf[n:], uint64(p.MhLength))
	if n != size {
		panic("size mismatch")
	}
	return buf
}

// PrefixFromBytes parses a Prefix-byte representation onto a
// Prefix.
func PrefixFromBytes(buf []byte) (Prefix, error) {
	r := bytes.NewReader(buf)
	vers, err := varint.ReadUvarint(r)
	if err != nil {
		return Prefix{}, err
	}

	codec, err := varint.ReadUvarint(r)
	if err != nil {
		return Prefix{}, err
	}

	mhtype, err := varint.ReadUvarint(r)
	if err != nil {
		return Prefix{}, err
	}

	mhlen, err := varint.ReadUvarint(r)
	if err != nil {
		return Prefix{}, err
	}

	return Prefix{
		Version:  vers,
		Codec:    codec,
		MhType:   mhtype,
		MhLength: int(mhlen),
	}, nil
}

func CidFromBytes(data []byte) (int, Cid, error) {
	if len(data) > 2 && data[0] == mh.SHA2_256 && data[1] == 32 {
		if len(data) < 34 {
			return 0, UndefCid, fmt.Errorf("not enough bytes for cid v0")
		}

		h, err := mh.Cast(data[:34])
		if err != nil {
			return 0, UndefCid, err
		}

		return 34, Cid{string(h)}, nil
	}

	vers, n, err := varint.FromUvarint(data)
	if err != nil {
		return 0, UndefCid, err
	}

	if vers != 1 {
		return 0, UndefCid, fmt.Errorf("expected 1 as the cid version number, got: %d", vers)
	}

	_, cn, err := varint.FromUvarint(data[n:])
	if err != nil {
		return 0, UndefCid, err
	}

	mhnr, _, err := mh.MHFromBytes(data[n+cn:])
	if err != nil {
		return 0, UndefCid, err
	}

	l := n + cn + mhnr

	return l, Cid{string(data[0:l])}, nil
}

func toBufByteReader(r io.Reader, dst []byte) *bufByteReader {
	// If the reader already implements ByteReader, use it directly.
	// Otherwise, use a fallback that does 1-byte Reads.
	if br, ok := r.(io.ByteReader); ok {
		return &bufByteReader{direct: br, dst: dst}
	}
	return &bufByteReader{fallback: r, dst: dst}
}

type bufByteReader struct {
	direct   io.ByteReader
	fallback io.Reader

	dst []byte
}

func (r *bufByteReader) ReadByte() (byte, error) {
	// The underlying reader has ReadByte; use it.
	if br := r.direct; br != nil {
		b, err := br.ReadByte()
		if err != nil {
			return 0, err
		}
		r.dst = append(r.dst, b)
		return b, nil
	}

	// Fall back to a one-byte Read.
	// once we have benchmarks and if they prove that to be faster.
	var p [1]byte
	if _, err := io.ReadFull(r.fallback, p[:]); err != nil {
		return 0, err
	}
	r.dst = append(r.dst, p[0])
	return p[0], nil
}

// CidFromReader reads a precise number of bytes for a CID from a given reader.
// It returns the number of bytes read, the CID, and any error encountered.
// The number of bytes read is accurate even if a non-nil error is returned.
//
// It's recommended to supply a reader that buffers and implements io.ByteReader,
// as CidFromReader has to do many single-byte reads to decode varints.
// If the argument only implements io.Reader, single-byte Read calls are used instead.
func CidFromReader(r io.Reader) (int, Cid, error) {
	// 64 bytes is enough for any CIDv0,
	// and it's enough for most CIDv1s in practice.
	// If the digest is too long, we'll allocate more.
	br := toBufByteReader(r, make([]byte, 0, 64))

	// We read the first varint, to tell if this is a CIDv0 or a CIDv1.
	// The varint package wants a io.ByteReader, so we must wrap our io.Reader.
	vers, err := varint.ReadUvarint(br)
	if err != nil {
		return len(br.dst), UndefCid, err
	}

	// If we have a CIDv0, read the rest of the bytes and cast the buffer.
	if vers == mh.SHA2_256 {
		if n, err := io.ReadFull(r, br.dst[1:34]); err != nil {
			return len(br.dst) + n, UndefCid, err
		}

		br.dst = br.dst[:34]
		h, err := mh.Cast(br.dst)
		if err != nil {
			return len(br.dst), UndefCid, err
		}

		return len(br.dst), Cid{string(h)}, nil
	}

	if vers != 1 {
		return len(br.dst), UndefCid, fmt.Errorf("expected 1 as the cid version number, got: %d", vers)
	}

	// CID block encoding multicodec.
	_, err = varint.ReadUvarint(br)
	if err != nil {
		return len(br.dst), UndefCid, err
	}

	// We could replace most of the code below with go-multihash's ReadMultihash.
	// Note that it would save code, but prevent reusing buffers.
	// Plus, we already have a ByteReader now.
	mhStart := len(br.dst)

	// Multihash hash function code.
	_, err = varint.ReadUvarint(br)
	if err != nil {
		return len(br.dst), UndefCid, err
	}

	// Multihash digest length.
	mhl, err := varint.ReadUvarint(br)
	if err != nil {
		return len(br.dst), UndefCid, err
	}

	// Refuse to make large allocations to prevent OOMs due to bugs.
	const maxDigestAlloc = 32 << 20 // 32MiB
	if mhl > maxDigestAlloc {
		return len(br.dst), UndefCid, fmt.Errorf("refusing to allocate %d bytes for a digest", mhl)
	}

	// Fine to convert mhl to int, given maxDigestAlloc.
	prefixLength := len(br.dst)
	cidLength := prefixLength + int(mhl)
	if cidLength > cap(br.dst) {
		// If the multihash digest doesn't fit in our initial 64 bytes,
		// efficiently extend the slice via append+make.
		br.dst = append(br.dst, make([]byte, cidLength-len(br.dst))...)
	} else {
		// The multihash digest fits inside our buffer,
		// so just extend its capacity.
		br.dst = br.dst[:cidLength]
	}

	if n, err := io.ReadFull(r, br.dst[prefixLength:cidLength]); err != nil {
		// We can't use len(br.dst) here,
		// as we've only read n bytes past prefixLength.
		return prefixLength + n, UndefCid, err
	}

	// This simply ensures the multihash is valid.
	// for now, it helps ensure consistency with CidFromBytes.
	_, _, err = mh.MHFromBytes(br.dst[mhStart:])
	if err != nil {
		return len(br.dst), UndefCid, err
	}

	return len(br.dst), Cid{string(br.dst)}, nil
}

// uvarint decodes a uint64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, then 0 is
// returned for both the value and the number of bytes read, and an
// error is returned.
func uvarint(buf string) (uint64, int, error) {
	var x uint64
	var s uint
	// we have a binary string so we can't use a range loop
	for i := 0; i < len(buf); i++ {
		b := buf[i]
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, 0, varint.ErrOverflow
			}
			if b == 0 && i > 0 {
				return 0, 0, varint.ErrNotMinimal
			}
			return x | uint64(b)<<s, i + 1, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0, varint.ErrUnderflow
}
