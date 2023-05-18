package common

import (
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self   t.ModuleID // id of this module
	Hasher t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	MinTransactionsInBatch int
	MaxTransactionsInBatch int

	// Divides transactions among N buckets, based on the transaction ID.
	// Upon batch request, it will fill the batch with transactions starting from a pseudo-randomly
	// chosen bucket.
	// Must be at least 1.
	IncomingTxBucketCount int

	// Used for choosing a bucket. Can be e.g. an index in the list of nodes.
	RandSeed int64

	// If this parameter is not nil, the mempool will not receive transactions directly (through NewTransactions) events.
	// On reception of such an event, it will report an error (making the system crash).
	// Instead, TxFetcher will be called to pull transactions from an external source
	// when they are needed to form a batch (upon the RequestBatch event).
	// Looking up transactions by ID will also always fail (return no transactions).
	TxFetcher func() []*trantorpbtypes.Transaction
}

// State represents the common state accessible to all parts of the module implementation.
type State struct {
	TxByID map[tt.TxID]*trantorpbtypes.Transaction
}
