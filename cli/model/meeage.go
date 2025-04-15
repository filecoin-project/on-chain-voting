package model

import (
	"github.com/filecoin-project/lotus/chain/types"
)

// Message represents a message to be sent in the Filecoin network.
type Message struct {
	Version    int             `json:"version"`    // Version of the message format
	To         string          `json:"to"`         // Recipient address
	From       string          `json:"from"`       // Sender address
	Nonce      int             `json:"nonce"`      // Nonce for transaction (used to prevent replay attacks)
	Value      string          `json:"value"`      // Value to send with the message (e.g., amount of FIL)
	GasLimit   int             `json:"gasLimit"`   // Gas limit for the transaction
	GasFeeCap  string          `json:"gasFeeCap"`  // Maximum fee cap for the transaction
	GasPremium string          `json:"gasPremium"` // Premium fee for the transaction
	Method     int             `json:"method"`     // Method ID for the specific method to call
	Params     string          `json:"params"`     // Encoded parameters for the method call
	Cid        types.TipSetKey `json:"cid"`        // CID (Transaction ID) for the message (optional)
}
