package common

import (
	"time"

	msctypes "github.com/filecoin-project/mir/pkg/availability/multisigcollector/types"
	cryptopbtypes "github.com/filecoin-project/mir/pkg/pb/cryptopb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// InstanceUID is used to uniquely identify an instance of multisig collector.
// It is used to prevent cross-instance signature replay attack and should be unique across all executions.
type InstanceUID []byte

// Bytes returns the binary representation of the InstanceUID.
func (uid InstanceUID) Bytes() []byte {
	return uid
}

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self t.ModuleID // id of this module

	BatchDB t.ModuleID
	Crypto  t.ModuleID
	Mempool t.ModuleID
	Net     t.ModuleID
}

// ModuleParams sets the values for the parameters of an instance of the protocol.
// All replicas are expected to use identical module parameters.
type ModuleParams struct {
	// Maximal time between receiving a batch request and emitting a batch.
	// On reception of a batch request, the mempool generally waits
	// until it contains enough transactions to fill a batch (by number or by payload size)
	// and only then emits the new batch.
	// If no batch has been filled by BatchTimeout, the mempool emits a non-full (even a completely empty) batch.
	BatchTimeout time.Duration

	InstanceUID []byte                     // unique identifier for this instance used to prevent replay attacks
	Membership  *trantorpbtypes.Membership // the list of participating nodes
	Limit       int                        // the maximum number of certificates to generate before a request is completed
	MaxRequests int                        // the maximum number of requests to be provided by this module
}

// SigData is the binary data that should be signed for forming a certificate.
func SigData(instanceUID InstanceUID, batchID msctypes.BatchID) *cryptopbtypes.SignedData {
	return &cryptopbtypes.SignedData{Data: [][]byte{instanceUID.Bytes(), []byte("BATCH_STORED"), []byte(batchID)}}
}
