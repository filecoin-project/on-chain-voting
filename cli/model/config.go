package model

// Config holds the configuration details for various services.
type Config struct {
	Drand   Drand   // Drand network configuration
	Network Network // Network configuration
	ABIPath ABIPath // Path to the ABI for PowerVoting contract
}

// Drand contains configuration for the Drand network.
type Drand struct {
	Urls      []string // List of URLs for the Drand network
	ChainHash string   // Chain hash for the Drand network
}

// Network defines the configuration for the blockchain network.
type Network struct {
	ChainID             int    // Chain identifier (e.g., mainnet, testnet)
	RPC                 string // RPC endpoint for the network
	Token               string // Token used in the network
	PowerVotingContract string // Contract address for PowerVoting
	PowerBackendURL     string // Backend URL for power service
}

// ABIPath holds the path or identifier for the PowerVoting contract's ABI.
type ABIPath struct {
	PowerVotingABI string // ABI for PowerVoting contract
}
