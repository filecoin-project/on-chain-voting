package main

import (
	"fmt"
	"testing"
)

var (
	secp256k1Aud           = "0xf5915bb5e97B857d55B4DC453A0F10490e81987e"
	secp256k1Act           = "add"
	secp256k1PrivateKeyStr = "zi2LNdLhA/MZKmNyUuZCuF6c0OUDqq8nIG8918ORHU4="
	secp256k1KeyTypeStr    = "secp256k1"
)

var (
	blsAud           = "0xf5915bb5e97B857d55B4DC453A0F10490e81987e"
	blsAct           = "add"
	blsPrivateKeyStr = "tWmF50J5RuGQlrO9g1nBkAWEnnL+jUuvPOu6ji6SaCA="
	blsKeyTypeStr    = "bls"
)

func TestSecp256k1Signature(t *testing.T) {
	ucan, err := Signature(secp256k1Aud, secp256k1Act, secp256k1PrivateKeyStr, secp256k1KeyTypeStr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ucan)
}

func TestBlsSignature(t *testing.T) {
	ucan, err := Signature(blsAud, blsAct, blsPrivateKeyStr, blsKeyTypeStr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ucan)
}
