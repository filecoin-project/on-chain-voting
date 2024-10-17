package client

import (
	"bytes"
	"net/http"
	"os"
	"powervoting-server/config"

	"github.com/ipfs/go-cid"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multihash"
	"github.com/web3-storage/go-ucanto/core/delegation"
	"github.com/web3-storage/go-ucanto/did"
	"github.com/web3-storage/go-ucanto/principal"
	"github.com/web3-storage/go-ucanto/principal/ed25519/signer"
	"github.com/web3-storage/go-w3up/capability/storeadd"
	"github.com/web3-storage/go-w3up/client"
	godelegation "github.com/web3-storage/go-w3up/delegation"
	"go.uber.org/zap"
)

var W3 W3Client

type W3Client struct {
	Space  did.DID
	Issuer principal.Signer
	Proof  delegation.Delegation
}

func InitW3Client() {
	space, err := did.Parse(config.Client.W3Client.Space)
	if err != nil {
		zap.L().Error("w3client parse space error:", zap.Error(err))
		return
	}

	issuer, err := signer.Parse(config.Client.W3Client.PrivateKey)
	if err != nil {
		zap.L().Error("w3client parse private error:", zap.Error(err))
		return
	}

	prfbytes, err := os.ReadFile(config.Client.W3Client.Proof)
	if err != nil {
		zap.L().Error("read proof.ucan file error:", zap.Error(err))
		return
	}

	proof, err := godelegation.ExtractProof(prfbytes)
	if err != nil {
		zap.L().Error("init w3storage error:", zap.Error(err))
		return
	}

	W3 = W3Client{
		Space:  space,
		Issuer: issuer,
		Proof:  proof,
	}
}

func (w *W3Client) Upload(path string) (string, error) {
	data, _ := os.ReadFile(path)

	// generate the CID for the CAR
	mh, _ := multihash.Sum(data, multihash.SHA2_256, -1)
	link := cidlink.Link{Cid: cid.NewCidV1(0x0202, mh)}
	rcpt, _ := client.StoreAdd(
		w.Issuer,
		w.Space,
		&storeadd.Caveat{Link: link, Size: uint64(len(data))},
		client.WithProofs([]delegation.Delegation{w.Proof}),
	)

	if rcpt.Out().Ok().Status == "upload" {
		hr, _ := http.NewRequest("PUT", *rcpt.Out().Ok().Url, bytes.NewReader(data))

		hdr := map[string][]string{}
		for k, v := range rcpt.Out().Ok().Headers.Values {
			hdr[k] = []string{v}
		}

		hr.Header = hdr
		hr.ContentLength = int64(len(data))
		httpClient := http.Client{}
		res, _ := httpClient.Do(hr)
		res.Body.Close()
	}

	return link.String(), nil
}
