// Package types defines the core MEV protocol types shared between client and gateway.
package types

import "time"

// TransactionMessage represents a transaction message.
type TransactionMessage struct {
	Type      string `json:"type"`
	TxID      string `json:"tx_id"`
	From      string `json:"from,omitempty"`
	To        string `json:"to,omitempty"`
	Amount    int64  `json:"amount,omitempty"`
	Fee       int64  `json:"fee,omitempty"`
	Value     string `json:"value,omitempty"`
	Chain     string `json:"chain"`
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature,omitempty"`
	PubKey    string `json:"pubkey,omitempty"`
	Hops      int    `json:"hops,omitempty"`
	MaxHops   int    `json:"max_hops,omitempty"`
	Trace     string `json:"trace,omitempty"`
}

// BundleMessage represents a bundle of transactions.
type BundleMessage struct {
	Type           string               `json:"type"`
	BundleID       string               `json:"bundle_id"`
	Transactions   []TransactionMessage `json:"transactions"`
	Timestamp      int64                `json:"timestamp"`
	Signature      string               `json:"signature,omitempty"`
	PubKey         string               `json:"pubkey,omitempty"`
	Hops           int                  `json:"hops,omitempty"`
	MaxHops        int                  `json:"max_hops,omitempty"`
	Trace          string               `json:"trace,omitempty"`
	BidWei         string               `json:"bid_wei,omitempty"`
	OriginatorID   string               `json:"originator_id,omitempty"`
	TargetBlock    string               `json:"target_block,omitempty"`
	RevertingTxIDs []string             `json:"reverting_tx_ids,omitempty"`
}

// Block represents a block in the blockchain.
type Block struct {
	Header       BlockHeader          `json:"header"`
	BlockID      string               `json:"block_id"`
	Bundles      []BundleMessage      `json:"bundles"`
	Transactions []TransactionMessage `json:"transactions"`
	Timestamp    int64                `json:"timestamp"`
	TotalProfit  string               `json:"total_profit,omitempty"`
	Signature    string               `json:"signature,omitempty"`
	PubKey       string               `json:"pubkey,omitempty"`
	Trace        string               `json:"trace,omitempty"`
}

// BlockHeader represents the header of a block.
type BlockHeader struct {
	BlockID    string `json:"block_id"`
	ParentID   string `json:"parent_id"`
	MerkleRoot string `json:"merkle_root"`
	Timestamp  int64  `json:"timestamp"`
}

// ProtectedTransaction represents a transaction submitted to the private protection pool.
type ProtectedTransaction struct {
	TxID            string   `json:"tx_id"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	Value           string   `json:"value"`
	Payload         string   `json:"payload"`
	Chain           string   `json:"chain"`
	OriginatorID    string   `json:"originator_id"`
	Timestamp       int64    `json:"timestamp"`
	Deadline        int64    `json:"deadline,omitempty"`
	MinOutputWei    string   `json:"min_output_wei,omitempty"`
	MaxSlippageBps  int      `json:"max_slippage_bps,omitempty"`
	PrivacyMode     string   `json:"privacy_mode,omitempty"`
	AllowedPaths    []string `json:"allowed_paths,omitempty"`
	RebateRecipient string   `json:"rebate_recipient,omitempty"`
	HintTokenPair   string   `json:"hint_token_pair,omitempty"`
	HintChain       string   `json:"hint_chain,omitempty"`
}

// BundleSummary is a lightweight view returned by ListBundles.
type BundleSummary struct {
	BundleID       string    `json:"bundle_id"`
	BidWei         string    `json:"bid_wei"`
	EffectiveValue string    `json:"effective_value"`
	Simulated      bool      `json:"simulated"`
	Reverted       bool      `json:"reverted"`
	CreatedAt      time.Time `json:"created_at"`
}

// BlockSummary is a lightweight view returned by ListBlocks.
type BlockSummary struct {
	BlockID     string    `json:"block_id"`
	Timestamp   int64     `json:"timestamp"`
	BundleCount int       `json:"bundle_count"`
	TotalProfit string    `json:"total_profit"`
	CreatedAt   time.Time `json:"created_at"`
}
